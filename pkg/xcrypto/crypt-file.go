package xcrypto

import (
	"errors"
	"os"

	"github.com/terefang/gommons/pkg/xfile"
)

func CrypterSetupSecretKey(l int) []byte {
	return GenerateSecretKeyWithSalt(l, nil)
}

func CrypterSetupSecretKeyFileWithPassphraseSeed(l int, _pass string, _salt []byte, f string) []byte {
	_key := GenerateSecretKeyWithSalt(l, _salt)
	if _pass == "" {
		WriteSecretToPem(_key, f, false)
	} else {
		_wkey, _ := WrapWithPassphrase(_pass, _key)
		WriteSecretToPem(_wkey, f, true)
	}
	return _key
}

func CrypterSetupFromFileWithPassphraseSeed(l int, _pass string, _salt []byte, f string, _make bool) []byte {
	_exists := xfile.FileExists(f)
	if !_exists && _make {
		return CrypterSetupSecretKeyFileWithPassphraseSeed(l, _pass, _salt, f)
	} else if !_exists {
		panic("file not found: " + f)
	}

	_pkey, _flag, _err := ReadSecretFromPem(f)
	if _err != nil {
		panic(_err)
	}

	if _pass != "" && _flag {
		_key, _err := UnwrapWithPassphrase(_pass, _pkey)
		if _err != nil {
			panic(_err)
		}
		return _key
	} else if !_flag {
		return _pkey
	}
	panic("encrypted secret but no passphrase given")
}

const ENCRYPTED_DATA_TAG = "ENCRYPTED DATA"

func CrypterEncryptFile(_key []byte, fin string, fout string, pem bool) error {
	if !xfile.FileExists(fin) {
		return errors.New("file not found: " + fin)
	}

	_data, _err := os.ReadFile(fin)
	if _err != nil {
		return _err
	}

	_edata, _err := WrapWithSecretKey(_key, _data)
	if _err != nil {
		return _err
	}

	if pem {
		return WriteTypeToPem(_edata, ENCRYPTED_DATA_TAG, fout)
	} else {
		return WriteTypeToDER(_edata, ENCRYPTED_DATA_TAG, fout)
	}
}

func CrypterDecryptFile(_key []byte, fin string, fout string, pem bool) error {
	if !xfile.FileExists(fin) {
		return errors.New("file not found: " + fin)
	}
	var _edata []byte
	var _type string
	var _err error
	if pem {
		_edata, _type, _err = ReadTypeFromPem(fin)
	} else {
		_edata, _type, _err = ReadTypeFromDER(fin)
	}
	if _err != nil {
		return _err
	}

	if _type != ENCRYPTED_DATA_TAG {
		return errors.New("encrypted file type not supported: " + _type)
	}

	_data, _err := UnwrapWithSecretKey(_key, _edata)
	if _err != nil {
		return _err
	}

	return os.WriteFile(fout, _data, os.FileMode(0600))
}
