package xcrypto

import (
	"encoding/asn1"
	"os"
)

type derValue struct {
	Type string
	Data []byte
}

func EncodeTypeToDER(_b []byte, t string) []byte {
	_pem := derValue{Type: t, Data: _b}
	_bytes, _err := asn1.Marshal(_pem)
	if _err != nil {
		panic(_err)
	}
	return _bytes
}

func WriteTypeToDER(_b []byte, t string, f string) error {
	_pem := EncodeTypeToDER(_b, t)
	return os.WriteFile(f, []byte(_pem), os.FileMode(0600))
}

func DecodeTypeFromDER(_b []byte) ([]byte, string, error) {
	_pem := derValue{}
	_, _err := asn1.Unmarshal(_b, &_pem)
	if _err != nil {
		return nil, "", _err
	}
	return _pem.Data, _pem.Type, nil
}

func ReadTypeFromDER(f string) ([]byte, string, error) {
	_pem, _err := os.ReadFile(f)
	if _err != nil {
		return nil, "", _err
	}
	return DecodeTypeFromDER(_pem)
}

// ---------------------------------------------------------------------

type EncValue struct {
	Type   string
	Cipher string
	Mode   string
	Pad    string
	Data   []byte
}

func EncodeEvToDER(_ev EncValue) []byte {
	_bytes, _err := asn1.Marshal(_ev)
	if _err != nil {
		panic(_err)
	}
	return _bytes
}

func WriteEvToDER(_ev EncValue, f string) error {
	_pem := EncodeEvToDER(_ev)
	return os.WriteFile(f, []byte(_pem), os.FileMode(0600))
}

func DecodeEvFromDER(_b []byte) (*EncValue, error) {
	_ev := &EncValue{}
	_, _err := asn1.Unmarshal(_b, _ev)
	if _err != nil {
		return nil, _err
	}
	return _ev, nil
}

func ReadEvFromDER(f string) (*EncValue, error) {
	_pem, _err := os.ReadFile(f)
	if _err != nil {
		return nil, _err
	}
	return DecodeEvFromDER(_pem)
}
