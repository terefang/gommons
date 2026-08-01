package xcrypto

import (
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

const SECRET_KEY_TAG = "SECRET KEY"
const ENCRYPTED_KEY_TAG = "ENCRYPTED KEY"

func EncodeSecretToPem(_b []byte, e bool) string {
	if e {
		return EncodeTypeToPem(_b, ENCRYPTED_KEY_TAG)
	}
	return EncodeTypeToPem(_b, SECRET_KEY_TAG)
}

func MakeSecretKeyFile(bits int, f string) []byte {
	return MakeSecretKeyFileWithSalt(nil, bits, f)
}

func MakeSecretKeyFileWithSalt(_salt []byte, bits int, f string) []byte {
	return MakeSecretKeyFileWithPassphraseSalt("", _salt, bits, f)
}

func MakeSecretKeyFileWithPassphrase(_pass string, bits int, f string) []byte {
	return MakeSecretKeyFileWithPassphraseSalt(_pass, nil, bits, f)
}

func MakeSecretKeyFileWithPassphraseSalt(_pass string, _salt []byte, bits int, f string) []byte {
	_len := CalcBytesFromBits(bits)
	_key := GenerateSecretKeyWithSalt(_len, _salt)
	if _pass != "" {
		_wkey, _err := WrapWithPassphrase(_pass, _key)
		if _err != nil {
			panic(_err)
		}
		WriteSecretToPem(_wkey, f, true)
	} else {
		WriteSecretToPem(_key, f, false)
	}
	return _key
}

func WriteSecretToPem(skey []byte, f string, e bool) error {
	//p7key, _ := xcrypt.PadPkcs7(skey, 8)
	_s := EncodeSecretToPem(skey, e)
	if f == "-" {
		fmt.Println(_s)
		return nil
	}
	return os.WriteFile(f, []byte(_s), os.FileMode(0600))
}

func DecodeSecretFromPem(_b []byte) ([]byte, bool, error) {
	_block, _ := pem.Decode(_b)
	if _block == nil {
		return nil, false, errors.New("failed to decode PEM block")
	}
	return _block.Bytes, _block.Type == ENCRYPTED_KEY_TAG, nil
}

func ReadSecretFromPem(f string) ([]byte, bool, error) {
	_pem, _err := os.ReadFile(f)
	if _err != nil {
		return nil, false, _err
	}
	return DecodeSecretFromPem(_pem)
}
