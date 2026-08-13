package xcrypto

import (
	"errors"
	"fmt"
)

// UnpadZero - remove zero padding
func UnpadZero(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("zero: Data is empty")
	}

	if blockSize != -1 && length%blockSize != 0 {
		return nil, errors.New("zero: Data is not block-aligned")
	}

	padLen := 0
	for (data[(length-1)-padLen] == 0) && (padLen < blockSize) {
		padLen++
	}
	return data[:length-padLen], nil
}

// PadZero - add zero padding
func PadZero(data []byte, blockSize int) ([]byte, error) {
	if blockSize == -1 {
		blockSize = 16
	}
	if blockSize <= 1 || blockSize >= 256 {
		return nil, fmt.Errorf("zero: Invalid block size %d", blockSize)
	} else {
		padLen := blockSize - len(data)%blockSize
		padding := make([]byte, padLen)
		return append(data, padding...), nil
	}
}
