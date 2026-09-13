package xcrypto

import (
	"crypto/hmac"
	"hash"
)

// HMAC returns a keyed hash digest of the data
func HMAC(key []byte, data []byte, _hash func() hash.Hash) []byte {
	mac := hmac.New(_hash, key)
	mac.Write(data)
	return mac.Sum(nil)
}
