package xotp

import (
    "crypto/hmac"
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "crypto/sha512"
    "encoding/binary"
    "hash"
)

const AlgorithmSHA1 = "SHA1"
const AlgorithmSHA256 = "SHA256"
const AlgorithmSHA512 = "SHA512"
const AlgorithmMD5 = "MD5"

func (f OtpFob) MakeAlgorithm() func() hash.Hash {
    switch f.Algorithm {
    case AlgorithmSHA1:
        return sha1.New
    case AlgorithmSHA256:
        return sha256.New
    case AlgorithmSHA512:
        return sha512.New
    case AlgorithmMD5:
        return md5.New
    default:
        return sha1.New
    }
}

func (f *OtpFob) HOTP() (string, error) {
    passcode, err := f.GenerateCode(uint64(f.Counter))
    if err != nil {
        return "", err
    }
    f.Counter++
    return passcode, nil
}

// GenerateCode uses a counter and secret value and options struct to
// create a passcode.
func (f *OtpFob) GenerateCode(counter uint64) (passcode string, err error) {
    if f.SymbolSet == "" {
        f.SymbolSet = DefaultSymbolSet
    }
    sum, err := f.GenerateDigest(counter)
    if err != nil {
        return "", err
    }

    // "Dynamic truncation" in RFC 4226
    // http://tools.ietf.org/html/rfc4226#section-5.4
    offset := sum[len(sum)-1] & 0xf
    value := int64(((int(sum[offset]) & 0x7f) << 24) |
        ((int(sum[offset+1] & 0xff)) << 16) |
        ((int(sum[offset+2] & 0xff)) << 8) |
        (int(sum[offset+3]) & 0xff))

    radix := int64(len(f.SymbolSet))
    _str := make([]byte, f.Digits)

    for i := 0; i < f.Digits; i++ {
        digit := value % radix
        value /= radix
        c := f.SymbolSet[digit]
        _str[(f.Digits-1)-i] = c
    }

    return string(_str), nil
}

func (f OtpFob) GenerateDigest(counter uint64) (digest []byte, err error) {
    buf := make([]byte, 8)
    mac := hmac.New(f.MakeAlgorithm(), f.Key)
    binary.BigEndian.PutUint64(buf, counter)

    mac.Write(buf)
    sum := mac.Sum(nil)

    return sum, nil
}
