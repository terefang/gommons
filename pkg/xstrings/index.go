package xstrings

import "strings"

func IndexOf(s, substr string, off int) int {
	if off == 0 {
		return strings.Index(s, substr)
	}

	n := strings.Index(s[off:], substr)
	if n == -1 {
		return -1
	}
	return (off + n)
}
