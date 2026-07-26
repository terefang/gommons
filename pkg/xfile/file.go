package xfile

import "os"

func MakePath(path string) error {
	if FileExists(path) {
		return nil
	}
	return os.MkdirAll(path, 0700)
}

func FileExists(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}
	return true
}
