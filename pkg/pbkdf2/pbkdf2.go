package pbkdf2

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"hash"
)

// GenerateKey generates a derived key from a password using the PBKDF2 algorithm. The
// inputs include salt bytes, the iteration count, desired key length, and a
// constructor for a hashing function.  For example, for a 32-byte key using
// SHA-256:
//
//	key := GenerateKey([]byte("trustNo1"), salt, 10000, 32, sha256.New)
func GenerateKey(password, salt []byte, iterCount, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hLen := prf.Size()
	numBlocks := keyLen / hLen
	// Get an extra block if keyLen is not an even number of hLen blocks.
	if keyLen%hLen > 0 {
		numBlocks++
	}

	Ti := make([]byte, hLen)
	Uj := make([]byte, hLen)
	dk := make([]byte, 0, hLen*numBlocks)
	buf := make([]byte, 4)

	for i := uint32(1); i <= uint32(numBlocks); i++ {
		// Initialize Uj for j == 1 from salt and block index.
		// Initialize Ti = U1.
		binary.BigEndian.PutUint32(buf, i)
		prf.Reset()
		prf.Write(salt)
		prf.Write(buf)
		Uj = Uj[:0]
		Uj = prf.Sum(Uj)

		// Ti = U1 ^ U2 ^ ... ^ Ux
		copy(Ti, Uj)
		for j := 2; j <= iterCount; j++ {
			prf.Reset()
			prf.Write(Uj)
			Uj = Uj[:0]
			Uj = prf.Sum(Uj)
			for k := range Uj {
				Ti[k] ^= Uj[k]
			}
		}

		// DK = concat(T1, T2, ... Tn)
		dk = append(dk, Ti...)
	}

	return dk[0:keyLen]
}

func GenerateKeySecure(password, salt []byte, iterCount, keyLen int) []byte {
	return GenerateKey(password, salt, iterCount, keyLen, sha512.New)
}
