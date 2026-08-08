package xcrypt

import (
    "errors"
    "os"
    "strings"

    "github.com/terefang/gommons/pkg/xfile"
)

// ReadFromHtpasswd loads user credentials and roles from an htpasswd-style file.
//
// - Each non-empty, non-comment line must have the form "username:password".
// - Lines beginning with '/', '#', ';', ':', '%', '!', or '$' are treated as comments and ignored.
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
