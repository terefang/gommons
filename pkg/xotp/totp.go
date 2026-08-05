package xotp

import "encoding/base64"

// DefaultTOTP generates a TOTP for the current time step using the default
// settings (compatible with Google Authenticator) based on the given key.
// An error is reported if the key is invalid.
func DefaultTOTP(key string) (string, error) {
    std, err := Config{}.WithKey(key)
    if err != nil {
        return "", err
    }
    return std.TOTP(), nil
}

// TOTP returns the TOTP code for the current time step.  If the current time
// step value is t, this is equivalent to c.HOTP(t).
func (c Config) TOTP() string {
    return c.HOTP(c.timeStepWindow())
}

var SteamAlphabet = "23456789BCDFGHJKMNPQRTVWXY"
var SteamDigits = 5

func SteamTOTP(key string) (string, error) {
    std, err := Config{}.WithKey(key)
    if err != nil {
        return "", err
    }
    std.Digits = SteamDigits
    std.Format = FormatAlphabet(SteamAlphabet)
    return std.TOTP(), nil
}

func SteamVerifyCode(key string) (string, error) {
    std, err := Config{}.WithKey(key)
    if err != nil {
        return "", err
    }
    std.Digits = 8
    std.Format = func(hash []byte, nb int) string {
        return base64.StdEncoding.EncodeToString(hash)[:nb]
    }
    return std.TOTP(), nil
}
