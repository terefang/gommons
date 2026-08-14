package xcrypt

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/go-crypt/crypt"
	tcrypt "github.com/tredoe/crypt"
	"github.com/tredoe/crypt/apr1_crypt"
	"github.com/tredoe/crypt/md5_crypt"
	"github.com/tredoe/crypt/sha512_crypt"

	"github.com/terefang/gommons/pkg/xfile"
	"github.com/terefang/gommons/pkg/xotp"
)

// ReadFromHtpasswd loads user credentials and roles from an htpasswd-style file.
//
// - Each non-empty, non-comment line must have the form "username:password".
// - Lines beginning with '/', '#', ';', ':', '%', '!', '-', or '$' are treated as comments and ignored.
// - Reading stops when a line containing only "END" is encountered.
// - Entries whose username begins with '*' are ignored unless isAllowAnyUser is true.
// - Existing entries for the same username are overwritten by later entries in the file.
func ReadFromHtpasswd(f string, isAllowAnyUser bool) (map[string]string, map[string]string, error) {
	if xfile.FileExists(f) {
		creds := make(map[string]string)
		roles := make(map[string]string)

		_buf, err := os.ReadFile(f)
		if err != nil {
			return nil, nil, err
		}
		_str := string(_buf)
		_lines := strings.Split(_str, "\n")
		for _, _line := range _lines {
			_sline := strings.TrimSpace(string(_line))
			if len(_sline) == 0 {
				continue
			}
			if _sline == "END" {
				break
			}
			if _sline[0] == '/' {
				continue
			}
			if _sline[0] == '#' {
				continue
			}
			if _sline[0] == ';' {
				continue
			}
			if _sline[0] == ':' {
				continue
			}
			if _sline[0] == '%' {
				continue
			}
			if _sline[0] == '!' {
				continue
			}
			if _sline[0] == '$' {
				continue
			}
			if _sline[0] == '-' {
				continue
			}
			// default any user
			if (_sline[0] == '*') && (!isAllowAnyUser) {
				continue
			}
			_upr := strings.SplitN(_sline, ":", 3)
			if len(_upr) == 2 {
				creds[_upr[0]] = _upr[1]
				if len(_upr) == 3 {
					roles[_upr[0]] = _upr[2]
				}
			} else {
				log.Println("htpasswd: invalid line:", _sline)
			}
		}
		return creds, roles, nil
	}
	return nil, nil, errors.New("file not found: " + f)
}

func ValidateCryptedCredential(_given string, _encrypted string) (bool, error) {
	if strings.HasPrefix(_encrypted, "{plain}") {
		return _given == _encrypted[7:], nil
	} else if strings.HasPrefix(_encrypted, "{type7}") {
		return ValidateType7Credential(_given, _encrypted)
	} else if _fob, err := xotp.FromMCF(_encrypted); err == nil {
		_codes, err := _fob.TOTPWithWindow(1)
		if err == nil {
			for _, _code := range _codes {
				if _code == _given {
					return true, nil
				}
			}
		}
	} else {
		_dec, _ := crypt.NewDefaultDecoder()
		_verifier, _verr := _dec.Decode(_encrypted)
		if _verr != nil {
			_verifier, _verr := tcrypt.NewFromHash(_encrypted)
			if _verr != nil {
				return false, _verr
			}
			_err := _verifier.Verify(_encrypted, []byte(_given))
			if _err != nil {
				return false, _err
			}
		} else {
			return _verifier.MatchAdvanced(_given)
		}
	}
	return false, nil
}

func ValidateCryptedCredentialSimple(pwd string, p string) bool {
	_b, _ := ValidateCryptedCredential(pwd, p)
	return _b
}

func CryptApr1Credential(_given string) (string, error) {
	_crypt := apr1_crypt.New()
	return _crypt.Generate([]byte(_given), make([]byte, 0))
}

func Crypt6Credential(_given string) (string, error) {
	_crypt := sha512_crypt.New()
	return _crypt.Generate([]byte(_given), make([]byte, 0))
}

func Crypt1Credential(_given string) (string, error) {
	_crypt := md5_crypt.New()
	return _crypt.Generate([]byte(_given), make([]byte, 0))
}

func TotpCredential(_given string) (string, error) {
	_otp, err := xotp.From([]byte(_given), xotp.DefaultDigits, xotp.DefaultAlgorithm)
	if err != nil {
		return "", err
	}
	return _otp.ToSimple(), nil
}

func TotpCredentialUrl(_given string) (string, error) {
	_otp, err := xotp.From([]byte(_given), xotp.DefaultDigits, xotp.DefaultAlgorithm)
	if err != nil {
		return "", err
	}
	return _otp.ToURL(), nil
}
