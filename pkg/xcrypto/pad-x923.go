package xcrypto

import (
	"bytes"
	"errors"
	"fmt"
)

// UnpadX923 - remove ansi x923 padding
func UnpadX923(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("x923: Data is empty")
	}

	if blockSize != -1 && length%blockSize != 0 {
		return nil, errors.New("x923: Data is not block-aligned")
	}
	padLen := int(data[length-1])
	ref := make([]byte, padLen)
	ref[padLen-1] = byte(padLen)
	if (blockSize != -1 && padLen > blockSize) || padLen == 0 || !bytes.HasSuffix(data, ref) {
		return nil, errors.New("x923: Invalid padding")
	}
	return data[:length-padLen], nil
}

// PadX923 - add ansi x923 padding
func PadX923(data []byte, blockSize int) ([]byte, error) {
	if blockSize == -1 {
		blockSize = 16
	}
	if blockSize <= 1 || blockSize >= 256 {
		return nil, fmt.Errorf("x923: Invalid block size %d", blockSize)
	} else {
		padLen := blockSize - len(data)%blockSize
		padding := make([]byte, padLen)
		padding[padLen-1] = byte(padLen)
		return append(data, padding...), nil
	}
}
