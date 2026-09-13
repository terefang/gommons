package xcrypto

import (
	"encoding/hex"
	"fmt"
	"os"
)

func WriteSecretToHex(skey []byte, f string) error {
	_s := EncodeSecretToHex(skey)
	if f == "-" {
		fmt.Println(_s)
		return nil
	}
	return os.WriteFile(f, []byte(_s), os.FileMode(0600))
}

func EncodeSecretToHex(skey []byte) string {
	return hex.EncodeToString(skey)
}

func DecodeSecretFromHex(_b []byte) ([]byte, error) {
	_block, err := hex.DecodeString(string(_b))
	if err != nil {
		return nil, err
	}
	return _block, nil
}

func ReadSecretFromHex(f string) ([]byte, error) {
	_pem, _err := os.ReadFile(f)
	if _err != nil {
		return nil, _err
	}
	return DecodeSecretFromHex(_pem)
}
