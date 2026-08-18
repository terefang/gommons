package xcrypt

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/go-crypt/crypt/algorithm/bcrypt"
)

func GenerateBcrypt(given string) string {
	_cry, _ := bcrypt.New(bcrypt.WithCost(14))
	_dgst, _ := _cry.Hash(given)
	return _dgst.String()
}

func GenerateBcryptWithSalt(given, salt string) string {
	_cry, _ := bcrypt.New(bcrypt.WithCost(14))
	_dgst, _ := _cry.HashWithSalt(given, []byte(salt))
	return _dgst.String()
}

func GenerateBcryptSha256(given string) string {
	_cry, _ := bcrypt.NewSHA256(bcrypt.WithCost(14))
	_dgst, _ := _cry.Hash(given)
	return _dgst.String()
}

func GenerateBcryptSha256WithSalt(given, salt string) string {
	_cry, _ := bcrypt.NewSHA256(bcrypt.WithCost(14))
	_dgst, _ := _cry.HashWithSalt(given, []byte(salt))
	return _dgst.String()
}

func GenerateBcrypt2Z(given string) string {
	given = Prehash(given)
	_cry, _ := bcrypt.New(bcrypt.WithCost(14))
	_dgst, _ := _cry.Hash(given)
	_str := _dgst.String()
	n := 1
	for _str[n] != '$' {
		n++
	}
	return "$2z" + _str[n:]
}

func Prehash(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	v := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return v
}

func VerifyBcrypt2Z(given, encoded string) (bool, error) {
	given = Prehash(given)
	_encoded := []byte(encoded)
	_encoded[2] = 'b'
	_dgst, _err := bcrypt.Decode(string(_encoded))
	if _err != nil {
		return false, _err
	}
	return _dgst.MatchAdvanced(given)
}
