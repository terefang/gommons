package xbytes

import (
    "encoding/base32"
    "encoding/base64"
    "encoding/hex"
    "strings"
)

func ToBase64(b []byte) []byte {
    return []byte(base64.StdEncoding.EncodeToString(b))
}

func FromBase64(b []byte) []byte {
    clean := strings.Join(strings.Fields(string(b)), "")
    _bytes, _err := base64.StdEncoding.DecodeString(clean)
    if _err != nil {
        return nil
    }
    return _bytes
}

func FromHex(s string) ([]byte, error) {
    clean := strings.ToUpper(strings.Join(strings.Fields(s), ""))
    return hex.DecodeString(clean)
}

func FromBase32(s string) ([]byte, error) {
    clean := strings.ToUpper(strings.Join(strings.Fields(s), ""))
    if n := len(clean) % 8; n != 0 {
        clean += "========"[:8-n]
    }
    return base32.StdEncoding.DecodeString(clean)
}
