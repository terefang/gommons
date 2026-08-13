package xstrings

import "strings"

func NormalizeUpper(s string) string {
	s = strings.ToUpper(s)
	_sb := strings.Builder{}
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			_sb.WriteByte(byte(s[i]))
		} else if s[i] >= '0' && s[i] <= '9' {
			_sb.WriteByte(byte(s[i]))
		} else {
			_sb.WriteByte('_')
		}
	}
	return strings.Trim(_sb.String(), "_")
}
