package xcrypto

import (
	"errors"
	"fmt"
	"math/rand/v2"
)

// UnpadRandom - remove Random padding
func UnpadRandom(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("padrnd: Data is empty")
	}

	if blockSize != -1 && length%blockSize != 0 {
		return nil, errors.New("padrnd: Data is not block-aligned")
	}
	padLen := int(data[length-1])
	return data[:length-padLen], nil
}

// PadRandom - add Random padding
func PadRandom(data []byte, blockSize int) ([]byte, error) {
	if blockSize == -1 {
		blockSize = 16
	}
	if blockSize <= 1 || blockSize >= 256 {
		return nil, fmt.Errorf("padrnd: Invalid block size %d", blockSize)
	} else {
		padLen := blockSize - len(data)%blockSize
		padding := make([]byte, padLen)
		for i := 0; i < (padLen - 1); i++ {
			padding[i] = byte(rand.IntN(256))
		}
		padding[padLen-1] = byte(padLen)
		return append(data, padding...), nil
	}
}
