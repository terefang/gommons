package xotp

import (
    "testing"

    "github.com/terefang/gommons/pkg/xbytes"
)

type hotpVector struct {
    C int
    P string
}

var testhotpVectors []hotpVector = []hotpVector{
    {0, "755224"},
    {1, "287082"},
    {2, "359152"},
    {3, "969429"},
    {4, "338314"},
    {5, "254676"},
    {6, "287922"},
    {7, "162583"},
    {8, "399871"},
    {9, "520489"},
    {10, "403154"},
    {11, "481090"},
}

func TestOtpFob_GenerateCode(t *testing.T) {
    _fob := &OtpFob{}
    _fob.Key = []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30}
    _fob.Digits = DefaultDigits
    _fob.SymbolSet = DefaultSymbolSet
    _fob.Algorithm = DefaultAlgorithm

    for _, v := range testhotpVectors {
        _code, err := _fob.GenerateCode(uint64(v.C))
        if err != nil {
            t.Errorf("GenerateCode(%d): %s", v.C, err)
        }
        if _code != v.P {
            t.Errorf("GenerateCode(%d): expected %s, got %s", v.C, v.P, _code)
        }
    }

    _fob.Key, _ = xbytes.FromBase32("GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ")
    for _, v := range testhotpVectors {
        _code, err := _fob.GenerateCode(uint64(v.C))
        if err != nil {
            t.Errorf("GenerateCode(%d): %s", v.C, err)
        }
        if _code != v.P {
            t.Errorf("GenerateCode(%d): expected %s, got %s", v.C, v.P, _code)
        }
    }
}
