package xmatch

import (
	"errors"
	"path"
	"regexp"
	"strings"
)

func convertAntPatternToRegex(_p string) string {
	_p = regexp.QuoteMeta(_p)

	_p = strings.ReplaceAll(_p, "\\*\\*", ".*") // **
	_p = strings.ReplaceAll(_p, "\\*", "[^/]*") // *
	_p = strings.ReplaceAll(_p, "\\?", ".")     // ?

	if !strings.HasPrefix(_p, "^") {
		_p = "^" + _p
	}
	if !strings.HasSuffix(_p, "$") {
		_p = _p + "$"
	}

	return _p
}

func AntCompile(pattern string) (*regexp.Regexp, error) {
	// clean parameters
	if strings.Contains(pattern, "//") {
		pattern = path.Clean(pattern)
	}

	// pattern without wildcards is useless
	if !strings.ContainsAny(pattern, "*?") {
		return nil, errors.New("no wildcards in pattern")
	}

	// make rx from ant
	regexPattern := convertAntPatternToRegex(pattern)
	return regexp.Compile(regexPattern)
}

func AntMatchList(_pattern string, _paths ...string) ([]string, error) {
	_ret := make([]string, 0)
	_rx, err := AntCompile(_pattern)
	if err != nil {
		return _ret, err
	}
	for _, _path := range _paths {
		if strings.Contains(_path, "//") {
			_path = path.Clean(_path)
		}
		if _rx.MatchString(_path) {
			_ret = append(_ret, _path)
		}
	}
	return _ret, nil
}

func AntMatch(pattern string, pathStr string) bool {
	// short-circuit direct match
	if pattern == pathStr {
		return true
	}

	// clean parameters
	if strings.Contains(pattern, "//") {
		pattern = path.Clean(pattern)
	}
	if strings.Contains(pathStr, "//") {
		pathStr = path.Clean(pathStr)
	}

	// short-circuit cleaned direct match
	if pattern == pathStr {
		return true
	}

	// make rx from ant
	regex, err := AntCompile(pattern)
	if err != nil {
		return false
	}

	return regex.MatchString(pathStr)
}
