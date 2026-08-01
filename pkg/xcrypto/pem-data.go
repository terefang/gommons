package xcrypto

import (
    "encoding/pem"
    "errors"
    "os"
)

func EncodeTypeToPem(_b []byte, t string) string {
    _pem := &pem.Block{Type: t, Bytes: _b, Headers: make(map[string]string)}
    pemBytes := pem.EncodeToMemory(_pem)
    return string(pemBytes)
}

func WriteTypeToPem(_b []byte, t string, f string) error {
    _pem := EncodeTypeToPem(_b, t)
    return os.WriteFile(f, []byte(_pem), os.FileMode(0600))
}

func DecodeTypeFromPem(_b []byte) ([]byte, string, error) {
    _block, _ := pem.Decode(_b)
    if _block == nil {
        return nil, "", errors.New("failed to decode PEM block")
    }
    return _block.Bytes, _block.Type, nil
}

func ReadTypeFromPem(f string) ([]byte, string, error) {
    _pem, _err := os.ReadFile(f)
    if _err != nil {
        return nil, "", _err
    }
    return DecodeTypeFromPem(_pem)
}