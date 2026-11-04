package md2

import (
	"hash"
)

const (
	Size      = 16
	BlockSize = 16
)

// piTable is the 256-byte permutation table constructed from the digits of pi (RFC 1319).
var piTable = [256]byte{
	41, 46, 67, 201, 162, 216, 124, 1, 61, 54, 84, 161, 236, 240, 6, 19,
	98, 167, 5, 243, 192, 199, 115, 140, 152, 147, 43, 217, 188, 76, 130, 202,
	30, 155, 87, 60, 253, 212, 224, 22, 103, 66, 111, 24, 138, 23, 229, 18,
	190, 78, 196, 214, 218, 158, 222, 73, 160, 251, 245, 142, 187, 47, 238, 122,
	169, 104, 121, 145, 21, 178, 7, 63, 148, 194, 16, 137, 11, 34, 95, 33,
	128, 127, 93, 154, 90, 144, 50, 39, 53, 62, 204, 231, 191, 247, 151, 3,
	255, 25, 48, 179, 72, 165, 181, 209, 215, 94, 146, 42, 172, 86, 170, 198,
	79, 184, 56, 210, 150, 164, 125, 182, 118, 252, 107, 226, 156, 116, 4, 241,
	69, 157, 112, 89, 100, 113, 135, 32, 134, 91, 207, 101, 230, 45, 168, 2,
	27, 96, 37, 173, 174, 176, 185, 246, 28, 70, 97, 105, 52, 64, 126, 15,
	85, 71, 163, 35, 221, 81, 175, 58, 195, 92, 249, 206, 186, 197, 234, 38,
	44, 83, 13, 110, 133, 40, 132, 9, 211, 223, 205, 244, 65, 129, 77, 82,
	106, 220, 55, 200, 108, 193, 171, 250, 36, 225, 123, 8, 12, 189, 177, 74,
	120, 136, 149, 139, 227, 99, 232, 109, 233, 203, 213, 254, 59, 0, 29, 57,
	242, 239, 183, 14, 102, 88, 208, 228, 166, 119, 114, 248, 235, 117, 75, 10,
	49, 68, 80, 180, 143, 237, 31, 26, 219, 153, 141, 51, 159, 17, 131, 20,
}

type digest struct {
	data []byte
}

// New returns a new hash.Hash computing the MD2 checksum.
func New() hash.Hash {
	return &digest{}
}

func (d *digest) Write(p []byte) (int, error) {
	d.data = append(d.data, p...)
	return len(p), nil
}

func (d *digest) Sum(b []byte) []byte {
	// Compute the hash on a copy/processed state to avoid changing running state
	h := computeHash(d.data)
	return append(b, h[:]...)
}

func (d *digest) Reset() {
	d.data = d.data[:0]
}

func (d *digest) Size() int {
	return Size
}

func (d *digest) BlockSize() int {
	return BlockSize
}

func computeHash(data []byte) [16]byte {
	// Step 1: Append Padding Bytes (PKCS#7-like padding to multiple of 16)
	padLen := 16 - (len(data) % 16)
	padded := make([]byte, len(data)+padLen)
	copy(padded, data)
	for i := len(data); i < len(padded); i++ {
		padded[i] = byte(padLen)
	}

	// Step 2: Append Checksum
	var checksum [16]byte
	var l byte = 0
	numBlocks := len(padded) / 16
	for i := 0; i < numBlocks; i++ {
		for j := 0; j < 16; j++ {
			c := padded[i*16+j]
			checksum[j] = piTable[c^l]
			l = checksum[j]
		}
	}
	padded = append(padded, checksum[:]...)

	// Step 3 & 4: Initialize MD Buffer and Process Message in 16-Byte Blocks
	var x [48]byte
	numBlocks = len(padded) / 16
	for i := 0; i < numBlocks; i++ {
		block := padded[i*16 : (i+1)*16]
		for j := 0; j < 16; j++ {
			x[16+j] = block[j]
			x[32+j] = x[16+j] ^ x[j]
		}
		var t byte = 0
		// 18 rounds of scrambling
		for j := 0; j < 18; j++ {
			for k := 0; k < 48; k++ {
				t = x[k] ^ piTable[t]
				x[k] = t
			}
			t = byte((int(t) + j) & 0xff)
		}
	}

	// Step 5: Output (first 16 bytes of X buffer)
	var digest [16]byte
	copy(digest[:], x[:16])
	return digest
}
