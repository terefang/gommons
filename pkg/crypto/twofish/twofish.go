package twofish

import (
	"crypto/cipher"
	"encoding/binary"
	"errors"
)

// BlockSize is the constant block size in bytes for Twofish (128 bits / 16 bytes).
const BlockSize = 16

type myCipher struct {
	kLen  uint32
	lKey  [40]uint32
	sKey  [4]uint32
	qTab  [2][256]byte
	mTab  [4][256]uint32
	mkTab [4][256]uint32
}

// NewCipher creates and returns a new cipher.Block implementing the Twofish cipher.
// The key argument must be 16, 24, or 32 bytes (128, 192, or 256 bits).
func NewCipher(key []byte) (cipher.Block, error) {
	kLen := len(key)
	if kLen != 16 && kLen != 24 && kLen != 32 {
		return nil, errors.New("twofish: Incorrect key size (must be 16, 24, or 32 bytes)")
	}

	c := &myCipher{}
	c.initKey(key)
	return c, nil
}

// BlockSize returns the cipher's block size.
func (c *myCipher) BlockSize() int {
	return BlockSize
}

// Encrypt encrypts the 16-byte block of src into dst.
func (c *myCipher) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic("crypto/cipher: input data too short")
	}
	if len(dst) < BlockSize {
		panic("crypto/cipher: output data too short")
	}

	blk := [4]uint32{
		binary.LittleEndian.Uint32(src[0:4]) ^ c.lKey[0],
		binary.LittleEndian.Uint32(src[4:8]) ^ c.lKey[1],
		binary.LittleEndian.Uint32(src[8:12]) ^ c.lKey[2],
		binary.LittleEndian.Uint32(src[12:16]) ^ c.lKey[3],
	}

	for i := 0; i < 8; i++ {
		c.encryptRound(i, &blk)
	}

	binary.LittleEndian.PutUint32(dst[0:4], blk[2]^c.lKey[4])
	binary.LittleEndian.PutUint32(dst[4:8], blk[3]^c.lKey[5])
	binary.LittleEndian.PutUint32(dst[8:12], blk[0]^c.lKey[6])
	binary.LittleEndian.PutUint32(dst[12:16], blk[1]^c.lKey[7])
}

// Decrypt decrypts the 16-byte block of src into dst.
func (c *myCipher) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic("crypto/cipher: input data too short")
	}
	if len(dst) < BlockSize {
		panic("crypto/cipher: output data too short")
	}

	blk := [4]uint32{
		binary.LittleEndian.Uint32(src[0:4]) ^ c.lKey[4],
		binary.LittleEndian.Uint32(src[4:8]) ^ c.lKey[5],
		binary.LittleEndian.Uint32(src[8:12]) ^ c.lKey[6],
		binary.LittleEndian.Uint32(src[12:16]) ^ c.lKey[7],
	}

	for i := 7; i >= 0; i-- {
		c.decryptRound(i, &blk)
	}

	binary.LittleEndian.PutUint32(dst[0:4], blk[2]^c.lKey[0])
	binary.LittleEndian.PutUint32(dst[4:8], blk[3]^c.lKey[1])
	binary.LittleEndian.PutUint32(dst[8:12], blk[0]^c.lKey[2])
	binary.LittleEndian.PutUint32(dst[12:16], blk[1]^c.lKey[3])
}

func (c *myCipher) encryptRound(i int, blk *[4]uint32) {
	t1 := c.g1Func(blk[1])
	t0 := c.g0Func(blk[0])
	blk[2] = rotr32(blk[2]^(t0+t1+c.lKey[4*i+8]), 1)
	blk[3] = rotl32(blk[3], 1) ^ (t0 + 2*t1 + c.lKey[4*i+9])

	t1 = c.g1Func(blk[3])
	t0 = c.g0Func(blk[2])
	blk[0] = rotr32(blk[0]^(t0+t1+c.lKey[4*i+10]), 1)
	blk[1] = rotl32(blk[1], 1) ^ (t0 + 2*t1 + c.lKey[4*i+11])
}

func (c *myCipher) decryptRound(i int, blk *[4]uint32) {
	t1 := c.g1Func(blk[1])
	t0 := c.g0Func(blk[0])
	blk[2] = rotl32(blk[2], 1) ^ (t0 + t1 + c.lKey[4*i+10])
	blk[3] = rotr32(blk[3]^(t0+2*t1+c.lKey[4*i+11]), 1)

	t1 = c.g1Func(blk[3])
	t0 = c.g0Func(blk[2])
	blk[0] = rotl32(blk[0], 1) ^ (t0 + t1 + c.lKey[4*i+8])
	blk[1] = rotr32(blk[1]^(t0+2*t1+c.lKey[4*i+9]), 1)
}

func (c *myCipher) g0Func(x uint32) uint32 {
	return c.mkTab[0][byteVal(x, 0)] ^
		c.mkTab[1][byteVal(x, 1)] ^
		c.mkTab[2][byteVal(x, 2)] ^
		c.mkTab[3][byteVal(x, 3)]
}

func (c *myCipher) g1Func(x uint32) uint32 {
	return c.mkTab[0][byteVal(x, 3)] ^
		c.mkTab[1][byteVal(x, 0)] ^
		c.mkTab[2][byteVal(x, 1)] ^
		c.mkTab[3][byteVal(x, 2)]
}

func byteVal(x uint32, n uint) byte {
	return byte(x >> (8 * n))
}

const gM = 0x0169

var tab5b = [4]byte{0, gM >> 2, gM >> 1, (gM >> 1) ^ (gM >> 2)}
var tabEf = [4]byte{0, (gM >> 1) ^ (gM >> 2), gM >> 1, gM >> 2}

func ffm5b(x byte) byte {
	return x ^ (x >> 2) ^ tab5b[x&3]
}

func ffmEf(x byte) byte {
	return x ^ (x >> 1) ^ (x >> 2) ^ tabEf[x&3]
}

var ror4 = [16]byte{0, 8, 1, 9, 2, 10, 3, 11, 4, 12, 5, 13, 6, 14, 7, 15}
var ashx = [16]byte{0, 9, 2, 11, 4, 13, 6, 15, 8, 1, 10, 3, 12, 5, 14, 7}

var qt0 = [2][16]byte{
	{8, 1, 7, 13, 6, 15, 3, 2, 0, 11, 5, 9, 14, 12, 10, 4},
	{2, 8, 11, 13, 15, 7, 6, 14, 3, 1, 9, 4, 0, 10, 12, 5},
}

var qt1 = [2][16]byte{
	{14, 12, 11, 8, 1, 2, 3, 5, 15, 4, 10, 6, 7, 0, 9, 13},
	{1, 14, 2, 11, 4, 12, 3, 7, 6, 13, 10, 5, 15, 9, 0, 8},
}

var qt2 = [2][16]byte{
	{11, 10, 5, 14, 6, 13, 9, 0, 12, 8, 15, 3, 2, 4, 7, 1},
	{4, 12, 7, 5, 1, 6, 9, 10, 0, 14, 13, 8, 2, 11, 3, 15},
}

var qt3 = [2][16]byte{
	{13, 7, 15, 4, 1, 2, 6, 14, 9, 11, 3, 0, 8, 5, 12, 10},
	{11, 9, 5, 1, 12, 3, 13, 14, 6, 4, 7, 15, 2, 0, 8, 10},
}

func qp(n uint32, x byte) byte {
	a0 := x >> 4
	b0 := x & 15
	a1 := a0 ^ b0
	b1 := ror4[b0] ^ ashx[a0]
	a2 := qt0[n][a1]
	b2 := qt1[n][b1]
	a3 := a2 ^ b2
	b3 := ror4[b2] ^ ashx[a2]
	a4 := qt2[n][a3]
	b4 := qt3[n][b3]
	return (b4 << 4) | a4
}

func (c *myCipher) genQTab() {
	for i := 0; i < 256; i++ {
		c.qTab[0][i] = qp(0, byte(i))
		c.qTab[1][i] = qp(1, byte(i))
	}
}

func (c *myCipher) q(n int, x byte) byte {
	return c.qTab[n][x]
}

func (c *myCipher) genMTab() {
	for i := 0; i < 256; i++ {
		f01 := uint32(c.q(1, byte(i)))
		f5b := uint32(ffm5b(byte(f01)))
		fef := uint32(ffmEf(byte(f01)))
		c.mTab[0][i] = f01 + (f5b << 8) + (fef << 16) + (fef << 24)
		c.mTab[2][i] = f5b + (fef << 8) + (f01 << 16) + (fef << 24)

		f01 = uint32(c.q(0, byte(i)))
		f5b = uint32(ffm5b(byte(f01)))
		fef = uint32(ffmEf(byte(f01)))
		c.mTab[1][i] = fef + (fef << 8) + (f5b << 16) + (f01 << 24)
		c.mTab[3][i] = f5b + (f01 << 8) + (fef << 16) + (f5b << 24)
	}
}

func (c *myCipher) mds(n int, x byte) uint32 {
	return c.mTab[n][x]
}

func (c *myCipher) hFunc(x uint32, key []uint32) uint32 {
	b0 := byteVal(x, 0)
	b1 := byteVal(x, 1)
	b2 := byteVal(x, 2)
	b3 := byteVal(x, 3)

	switch c.kLen {
	case 4:
		b0 = c.q(1, b0) ^ byteVal(key[3], 0)
		b1 = c.q(0, b1) ^ byteVal(key[3], 1)
		b2 = c.q(0, b2) ^ byteVal(key[3], 2)
		b3 = c.q(1, b3) ^ byteVal(key[3], 3)
		fallthrough
	case 3:
		b0 = c.q(1, b0) ^ byteVal(key[2], 0)
		b1 = c.q(1, b1) ^ byteVal(key[2], 1)
		b2 = c.q(0, b2) ^ byteVal(key[2], 2)
		b3 = c.q(0, b3) ^ byteVal(key[2], 3)
		fallthrough
	case 2:
		b0 = c.q(0, c.q(0, b0)^byteVal(key[1], 0)) ^ byteVal(key[0], 0)
		b1 = c.q(0, c.q(1, b1)^byteVal(key[1], 1)) ^ byteVal(key[0], 1)
		b2 = c.q(1, c.q(0, b2)^byteVal(key[1], 2)) ^ byteVal(key[0], 2)
		b3 = c.q(1, c.q(1, b3)^byteVal(key[1], 3)) ^ byteVal(key[0], 3)
	}

	return c.mds(0, b0) ^ c.mds(1, b1) ^ c.mds(2, b2) ^ c.mds(3, b3)
}

const gMod = 0x0000014d

func mdsRem(p0, p1 uint32) uint32 {
	for i := 0; i < 8; i++ {
		t := p1 >> 24
		p1 = (p1 << 8) | (p0 >> 24)
		p0 <<= 8

		u := t << 1
		if (t & 0x80) != 0 {
			u ^= gMod
		}

		p1 ^= t ^ (u << 16)
		u ^= (t >> 1)
		if (t & 0x01) != 0 {
			u ^= gMod >> 1
		}
		p1 ^= (u << 24) | (u << 8)
	}
	return p1
}

func (c *myCipher) genMkTab(key []uint32) {
	switch c.kLen {
	case 2:
		for i := 0; i < 256; i++ {
			by := byte(i)
			c.mkTab[0][i] = c.mds(0, c.q(0, c.q(0, by)^byteVal(key[1], 0))^byteVal(key[0], 0))
			c.mkTab[1][i] = c.mds(1, c.q(0, c.q(1, by)^byteVal(key[1], 1))^byteVal(key[0], 1))
			c.mkTab[2][i] = c.mds(2, c.q(1, c.q(0, by)^byteVal(key[1], 2))^byteVal(key[0], 2))
			c.mkTab[3][i] = c.mds(3, c.q(1, c.q(1, by)^byteVal(key[1], 3))^byteVal(key[0], 3))
		}
	case 3:
		for i := 0; i < 256; i++ {
			by := byte(i)
			c.mkTab[0][i] = c.mds(0, c.q(0, c.q(0, c.q(1, by)^byteVal(key[2], 0))^byteVal(key[1], 0))^byteVal(key[0], 0))
			c.mkTab[1][i] = c.mds(1, c.q(0, c.q(1, c.q(1, by)^byteVal(key[2], 1))^byteVal(key[1], 1))^byteVal(key[0], 1))
			c.mkTab[2][i] = c.mds(2, c.q(1, c.q(0, c.q(0, by)^byteVal(key[2], 2))^byteVal(key[1], 2))^byteVal(key[0], 2))
			c.mkTab[3][i] = c.mds(3, c.q(1, c.q(1, c.q(0, by)^byteVal(key[2], 3))^byteVal(key[1], 3))^byteVal(key[0], 3))
		}
	case 4:
		for i := 0; i < 256; i++ {
			by := byte(i)
			c.mkTab[0][i] = c.mds(0, c.q(0, c.q(0, c.q(1, c.q(1, by)^byteVal(key[3], 0))^byteVal(key[2], 0))^byteVal(key[1], 0))^byteVal(key[0], 0))
			c.mkTab[1][i] = c.mds(1, c.q(0, c.q(1, c.q(1, c.q(0, by)^byteVal(key[3], 1))^byteVal(key[2], 1))^byteVal(key[1], 1))^byteVal(key[0], 1))
			c.mkTab[2][i] = c.mds(2, c.q(1, c.q(0, c.q(0, c.q(0, by)^byteVal(key[3], 2))^byteVal(key[2], 2))^byteVal(key[1], 2))^byteVal(key[0], 2))
			c.mkTab[3][i] = c.mds(3, c.q(1, c.q(1, c.q(0, c.q(1, by)^byteVal(key[3], 3))^byteVal(key[2], 3))^byteVal(key[1], 3))^byteVal(key[0], 3))
		}
	}
}

func (c *myCipher) initKey(rawKey []byte) {
	c.genQTab()
	c.genMTab()

	keyLenWords := uint32(len(rawKey) / 4)
	c.kLen = keyLenWords / 2 // 2, 3, or 4

	inKey := make([]uint32, len(rawKey)/4)
	for i := range inKey {
		inKey[i] = binary.LittleEndian.Uint32(rawKey[i*4 : (i+1)*4])
	}

	meKey := make([]uint32, c.kLen)
	moKey := make([]uint32, c.kLen)
	sKey := make([]uint32, c.kLen)

	for i := uint32(0); i < c.kLen; i++ {
		a := inKey[i*2]
		meKey[i] = a
		b := inKey[i*2+1]
		moKey[i] = b
		sKey[c.kLen-i-1] = mdsRem(a, b)
	}
	copy(c.sKey[:], sKey)

	for i := 0; i < 40; i += 2 {
		a := uint32(0x01010101 * i)
		b := a + 0x01010101
		a = c.hFunc(a, meKey)
		b = rotl32(c.hFunc(b, moKey), 8)
		c.lKey[i] = a + b
		c.lKey[i+1] = rotl32(a+2*b, 9)
	}

	c.genMkTab(sKey)
}

func rotl32(x uint32, n uint) uint32 {
	return (x << n) | (x >> (32 - n))
}

func rotr32(x uint32, n uint) uint32 {
	return (x >> n) | (x << (32 - n))
}
