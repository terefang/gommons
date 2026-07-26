package stemplate

import (
	"os"
	"strings"
)

// gtthEnv returns the contents of an environmental variable.
func gtthEnv(s string) string {
	return (os.Getenv(s))
}

// gtthSplit converts a string to an array.
func gtthSplit(in string, delim string) []string {
	return strings.Split(in, delim)
}
