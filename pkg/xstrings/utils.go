package xstrings

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/terefang/gommons/pkg/util"
)

// NormalizeNewlinesInPlace changes CRLF (Windows) and
// CR (Mac) to LF (Unix)
// Optimized for speed, modifies data in place
func NormalizeNewlinesInPlace(d []byte) []byte {
	wi := 0
	n := len(d)
	for i := 0; i < n; i++ {
		c := d[i]
		// 13 is CR
		if c != 13 {
			d[wi] = c
			wi++
			continue
		}
		// replace CR (mac / win) with LF (unix)
		d[wi] = 10
		wi++
		if i < n-1 && d[i+1] == 10 {
			// this was CRLF, so skip the LF
			i++
		}

	}
	return d[:wi]
}

// NormalizeNewlines is like NormalizeNewlinesInPlace but
// slower because it makes a copy of data
func NormalizeNewlines(d []byte) []byte {
	d = append([]byte{}, d...)
	return NormalizeNewlinesInPlace(d)
}

// Capitalize does foo => Foo, BAR => Bar etc.
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	s = strings.ToLower(s)
	return strings.ToUpper(s[0:1]) + s[1:]
}

// TrimPrefix is like strings.TrimPrefix but also returns a bool
// indicating that the string was trimmed
func TrimPrefix(s string, prefix string) (string, bool) {
	s2 := strings.TrimPrefix(s, prefix)
	return s2, len(s) != len(s2)
}

func ToTrimmedLines(d []byte) []string {
	lines := strings.Split(string(d), "\n")
	i := 0
	for _, l := range lines {
		l = strings.TrimSpace(l)
		// remove empty lines
		if len(l) > 0 {
			lines[i] = l
			i++
		}
	}
	return lines[:i]
}

func AppendNewline(s *string) string {
	if strings.HasSuffix(*s, "\n") {
		return *s
	}
	*s = *s + "\n"
	return *s
}

func CollapseMultipleNewlines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n") // CRLF => CR
	prev := ""
	for prev != s {
		prev = s
		s = strings.ReplaceAll(s, "\n\n\n", "\n\n")
	}
	return s
}

func AppendOrReplaceInText(orig string, toAppend string, delim string) string {
	AppendNewline(&toAppend)
	AppendNewline(&delim)
	content := "\n\n" + delim + toAppend + delim
	if strings.Contains(orig, content) {
		return CollapseMultipleNewlines(orig)
	}
	start := strings.Index(orig, delim)
	if start < 0 {
		return CollapseMultipleNewlines(orig + content)
	}
	end := strings.Index(orig[start+1:], delim)
	util.PanicIf(end == -1, "didn't find end delim")
	end += start + 1
	orig = orig[:start] + "\n\n" + orig[end+len(delim):]
	res := AppendNewline(&orig) + content
	return CollapseMultipleNewlines(res)
}

func AppendOrReplaceInFileMust(path string, toAppend string, delim string) bool {
	st, err := os.Lstat(path)
	util.Must(err)
	perm := st.Mode().Perm()
	orig, err := os.ReadFile(path)
	util.Must(err)
	newContent := AppendOrReplaceInText(string(orig), toAppend, delim)
	if newContent == string(orig) {
		return false
	}
	err = os.WriteFile(path, []byte(newContent), perm)
	util.Must(err)
	return true
}

/*
DeleteWhiteSpace deletes all whitespaces from a string as defined by unicode.IsSpace(rune).
It returns the string without whitespaces.

Parameter:

	str - the string to delete whitespace from, may be nil

Returns:

	the string without whitespaces
*/
func DeleteWhiteSpace(str string) string {
	if str == "" {
		return str
	}
	sz := len(str)
	var chs bytes.Buffer
	count := 0
	for i := 0; i < sz; i++ {
		ch := rune(str[i])
		if !unicode.IsSpace(ch) {
			chs.WriteRune(ch)
			count++
		}
	}
	if count == sz {
		return str
	}
	return chs.String()
}

// Typically returned by functions where a searched item cannot be found
const INDEX_NOT_FOUND = -1

/*
IndexOfDifference compares two strings, and returns the index at which the strings begin to differ.

Parameters:

	str1 - the first string
	str2 - the second string

Returns:

	the index where str1 and str2 begin to differ; -1 if they are equal
*/
func IndexOfDifference(str1 string, str2 string) int {
	if str1 == str2 {
		return INDEX_NOT_FOUND
	}
	if IsEmpty(str1) || IsEmpty(str2) {
		return 0
	}
	var i int
	for i = 0; i < len(str1) && i < len(str2); i++ {
		if rune(str1[i]) != rune(str2[i]) {
			break
		}
	}
	if i < len(str2) || i < len(str1) {
		return i
	}
	return INDEX_NOT_FOUND
}

// DefaultString returns the provided default string if the target string is empty;
// otherwise, it returns the original string.
//
// Parameters:
//   - str: The primary string to evaluate.
//   - defaultStr: The fallback string returned if `str` is empty.
//
// Behavior:
//   - Uses IsEmpty to check if `str` has a length of zero.
//   - Returns `defaultStr` if `str` is empty, otherwise returns `str` untouched.
//
// Returns:
//   - string: The original string `str` or the fallback `defaultStr`.
func DefaultString(str string, defaultStr string) string {
	if IsEmpty(str) {
		return defaultStr
	}
	return str
}

// DefaultIfBlank returns the provided default string if the target string is blank
// (empty or consisting entirely of whitespace); otherwise, it returns the original string.
//
// Parameters:
//   - str: The primary string to evaluate.
//   - defaultStr: The fallback string returned if `str` is blank.
//
// Behavior:
//   - Uses IsBlank to check if `str` is empty or contains only whitespace characters.
//   - Returns `defaultStr` if `str` is blank, otherwise returns `str` untouched.
//
// Returns:
//   - string: The original string `str` or the fallback `defaultStr`.
func DefaultIfBlank(str string, defaultStr string) string {
	if IsBlank(str) {
		return defaultStr
	}
	return str
}

// ParseEnvMust parses a byte slice containing environment variable definitions (.env format)
// into a key-value map, panicking if syntax errors are encountered.
//
// Parameters:
//   - d: The raw .env file byte slice.
//
// Behavior:
//   - Follows identical parsing rules to ParseEnv.
//   - Panics if any non-empty, non-comment line lacks an '=' separator.
//
// Returns:
//   - map[string]string: A map populated with parsed environment keys and values.
//
// Panics:
//   - Panics if a line cannot be parsed into a key-value pair.
func ParseEnvMust(d []byte) map[string]string {
	m, err := parseEnv(d, true)
	if err != nil {
		panic(err)
	}
	return m
}

// ParseEnv parses a byte slice containing environment variable definitions (.env format)
// into a key-value map.
//
// Parameters:
//   - d: The raw .env file byte slice.
//
// Behavior:
//   - Normalizes all line endings (CRLF/CR -> LF).
//   - Ignores empty lines and lines starting with comment markers ('#').
//   - Splits lines on the first '=' character into key and value pairs.
//   - Trims leading and trailing whitespace from keys and values.
//   - Silently ignores invalid lines (lines without an '=') without returning an error.
//
// Returns:
//   - map[string]string: A map populated with parsed environment keys and values.
func ParseEnv(d []byte) map[string]string {
	m, _ := parseEnv(d, false)
	return m
}

func parseEnv(d []byte, p bool) (map[string]string, error) {
	d = NormalizeNewlines(d)
	s := string(d)
	lines := strings.Split(s, "\n")
	m := make(map[string]string)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			m[key] = val
		} else if p {
			return nil, errors.New(fmt.Sprintf("invalid line '%s' in .env", line))
		}
	}
	return m, nil
}
