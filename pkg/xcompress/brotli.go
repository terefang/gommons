package xcompress

import (
	"bytes"
	"io"
	"os"

	"github.com/andybalholm/brotli"
)

func BrCompressData(d []byte, level int) ([]byte, error) {
	var dst bytes.Buffer
	w := brotli.NewWriterLevel(&dst, level)
	_, err := w.Write(d)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return dst.Bytes(), nil
}

func BrCompressDataBest(d []byte) ([]byte, error) {
	return BrCompressData(d, brotli.BestCompression)
}

func BrCompressDataDefault(d []byte) ([]byte, error) {
	return BrCompressData(d, brotli.DefaultCompression)
}

func BrDecompressData(d []byte) ([]byte, error) {
	r := bytes.NewReader(d)
	zr := brotli.NewReader(r)
	return io.ReadAll(zr)
}

func BrCompressFile(dstPath string, path string, level int) error {
	d, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	d2, err := BrCompressData(d, level)
	if err != nil {
		return err
	}
	return os.WriteFile(dstPath, d2, 0644)
}

func BrCompressFileDefault(dstPath string, path string) error {
	return BrCompressFile(dstPath, path, brotli.DefaultCompression)
}

func BrCompressFileBest(dstPath string, path string) error {
	return BrCompressFile(dstPath, path, brotli.BestCompression)
}

func BrReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := brotli.NewReader(f)
	return io.ReadAll(r)
}
