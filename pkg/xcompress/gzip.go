package xcompress

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
)

// WriteFileGzipped writes data to a path, using best gzip compression
func WriteFileGzipped(path string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	w, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		f.Close()
		os.Remove(path)
		return err
	}
	err = w.Close()
	if err != nil {
		f.Close()
		os.Remove(path)
		return err
	}
	err = f.Close()
	if err != nil {
		os.Remove(path)
		return err
	}
	return nil
}

func GzipCompressData(d []byte) ([]byte, error) {
	var dst bytes.Buffer
	w, err := gzip.NewWriterLevel(&dst, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(d)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return dst.Bytes(), nil
}

func GzipDecompressData(d []byte) ([]byte, error) {
	r := bytes.NewReader(d)
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	return io.ReadAll(zr)
}

func GzipCompressFile(dstPath, srcPath string) error {
	fSrc, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer fSrc.Close()
	fDst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer fDst.Close()
	w, err := gzip.NewWriterLevel(fDst, gzip.BestCompression)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, fSrc)
	if err != nil {
		return err
	}
	return w.Close()
}

func GzipReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}
