package xcrypto

import (
	"errors"
	"hash"
)

// S2KType defines OpenPGP S2K specifier types.
type S2KType byte

const (
	SimpleS2K         S2KType = 0
	SaltedS2K         S2KType = 1
	IteratedSaltedS2K S2KType = 3
)

// S2kSimple derives a key using Simple S2K (Type 0).
func S2kSimple(passphrase []byte, keyLen int, h func() hash.Hash) ([]byte, error) {
	return expandS2kHash(passphrase, keyLen, h)
}

// S2kSalted derives a key using Salted S2K (Type 1).
func S2kSalted(passphrase, salt []byte, keyLen int, h func() hash.Hash) ([]byte, error) {
	if len(salt) != 8 {
		return nil, errors.New("s2k: salt must be exactly 8 bytes")
	}
	data := append(salt, passphrase...)
	return expandS2kHash(data, keyLen, h)
}

// S2kIteratedAndSalted derives a key using Iterated and Salted S2K (Type 3).
// count specifies the encoded 1-byte count parameter (0-255).
func S2kIteratedAndSalted(passphrase, salt []byte, countByte byte, keyLen int, h func() hash.Hash) ([]byte, error) {
	if len(salt) != 8 {
		return nil, errors.New("s2k: salt must be exactly 8 bytes")
	}

	decodedCount := S2kDecodeCount(countByte)
	combined := append(salt, passphrase...)
	if len(combined) == 0 {
		return nil, errors.New("s2k: passphrase and salt cannot both be empty")
	}

	hasher := h()
	bytesHashed := 0

	// Loop until we have fed 'decodedCount' bytes into the hash context
	for bytesHashed < decodedCount {
		needed := decodedCount - bytesHashed
		if needed >= len(combined) {
			hasher.Write(combined)
			bytesHashed += len(combined)
		} else {
			hasher.Write(combined[:needed])
			bytesHashed += needed
		}
	}

	// First chunk of derived key material
	key := hasher.Sum(nil)

	// If requested key length exceeds single hash digest size, perform OpenPGP multiple hash passes
	if len(key) < keyLen {
		out := make([]byte, 0, keyLen)
		out = append(out, key...)
		pass := 1

		for len(out) < keyLen {
			hasher.Reset()

			// OpenPGP S2K multi-pass: Prepend 'pass' number of zero bytes for pass N
			zeroes := make([]byte, pass)
			hasher.Write(zeroes)

			// Re-hash the iterated salted block
			bytesHashed = 0
			for bytesHashed < decodedCount {
				needed := decodedCount - bytesHashed
				if needed >= len(combined) {
					hasher.Write(combined)
					bytesHashed += len(combined)
				} else {
					hasher.Write(combined[:needed])
					bytesHashed += needed
				}
			}

			out = append(out, hasher.Sum(nil)...)
			pass++
		}
		key = out
	}

	return key[:keyLen], nil
}

// S2kEncodeCount converts an actual iteration count to the RFC 4880 1-byte S2K count representation.
func S2kEncodeCount(count int) byte {
	if count <= 1024 {
		return 0
	}
	if count >= 65011712 {
		return 255
	}

	// Find exponent c such that (16 + (count & 15)) << (c + 6) >= count
	var c byte
	for c = 0; c < 16; c++ {
		if (16+15)<<(c+6) >= count {
			break
		}
	}

	// Calculate 4-bit mantissa
	bits := count >> (c + 6)
	if bits > 31 {
		bits = 31
	}
	val := bits - 16

	return (c << 4) | byte(val)
}

// S2kDecodeCount decodes RFC 4880 1-byte S2K iteration count parameter into absolute byte count.
func S2kDecodeCount(c byte) int {
	return (16 + int(c&15)) << (uint32(c>>4) + 6)
}

// Helper to expand output to keyLen if hash output is smaller than target length (RFC 4880 Section 3.7.1)
func expandS2kHash(data []byte, keyLen int, h func() hash.Hash) ([]byte, error) {
	hasher := h()
	hasher.Write(data)
	digest := hasher.Sum(nil)

	if len(digest) >= keyLen {
		return digest[:keyLen], nil
	}

	out := make([]byte, 0, keyLen)
	out = append(out, digest...)
	pass := 1

	for len(out) < keyLen {
		hasher.Reset()
		// Prepend zero bytes equal to the pass index
		zeroes := make([]byte, pass)
		hasher.Write(zeroes)
		hasher.Write(data)
		out = append(out, hasher.Sum(nil)...)
		pass++
	}

	return out[:keyLen], nil
}
