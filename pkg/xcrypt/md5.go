package xcrypt

import (
	"github.com/go-crypt/crypt/algorithm/md5crypt"
	"github.com/tredoe/crypt/apr1_crypt"
)

func GenerateCrypt1(_plain string) (string, error) {
	_cry, _err := md5crypt.New()
	if _err != nil {
		return "", _err
	}
	_dgst, _err := _cry.Hash(_plain)
	if _err != nil {
		return "", _err
	}
	return _dgst.String(), nil
}

func GenerateCrypt1WithSalt(_plain string, _salt string) (string, error) {
	_cry, _err := md5crypt.New()
	if _err != nil {
		return "", _err
	}
	_dgst, _err := _cry.HashWithSalt(_plain, []byte(_salt))
	if _err != nil {
		return "", _err
	}
	return _dgst.String(), nil
}

func VerifyCrypt1(given, encrypted string) (bool, error) {
	_dgst, err := md5crypt.Decode(encrypted)
	if err != nil {
		return false, err
	}
	return _dgst.MatchAdvanced(given)
}

func VerifyApr1(given, encrypted string) bool {
	err := apr1_crypt.New().Verify(encrypted, []byte(given))
	if err != nil {
		return false
	}
	return true
}
