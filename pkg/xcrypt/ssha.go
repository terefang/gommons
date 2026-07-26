package xcrypt

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

func SshaCrypt(plain string) (string, error) {
	_enc := SSHAEncoder{}
	_crypted, _err := _enc.Encode([]byte(plain))
	if _err != nil {
		return "", _err
	}
	return string(_crypted), nil
}

func SshaCryptMatch(given string, stored string) bool {
	_enc := SSHAEncoder{}
	return _enc.Matches([]byte(stored), []byte(given))
}

type SSHAEncoder struct {
}

// Encode encodes the []byte of raw password
func (enc SSHAEncoder) Encode(rawPassPhrase []byte) ([]byte, error) {
	hash := makeSSHAHash(rawPassPhrase, makeSalt())
	b64 := base64.StdEncoding.EncodeToString(hash)
	return []byte(fmt.Sprintf("{SSHA}%s", b64)), nil
}

// Matches matches the encoded password and the raw password
func (enc SSHAEncoder) Matches(encodedPassPhrase, rawPassPhrase []byte) bool {
	//strip the {SSHA}
	eppS := string(encodedPassPhrase)[6:]
	hash, err := base64.StdEncoding.DecodeString(eppS)
	if err != nil {
		return false
	}
	salt := hash[len(hash)-4:]

	sha := sha1.New()
	sha.Write(rawPassPhrase)
	sha.Write(salt)
	sum := sha.Sum(nil)

	if bytes.Compare(sum, hash[:len(hash)-4]) != 0 {
		return false
	}
	return true
}

// makeSalt make a 4 byte array containing random bytes.
func makeSalt() []byte {
	sbytes := make([]byte, 4)
	rand.Read(sbytes)
	return sbytes
}

// makeSSHAHash make hasing using SHA-1 with salt. This is not the final output though. You need to append {SSHA} string with base64 of this hash.
func makeSSHAHash(passphrase, salt []byte) []byte {
	sha := sha1.New()
	sha.Write(passphrase)
	sha.Write(salt)

	h := sha.Sum(nil)
	return append(h, salt...)
}

// ---------------------------------------------------------------------------

func ShaCrypt(plain string) (string, error) {
	_enc := SHAEncoder{}
	_crypted, _err := _enc.Encode([]byte(plain))
	if _err != nil {
		return "", _err
	}
	return string(_crypted), nil
}

func ShaCryptMatch(given string, stored string) bool {
	_enc := SHAEncoder{}
	return _enc.Matches([]byte(stored), []byte(given))
}

type SHAEncoder struct {
}

// Encode encodes the []byte of raw password
func (enc SHAEncoder) Encode(rawPassPhrase []byte) ([]byte, error) {
	hash := makeSHAHash(rawPassPhrase)
	b64 := base64.StdEncoding.EncodeToString(hash)
	return []byte(fmt.Sprintf("{SHA}%s", b64)), nil
}

// Matches matches the encoded password and the raw password
func (enc SHAEncoder) Matches(encodedPassPhrase, rawPassPhrase []byte) bool {
	//strip the {SSHA}
	eppS := string(encodedPassPhrase)[6:]
	hash, err := base64.StdEncoding.DecodeString(eppS)
	if err != nil {
		return false
	}
	sha := sha1.New()
	sha.Write(rawPassPhrase)
	sum := sha.Sum(nil)

	if bytes.Compare(sum, hash) != 0 {
		return false
	}
	return true
}

// makeSHAHash make hasing using SHA-1. This is not the final output though. You need to append {SHA} string with base64 of this hash.
func makeSHAHash(passphrase []byte) []byte {
	sha := sha1.New()
	sha.Write(passphrase)

	h := sha.Sum(nil)
	return h
}
