package xcrypt

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/terefang/gommons/pkg/xfile"
	"github.com/terefang/gommons/pkg/xstrings"
)

// ReadFromHtgroup loads user groups/roles from an htgroup-style file.
//
// - Each non-empty, non-comment line must have the form "groupname:user1 .... userN".
// - Lines beginning with '/', '#', ';', ':', '%', '!', '-', or '$' are treated as comments and ignored.
// - Reading stops when a line containing only "END" is encountered.
// - Existing entries for the same groupname are overwritten by later entries in the file.
func ReadFromHtgroup(f string) (map[string][]string, error) {
	if xfile.FileExists(f) {
		groups := make(map[string][]string)

		_buf, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		_str := string(_buf)
		_lines := strings.Split(_str, "\n")
		for _, _line := range _lines {
			_sline := strings.TrimSpace(_line)
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
			if _sline[0] == '*' {
				continue
			}
			_gu := strings.SplitN(_sline, ":", 3)
			if len(_gu) >= 2 {
				groups[strings.ToLower(_gu[0])] = xstrings.SplitByDefaultSet(strings.ToLower(_gu[1]))
			} else {
				log.Println("htgroup: invalid line:", _sline)
			}
		}
		return groups, nil
	}
	return nil, errors.New("file not found: " + f)
}
