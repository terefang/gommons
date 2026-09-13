package xcompress

import (
	"bytes"
	"io"
	"os"

	"github.com/klauspost/compress/zstd"
	"github.com/terefang/gommons/pkg/util"
)

func zstdNewWriter(dst io.Writer, level zstd.EncoderLevel) (*zstd.Encoder, error) {
	// in my tests:
	// - zstd.SpeedBestCompression is much slower and not much better
	// - default concurrency is GONUMPROCS() but adding concurrency of any value
	//   doesn't consistently speed things up
	return zstd.NewWriter(dst, zstd.WithEncoderLevel(level))
}

func ZstdCompressData(d []byte, level zstd.EncoderLevel) ([]byte, error) {
	var dst bytes.Buffer
	w, err := zstdNewWriter(&dst, level)
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

func ZstdCompressDataBest(d []byte) ([]byte, error) {
	return ZstdCompressData(d, zstd.SpeedBestCompression)
}

func ZstdCompressDataDefault(d []byte) ([]byte, error) {
	return ZstdCompressData(d, zstd.SpeedDefault)
}

func ZstdDecompressData(d []byte) ([]byte, error) {
	r := bytes.NewReader(d)
	zr, err := zstd.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	return io.ReadAll(zr)
}

func ZstdCompressFileBest(dst string, src string) error {
	return ZstdCompressFile(dst, src, zstd.SpeedBestCompression)
}

func ZstdCompressFileDefault(dst string, src string) error {
	return ZstdCompressFile(dst, src, zstd.SpeedDefault)
}

func ZstdCompressFile(dst string, src string, level zstd.EncoderLevel) error {
	fSrc, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fSrc.Close()
	fDst, err := os.Create(dst)
	if err != nil {
		return err
	}
	zw, err := zstdNewWriter(fDst, level)
	if err != nil {
		return err
	}
	_, err = io.Copy(zw, fSrc)
	err2 := zw.Close()
	err3 := fDst.Close()

	err = util.GetErr(err, err2, err3)
	if err != nil {
		os.Remove(dst)
		return err
	}
	return nil
}

func ZstdReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r, err := zstd.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}
