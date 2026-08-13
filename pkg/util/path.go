package util

import (
	"os"
	"path/filepath"
	"strings"
)

// TrimExt removes extension from s
func TrimExt(s string) string {
	idx := strings.LastIndex(s, ".")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

// ExtEqualFold returns true if s ends with extension (e.g. ".html")
// case-insensitive
func ExtEqualFold(s string, ext string) bool {
	e := filepath.Ext(s)
	return strings.EqualFold(e, ext)
}

func ExpandTildeInPath(s string) string {
	if strings.HasPrefix(s, "~") {
		dir, err := os.UserHomeDir()
		Must(err)
		return dir + s[1:]
	}
	return s
}
