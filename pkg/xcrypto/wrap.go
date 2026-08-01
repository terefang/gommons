package xcrypto

import (
	"bytes"
	"crypto/cipher"
	"crypto/sha1"
	"errors"
)

// WrapWithCipher encrypts the provided key data (cek) with the given AES cipher
// (and corresponding key).
func WrapWithCipher(block cipher.Block, cek []byte) ([]byte, error) {
	_clen := block.BlockSize()
	// wrap in pkcs7padding to ensure blocksize
	cek, _ = PadPkcs7(cek, _clen)
	_IV := GenerateSalt(_clen)
	_cyp := cipher.NewCTR(block, _IV)

	_out := make([]byte, 0)
	_buf := bytes.NewBuffer(_out)
	_buf.WriteByte(byte(_clen))
	_buf.Write(_IV)

	_tmp := make([]byte, _clen)
	for i := 0; i < len(cek); i += _clen {
		_cyp.XORKeyStream(_tmp, cek[i:i+_clen])
		_buf.Write(_tmp)
	}
	// wrap again in pkcs7 to adjust for iv
	return PadPkcs7(_buf.Bytes(), _clen)
}

// UnwrapWithCipher decrypts the provided cipher text with the given cipher
// (and corresponding key).
func UnwrapWithCipher(block cipher.Block, cipherText []byte) ([]byte, error) {
	_clen := block.BlockSize()
	//unwrap from pkcs7
	cek, _ := UnpadPkcs7(cipherText, _clen)
	if int(cek[0]) != _clen {
		return nil, errors.New("iv len != block len")
	}
	_cyp := cipher.NewCTR(block, cek[1:1+_clen])

	_out := make([]byte, 0)
	_buf := bytes.NewBuffer(_out)

	_tmp := make([]byte, _clen)
	for i := 1 + _clen; i < len(cek); i += _clen {
		_cyp.XORKeyStream(_tmp, cek[i:i+_clen])
		_buf.Write(_tmp)
	}

	// unwrap from initial pkcs7 padding
	// will fail if cipher/key are incorrect
	return UnpadPkcs7(_buf.Bytes(), _clen)
}

func WrapWithSecretKey(_key []byte, plainText []byte) ([]byte, error) {
	_key = Kdf(sha1.New, 16, _key)
	return WrapWithSecretKeyCipher( /*aes.NewCipher*/ NewKdfSha1E, _key, plainText)
}

func WrapWithPassphrase(_pass string, plainText []byte) ([]byte, error) {
	_key := KdfI(sha1.New, 16, 1024, []byte(_pass))
	return WrapWithSecretKeyCipher( /*aes.NewCipher*/ NewKdfSha1E, _key, plainText)
}

func WrapWithSecretKeyCipher(_cypher func([]byte) (cipher.Block, error), _key []byte, plainText []byte) ([]byte, error) {
	_cyp, _err := _cypher(_key)
	if _err != nil {
		return nil, _err
	}
	return WrapWithCipher(_cyp, plainText)
}

func UnwrapWithPassphrase(_pass string, plainText []byte) ([]byte, error) {
	_key := KdfI(sha1.New, 16, 1024, []byte(_pass))
	return UnwrapWithSecretKeyCipher( /*aes.NewCipher*/ NewKdfSha1D, _key, plainText)
}

func UnwrapWithSecretKey(_key []byte, cipherText []byte) ([]byte, error) {
	_key = Kdf(sha1.New, 16, _key)
	return UnwrapWithSecretKeyCipher( /*aes.NewCipher*/ NewKdfSha1D, _key, cipherText)
}

func UnwrapWithSecretKeyCipher(_cypher func([]byte) (cipher.Block, error), _key []byte, cipherText []byte) ([]byte, error) {
	_cyp, _err := _cypher(_key)
	if _err != nil {
		return nil, _err
	}
	return UnwrapWithCipher(_cyp, cipherText)
}
