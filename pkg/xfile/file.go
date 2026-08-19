package xfile

import (
	"errors"
	"io/fs"
	"os"
)

func MakePath(path string, perm fs.FileMode) error {
	if FileExists(path) {
		return nil
	}
	return os.MkdirAll(path, perm)
}

func FileExists(filepath string) bool {
	if _fi, err := os.Stat(filepath); (!_fi.IsDir()) && (!os.IsNotExist(err)) {
		return true
	}
	return false
}

func PathExists(filepath string) bool {
	if _fi, err := os.Stat(filepath); _fi.IsDir() && (!os.IsNotExist(err)) {
		return true
	}
	return false
}

func CreatePath(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func CreatePathIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return CreatePath(dir)
	}
	return errors.New("dir already exists")
}

func IsEmpty(dir string) (bool, error) {
	es, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	return len(es) == 0, nil
}

func DeletePath(dir string) error {
	return os.RemoveAll(dir)
}

func DeletePathIfExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return errors.New("dir does not exist")
	}
	return DeletePath(dir)
}

func DeletePathIfEmpty(dir string) error {
	if b, err := IsEmpty(dir); err != nil {
		return err
	} else if !b {
		return errors.New("dir is not empty")
	}
	return DeletePath(dir)
}

func IsRegularFile(path string) bool {
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return fi.Mode().IsRegular()
}
