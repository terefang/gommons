package xcompress

import (
	"bytes"
	"io"

	"github.com/therootcompany/xz"
)

func XzDecompressData(d []byte) ([]byte, error) {
	r := bytes.NewReader(d)
	zr, err := xz.NewReader(r, 0)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(zr)
}
