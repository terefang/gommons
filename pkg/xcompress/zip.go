package xcompress

// originally taken from https://github.com/kjk/common/ under:
//MIT License
//
//Copyright (c) 2021 Krzysztof Kowalczyk
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/terefang/gommons/pkg/util"
	"github.com/terefang/gommons/pkg/xfile"
)

func ZipDir(dirToZip string) ([]byte, error) {
	var buf bytes.Buffer
	err := ZipDirToWriter(&buf, dirToZip)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ZipDirToWriter(w io.Writer, dirToZip string) error {
	zw := zip.NewWriter(w)
	err := filepath.Walk(dirToZip, func(pathToZip string, info os.FileInfo, err error) error {
		if err != nil {
			//fmt.Printf("WalkFunc() received err %s from filepath.Wath()\n", err.Error())
			return err
		}
		//fmt.Printf("%s\n", path)
		isDir, err := xfile.PathIsDir(pathToZip)
		if err != nil {
			//fmt.Printf("PathIsDir() for %s failed with %s\n", pathToZip, err.Error())
			return err
		}
		if isDir {
			return nil
		}
		toZipReader, err := os.Open(pathToZip)
		if err != nil {
			//fmt.Printf("os.Open() %s failed with %s\n", pathToZip, err.Error())
			return err
		}
		defer toZipReader.Close()

		zipName := pathToZip[len(dirToZip)+1:] // +1 for '/' in the path
		// per zip.Writer.Create(), name must be slash-separated so on Windows
		// we have to convert backslashes to slashes
		zipName = filepath.ToSlash(zipName)
		inZipWriter, err := zw.Create(zipName)
		if err != nil {
			//fmt.Printf("Error in zipWriter(): %s\n", err.Error())
			return err
		}
		_, err = io.Copy(inZipWriter, toZipReader)
		if err != nil {
			return err
		}
		//fmt.Printf("Added %s to zip file\n", pathToZip)
		return nil
	})
	err2 := zw.Close()
	if err2 != nil {
		return err2
	}
	return err
}

// CreateZipWithDirContent creates a zip file with the content of a directory.
// The names of files inside the zip file are relatitve to dirToZip e.g.
// if dirToZip is foo and there is a file foo/bar.txt, the name in the zip
// will be bar.txt
func CreateZipWithDirContent(zipFilePath, dirToZip string) error {
	if isDir, err := xfile.PathIsDir(dirToZip); err != nil || !isDir {
		// TODO: should return an error if err == nil && !isDir
		return err
	}
	zf, err := os.Create(zipFilePath)
	if err != nil {
		//fmt.Printf("Failed to os.Create() %s, %s\n", zipFilePath, err.Error())
		return err
	}
	defer zf.Close()
	return ZipDirToWriter(zf, dirToZip)
}

func ReadZipFileMust(path string) map[string][]byte {
	r, err := zip.OpenReader(path)
	util.Must(err)
	defer xfile.CloseNoError(r)
	res := map[string][]byte{}
	for _, f := range r.File {
		rc, err := f.Open()
		util.Must(err)
		d, err := io.ReadAll(rc)
		util.Must(err)
		_ = rc.Close()
		res[f.Name] = d
	}
	return res
}

func zipAddFile(zw *zip.Writer, zipName string, path string) {
	zipName = filepath.ToSlash(zipName)
	d, err := os.ReadFile(path)
	util.Must(err)
	w, err := zw.Create(zipName)
	util.Must(err)
	_, err = w.Write(d)
	util.Must(err)
	fmt.Printf("  added %s from %s\n", zipName, path)
}

func zipDirRecur(zw *zip.Writer, baseDir string, dirToZip string) {
	dir := filepath.Join(baseDir, dirToZip)
	files, err := os.ReadDir(dir)
	util.Must(err)
	for _, fi := range files {
		if fi.IsDir() {
			zipDirRecur(zw, baseDir, filepath.Join(dirToZip, fi.Name()))
		} else if fi.Type().IsRegular() {
			zipName := filepath.Join(dirToZip, fi.Name())
			path := filepath.Join(baseDir, zipName)
			zipAddFile(zw, zipName, path)
		} else {
			path := filepath.Join(baseDir, fi.Name())
			s := fmt.Sprintf("%s is not a dir or regular file", path)
			panic(s)
		}
	}
}

// toZip is a list of files and directories in baseDir
// Directories are added recursively
func CreateZipFile(dst string, baseDir string, toZip ...string) {
	os.Remove(dst)
	util.PanicIf(len(toZip) == 0, "must provide toZip args")
	w, err := os.Create(dst)
	util.Must(err)
	defer xfile.CloseNoError(w)
	zw := zip.NewWriter(w)
	util.Must(err)
	for _, name := range toZip {
		path := filepath.Join(baseDir, name)
		fi, err := os.Stat(path)
		util.Must(err)
		if fi.IsDir() {
			zipDirRecur(zw, baseDir, name)
		} else if fi.Mode().IsRegular() {
			zipAddFile(zw, name, path)
		} else {
			s := fmt.Sprintf("%s is not a dir or regular file", path)
			panic(s)
		}
	}
	err = zw.Close()
	util.Must(err)
}

func UnzipDataToDir(zipData []byte, dir string) error {
	writeFile := func(f *zip.File, data []byte) error {
		// names in zip are unix-style, convert to windows-style
		name := filepath.FromSlash(f.Name)
		path := filepath.Join(dir, name)
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return err
		}
		return os.WriteFile(path, data, 0644)
	}
	return IterZipData(zipData, writeFile)
}

func IterZipReader(r *zip.Reader, cb func(f *zip.File, data []byte) error) error {
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		d, err := io.ReadAll(rc)
		err2 := rc.Close()
		if err != nil {
			return err
		}
		if err2 != nil {
			return err2
		}
		err = cb(f, d)
		if err != nil {
			return err
		}
	}
	return nil
}

func IterZipData(zipData []byte, cb func(f *zip.File, data []byte) error) error {
	dr := bytes.NewReader(zipData)
	r, err := zip.NewReader(dr, int64(len(zipData)))
	if err != nil {
		return err
	}
	return IterZipReader(r, cb)
}

func ReadZipData(zipData []byte) (map[string][]byte, error) {
	res := map[string][]byte{}
	err := IterZipData(zipData, func(f *zip.File, data []byte) error {
		res[f.Name] = data
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
