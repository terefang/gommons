package xcrypt

import (
	"errors"
	"os"
	"strings"

	"github.com/terefang/gommons/pkg/xfile"
)

func ReadFromHtpasswd(f string, isAllowAnyUser bool) (creds map[string]string, roles map[string]string, err error) {
	if xfile.FileExists(f) {
		creds = make(map[string]string)
		roles = make(map[string]string)

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
				return
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
			// default any user
			if (_sline[0] == '*') && (!isAllowAnyUser) {
				continue
			}
			_upr := strings.SplitN(_sline, ":", 3)
			creds[_upr[0]] = _upr[1]
			roles[_upr[0]] = _upr[2]
		}
		return creds, roles, nil
	}
	return nil, nil, errors.New("file not found: " + f)
}
