package xcrypto

import (
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"hash"
)

func Kdf(_hash func() hash.Hash, l int, s []byte) []byte {
	_h := _hash()
	_bl := _h.Size()
	_bc := l / _bl
	_rl := (_bc + 1) * _bl
	_ret := make([]byte, _rl)

	for i := 0; i < l; i += _bl {
		_h.Reset()
		_h.Write(_ret)
		_h.Write(s)
		_md := _h.Sum(nil)
		copy(_ret[i:i+_bl], _md)
	}
	return _ret[:l]
}

func KdfI(_hash func() hash.Hash, l int, _itr int, s []byte) []byte {
	_h := _hash()
	_bl := _h.Size()
	_bc := l / _bl
	_rl := (_bc + 1) * _bl
	_ret := make([]byte, _rl)
	for j := 0; j < _itr; j++ {
		_m := hmac.New(_hash, s)
		for i := 0; i < l; i += _bl {
			_m.Reset()
			_m.Write(_ret)
			_md := _m.Sum(nil)
			copy(_ret[i:i+_bl], _md)
		}
	}
	return _ret[:l]
}

func GenerateSalt(l int) []byte {
	_salt := make([]byte, l)
	rand.Reader.Read(_salt)
	return _salt
}

func GenerateSecretKeyWithSalt(l int, _salt []byte) []byte {
	return GenerateSecretKeyWithAlgoSalt(sha1.New, l, _salt)
}

func GenerateSecretKeyWithAlgoSalt(h func() hash.Hash, l int, _salt []byte) []byte {
	if _salt == nil {
		_salt = GenerateSalt(64)
	}
	return Kdf(h, l, _salt)
}

type kdfBlock struct {
	state    []byte
	blockLen int
	mode     bool
	hash     func() hash.Hash
	ctr      int
}

func (k kdfBlock) NewCTR(_iv []byte) cipher.Stream {
	_tmp := append(k.state, _iv...)
	k.state = Kdf(k.hash, k.blockLen, _tmp)
	return &k
}

func (k *kdfBlock) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for i := 0; i < len(src); i++ {
		if k.ctr >= k.blockLen {
			k.Refill()
		}
		dst[i] = src[i] ^ k.state[k.ctr]
		k.ctr++
	}
}

func NewKdfSha1E(_key []byte) (cipher.Block, error) {
	return NewKdfCipher(_key, sha1.New, true), nil
}
func NewKdfSha1D(_key []byte) (cipher.Block, error) {
	return NewKdfCipher(_key, sha1.New, false), nil
}

func NewKdfCipher(_key []byte, _h func() hash.Hash, encrypt bool) cipher.Block {
	_kb := &kdfBlock{}
	_kb.mode = encrypt
	_kb.blockLen = _h().BlockSize()
	_kb.hash = _h
	_kb.state = Kdf(_kb.hash, _kb.blockLen, _key)
	return _kb
}

func (k *kdfBlock) Refill() {
	k.state = Kdf(k.hash, k.blockLen, k.state)
	k.ctr = 0
}

func (k kdfBlock) BlockSize() int {
	return k.blockLen
}

func (k kdfBlock) Encrypt(dst, src []byte) {
	if k.mode == false {
		panic("kdfcipher: mode is not encryption")
	}
	for i := 0; i < k.blockLen; i++ {
		dst[i] = src[i] ^ k.state[i]
	}
	k.Refill()
}

func (k kdfBlock) Decrypt(dst, src []byte) {
	if k.mode == true {
		panic("kdfcipher: mode is not decryption")
	}
	for i := 0; i < k.blockLen; i++ {
		dst[i] = src[i] ^ k.state[i]
	}
	k.Refill()
}
