package xotp

import (
    "errors"
    "fmt"
    "net/url"
    "strconv"
    "strings"

    "github.com/terefang/gommons/pkg/xbytes"
)

type OtpFob struct {
    Type      string
    Account   string
    Issuer    string
    Algorithm string
    Secret    string
    Key       []byte
    Digits    int
    SymbolSet string
    Counter   int
    Skew      int
    Period    int
}

const DefaultAlgorithm = "SHA1"
const DefaultPeriod = 30
const DefaultDigits = 6
const DefaultSymbolSet = `0123456789`
const SteamSymbolSet = `23456789BCDFGHJKMNPQRTVWXY`

// otpauth://totp/issuer:account?algorithm=SHA256&issuer=issuer&secret=WTQWMR5ZB7EP6BHH

func FromURL(uri string) (*OtpFob, error) {
    _uri, err := url.Parse(uri)
    if err != nil {
        return nil, err
    }
    if _uri.Scheme != "otpauth" && _uri.Scheme != "gauth" {
        return nil, errors.New("invalid uri scheme: " + _uri.Scheme)
    }

    _fob := &OtpFob{}
    _fob.Type = _uri.Host
    _fob.Account = _uri.Path

    if _uri.Query().Has("account") {
        _fob.Account = _uri.Query().Get("account")
    }

    if _uri.Query().Has("issuer") {
        _fob.Issuer = _uri.Query().Get("issuer")
    }

    _fob.Algorithm = DefaultAlgorithm
    if _uri.Query().Has("algorithm") {
        _fob.Algorithm = _uri.Query().Get("algorithm")
    }

    _fob.Digits = DefaultDigits
    if _uri.Query().Has("digits") {
        _fob.Digits, err = strconv.Atoi(_uri.Query().Get("digits"))
        if err != nil {
            return nil, err
        }
    }

    _fob.Counter = 0
    if _fob.Type == "hotp" {
        if _uri.Query().Has("counter") {
            _fob.Counter, err = strconv.Atoi(_uri.Query().Get("counter"))
            if err != nil {
                return nil, err
            }
        }
    }

    _fob.Period = DefaultPeriod
    _fob.Skew = 0
    if _fob.Type == "totp" {
        if _uri.Query().Has("period") {
            _fob.Period, err = strconv.Atoi(_uri.Query().Get("period"))
            if err != nil {
                return nil, err
            }
        }
        if _uri.Query().Has("skew") {
            _fob.Skew, err = strconv.Atoi(_uri.Query().Get("skew"))
            if err != nil {
                return nil, err
            }
        }
    }

    if _uri.Query().Has("secret") {
        _fob.Secret = _uri.Query().Get("secret")
        // url is always b32 encoded
        _fob.Key, err = xbytes.FromBase32(_fob.Secret)
        if err != nil {
            return nil, err
        }
    }

    return _fob, nil
}

func (f OtpFob) ToURL() string {
    _sb := strings.Builder{}
    _sb.WriteString("otpauth://")
    _sb.WriteString(f.Type)
    _sb.WriteString("/")
    if f.Issuer != "" {
        _sb.WriteString(url.QueryEscape(f.Issuer))
        _sb.WriteString(":")
    }
    if f.Account != "" {
        _sb.WriteString(url.QueryEscape(f.Account))
    } else {
        _sb.WriteString("TOKEN")
    }
    _sb.WriteString("?algorithm=")
    _sb.WriteString(url.QueryEscape(f.Algorithm))
    if f.Issuer != "" {
        _sb.WriteString("&issuer=")
        _sb.WriteString(url.QueryEscape(f.Issuer))
    }
    if f.Account != "" {
        _sb.WriteString("&account=")
        _sb.WriteString(url.QueryEscape(f.Account))
    }
    _sb.WriteString("&digits=")
    _sb.WriteString(fmt.Sprintf("%d", f.Digits))
    _sb.WriteString("&period=")
    _sb.WriteString(fmt.Sprintf("%d", f.Period))
    _sb.WriteString("&secret=")
    _sb.WriteString(string(xbytes.ToBase32(f.Key)))
    return _sb.String()
}

func FromHex(_b16 string, _digits int, _algo string) (*OtpFob, error) {
    _str, _ := xbytes.FromHex(_b16)
    return From(_str, _digits, _algo)
}

func FromB64(_b64 string, _digits int, _algo string) (*OtpFob, error) {
    _str := xbytes.FromBase64([]byte(_b64))
    return From(_str, _digits, _algo)
}

func FromB32(_b32 string, _digits int, _algo string) (*OtpFob, error) {
    _str, _ := xbytes.FromBase32(_b32)
    return From(_str, _digits, _algo)
}

func From(_key []byte, _digits int, _algo string) (*OtpFob, error) {

    _fob := &OtpFob{}
    _fob.Type = "totp"

    _fob.Algorithm = DefaultAlgorithm
    if _algo != "" {
        _fob.Algorithm = _algo
    }

    _fob.Digits = DefaultDigits
    if _digits > 0 {
        _fob.Digits = _digits
    }

    _fob.Counter = 0

    _fob.Period = DefaultPeriod
    _fob.Skew = 0
    _fob.Key = _key

    return _fob, nil
}
