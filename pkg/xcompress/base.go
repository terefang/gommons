package xcompress

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"compress/zlib"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

// -----------------------------------------------
// added by Alfred Reibenschuh

func InlineFlateCompress(plain []byte) []byte {
	var _buf bytes.Buffer
	_bufwr := bufio.NewWriter(&_buf)
	_flate := zlib.NewWriter(_bufwr)
	_flate.Write(plain)
	_flate.Flush()
	_flate.Close()
	_bufwr.Flush()
	return _buf.Bytes()
}

//func InlineFlateDecompress(plain []byte) []byte {
//	var _buf bytes.Buffer
//	_buf.Write(plain)
//	_bufrd := bufio.NewReader(&_buf)
//	_flate, _ := zlib.NewReader(_bufrd)
//
//	_out := make([]byte, 0)
//	_b := make([]byte, 8192)
//	for _n, _err := _flate.Read(_b); _err == nil; {
//		_out = append(_out, _b[:_n]...)
//	}
//	_flate.Close()
//	return _out
//}

// implement io.ReadCloser over os.File wrapped with io.Reader.
// io.Closer goes to os.File, io.Reader goes to wrapping reader
type readerWrappedFile struct {
	f *os.File
	r io.Reader
}

func (rc *readerWrappedFile) Close() error {
	return rc.f.Close()
}

func (rc *readerWrappedFile) Read(p []byte) (int, error) {
	return rc.r.Read(p)
}

func wrapInReadCloser(f *os.File, r io.Reader, err error) (io.ReadCloser, error) {
	if err != nil {
		f.Close()
		return nil, err
	}
	return &readerWrappedFile{
		f: f,
		r: r,
	}, nil
}

// OpenFileMaybeCompressed opens a file that might be compressed with gzip
// or bzip2 or zstd or brotli
// TODO: could sniff file content instead of checking file extension
func OpenFileMaybeCompressed(path string) (io.ReadCloser, error) {
	ext := strings.ToLower(filepath.Ext(path))
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if ext == ".gz" {
		r, err := gzip.NewReader(f)
		return wrapInReadCloser(f, r, err)
	}
	if ext == ".bz2" {
		r := bzip2.NewReader(f)
		return wrapInReadCloser(f, r, nil)
	}
	if ext == ".zstd" {
		r, err := zstd.NewReader(f)
		return wrapInReadCloser(f, r, err)
	}
	if ext == ".br" {
		r := brotli.NewReader(f)
		return wrapInReadCloser(f, r, nil)
	}
	return f, nil
}

// ReadFileMaybeCompressed reads file. Ungzips if it's gzipped.
func ReadFileMaybeCompressed(path string) ([]byte, error) {
	r, err := OpenFileMaybeCompressed(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}
