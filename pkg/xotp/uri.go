package xotp

import (
    "errors"
    "fmt"
    "net/url"
    "strconv"

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

    fmt.Println(_fob)
    return _fob, nil
}
