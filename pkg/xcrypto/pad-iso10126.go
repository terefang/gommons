package xcrypto

import (
	"bytes"
	"errors"
	"fmt"
)

// UnpadIso10126 - remove Iso10126 padding
func UnpadIso10126(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("Iso10126: Data is empty")
	}

	if blockSize != -1 && length%blockSize != 0 {
		return nil, errors.New("Iso10126: Data is not block-aligned")
	}
	padLen := int(data[length-1])
	ref := make([]byte, padLen)
	for i := padLen - 1; i >= 0; i-- {
		ref[i] = byte(i + 1)
	}
	if (blockSize != -1 && padLen > blockSize) || padLen == 0 || !bytes.HasSuffix(data, ref) {
		return nil, errors.New("Iso10126: Invalid padding")
	}
	return data[:length-padLen], nil
}

// PadIso10126 - add Iso10126 padding
func PadIso10126(data []byte, blockSize int) ([]byte, error) {
	if blockSize == -1 {
		blockSize = 16
	}
	if blockSize <= 1 || blockSize >= 256 {
		return nil, fmt.Errorf("Iso10126: Invalid block size %d", blockSize)
	} else {
		padLen := blockSize - len(data)%blockSize
		padding := make([]byte, padLen)
		for i := padLen - 1; i >= 0; i-- {
			padding[i] = byte(i + 1)
		}
		return append(data, padding...), nil
	}
}
