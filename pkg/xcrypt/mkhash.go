package xcrypt

import (
	"github.com/go-crypt/crypt/algorithm/argon2"
	"github.com/go-crypt/crypt/algorithm/bcrypt"
	"github.com/go-crypt/crypt/algorithm/md5crypt"
	"github.com/go-crypt/crypt/algorithm/pbkdf2"
	"github.com/go-crypt/crypt/algorithm/scrypt"
	"github.com/go-crypt/crypt/algorithm/sha1crypt"
	"github.com/go-crypt/crypt/algorithm/shacrypt"
)

func GenerateCisco1Crypt(given string) string {
	_cry, _ := md5crypt.New(md5crypt.WithSaltLength(4))
	_dgst, _ := _cry.Hash(given)
	return _dgst.String()
}

func GenerateCisco1CryptWithSalt(given, salt string) string {
	_cry, _ := md5crypt.New(md5crypt.WithSaltLength(4))
	_dgst, _ := _cry.HashWithSalt(given, []byte(salt))
	return _dgst.String()
}

func GenerateType7Crypt(given string) string {
	return "{type7}" + Type7_encrypt(given)
}

func GenerateMd5Crypt(given string) string {
	_cry, _ := md5crypt.New()
	_dgst, _ := _cry.Hash(given)
	return _dgst.String()
}

func GenerateMd5CryptWithSalt(given, salt string) string {
	_cry, _ := md5crypt.New()
	_dgst, _ := _cry.HashWithSalt(given, []byte(salt))
	return _dgst.String()
}

func GenerateApr1Crypt(given string) string {
	_cry, _ := md5crypt.New()
	_dgst, _ := _cry.Hash(given)
	return "$apr1" + _dgst.String()[2:]
}

func GenerateApr1CryptWithSalt(given, salt string) string {
	_cry, _ := md5crypt.New()
	_dgst, _ := _cry.HashWithSalt(given, []byte(salt))
	return "$apr1" + _dgst.String()[2:]
}

func GenerateSha1Crypt(given string) string {
	_cry, _ := sha1crypt.New(sha1crypt.WithIterations(1000), sha1crypt.WithSaltLength(16))
	_dgst, _ := _cry.Hash(given)
	return _dgst.String()
}

func GenerateSha1CryptWithSalt(given, salt string) string {
	_cry, _ := sha1crypt.New(sha1crypt.WithIterations(1000))
	_dgst, _ := _cry.HashWithSalt(given, []byte(salt))
	return _dgst.String()
}

func GenerateSshaCrypt(given string) string {
	_dgst, _ := SshaCrypt(given)
	return _dgst
}

func GenerateSha256Crypt(given string) string {
	_cry, _ := shacrypt.New(shacrypt.WithSHA256(), shacrypt.WithIterations(5000), shacrypt.WithSaltLength(16))
	_dgst, _ := _cry.Hash(given)
	return _dgst.String()
}

func GenerateSha256CryptWithSalt(given, salt string) string {
	_cry, _ := shacrypt.New(shacrypt.WithSHA256(), shacrypt.WithIterations(5000))
	_dgst, _ := _cry.HashWithSalt(given, []byte(salt))
	return _dgst.String()
}

func GenerateSha512Crypt(given string) string {
	_cry, _ := shacrypt.New(shacrypt.WithSHA512(), shacrypt.WithIterations(5000), shacrypt.WithSaltLength(16))
	_dgst, _ := _cry.Hash(given)
	return _dgst.String()
}

func GenerateSha512CryptWithSalt(given, salt string) string {
	_cry, _ := shacrypt.New(shacrypt.WithSHA512(), shacrypt.WithIterations(5000))
	_dgst, _ := _cry.HashWithSalt(given, []byte(salt))
	return _dgst.String()
}

func GenerateScrypt(given string) string {
	_cry, _ := scrypt.New()
	_dgst, _ := _cry.Hash(given)
	return _dgst.String()
}

func GenerateScryptWithSalt(given, salt string) string {
	_cry, _ := scrypt.New()
	_dgst, _ := _cry.HashWithSalt(given, []byte(salt))
	return _dgst.String()
}

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

func GeneratePbkdf2Sha256(given string) string {
	_cry, _ := pbkdf2.NewSHA256()
	_dgst, _ := _cry.Hash(given)
	return _dgst.String()
}

func GeneratePbkdf2Sha256WithSalt(given, salt string) string {
	_cry, _ := pbkdf2.NewSHA256()
	_dgst, _ := _cry.HashWithSalt(given, []byte(salt))
	return _dgst.String()
}

func GeneratePbkdf2Sha1(given string) string {
	_cry, _ := pbkdf2.NewSHA1()
	_dgst, _ := _cry.Hash(given)
	return _dgst.String()
}

func GeneratePbkdf2Sha1WithSalt(given, salt string) string {
	_cry, _ := pbkdf2.NewSHA1()
	_dgst, _ := _cry.HashWithSalt(given, []byte(salt))
	return _dgst.String()
}

func GenerateArgon2(given string) string {
	_cry, _ := argon2.New()
	_dgst, _ := _cry.Hash(given)
	return _dgst.String()
}

func GenerateArgon2WithSalt(given, salt string) string {
	_cry, _ := argon2.New()
	_dgst, _ := _cry.HashWithSalt(given, []byte(salt))
	return _dgst.String()
}
