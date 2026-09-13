package xcompress

import (
	"bytes"
	"compress/bzip2"
	"io"
	"os"
)

func Bzip2DecompressData(d []byte) ([]byte, error) {
	r := bytes.NewReader(d)
	zr := bzip2.NewReader(r)
	return io.ReadAll(zr)
}

func Bzip2ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bzip2.NewReader(f)
	return io.ReadAll(r)
}
