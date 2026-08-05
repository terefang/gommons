package xotp

// DefaultHOTP generates an HTOP for the specified counter using the default
// settings (compatible with Google Authenticator) based on the given key.
// An error is reported if the key is invalid.
func DefaultHOTP(key string, counter uint64) (string, error) {
    std, err := Config{}.WithKey(key)
    if err != nil {
        return "", err
    }
    return std.HOTP(counter), nil
}

// HOTP returns the HOTP code for the specified counter value.
func (c Config) HOTP(counter uint64) string {
    nd := c.digits()
    code := c.format(c.hmac(counter), nd)
    //if len(code) != nd {
    //    panic(fmt.Sprintf("invalid code length: got %d, want %d", len(code), nd))
    //}
    if len(code) > nd {
        return code[:nd]
    }
    return code
}

// Next increments the counter and returns the HOTP corresponding to its new value.
func (c *Config) Next() string { c.Counter++; return c.HOTP(c.Counter) }
