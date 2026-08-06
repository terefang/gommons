package xotp

import "time"

func (f *OtpFob) TOTP() (string, error) {
    counter := time.Now().Unix()
    return f.TOTPWithTime(uint64(counter))
}

func (f *OtpFob) TOTPWithTime(t uint64) (string, error) {
    counter := t
    counter /= uint64(f.Period)
    passcode, err := f.GenerateCode(counter)
    if err != nil {
        return "", err
    }
    return passcode, nil
}

func (f *OtpFob) TOTPWithWindow(w int) ([]string, error) {
    counter := time.Now().Unix()
    return f.TOTPWithTimeAndWindow(uint64(counter), w)
}

func (f *OtpFob) TOTPWithTimeAndWindow(t uint64, w int) ([]string, error) {
    _ret := make([]string, 0)
    counter := t
    counter /= uint64(f.Period)
    for i := -w; i <= w; i++ {
        c0 := int64(counter) + int64(i)
        passcode, err := f.GenerateCode(uint64(c0))
        if err != nil {
            return nil, err
        }
        _ret = append(_ret, passcode)
    }
    return _ret, nil
}
