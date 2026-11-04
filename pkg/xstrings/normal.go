package xstrings

import (
	"strings"
	"unicode"

	anyascii "github.com/anyascii/go"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// NormalizeUpper converts the input string `s` into an uppercase, alphanumeric-friendly format.
//
// Parameters:
//   - s: The input string to be normalized.
//
// Behavior:
//   - Converts all ASCII/Unicode characters in `s` to uppercase using strings.ToUpper.
//   - Iterates through the runes of the uppercase string:
//   - Retains uppercase ASCII letters ('A'-'Z') and digits ('0'-'9').
//   - Replaces all other characters (punctuation, whitespace, special symbols) with an underscore ('_').
//   - Trims any leading and trailing underscores from the resulting string.
//
// Returns:
//   - string: The normalized uppercase string containing only ASCII letters, digits, and inner underscores.
//
// Example:
//   - NormalizeUpper("  hello-world! 123  ", '_') -> "HELLO_WORLD__123"
func NormalizeUpper(s string, p rune) string {
	s = strings.ToUpper(s)
	rs := []rune(s)
	_sb := strings.Builder{}
	for i := 0; i < len(rs); i++ {
		if rs[i] >= 'A' && rs[i] <= 'Z' {
			_sb.WriteRune(rs[i])
		} else if rs[i] >= '0' && rs[i] <= '9' {
			_sb.WriteRune(rs[i])
		} else {
			_sb.WriteRune(p)
		}
	}
	return strings.Trim(_sb.String(), string(p))
}

// NormalizeUpperCompact transforms an input string `s` into a compact, uppercase alphanumeric string.
//
// Parameters:
//   - s: The input string to be normalized.
//
// Behavior:
//   - Converts all characters in `s` to uppercase using `strings.ToUpper`.
//   - Iterates through the runes of the uppercase string:
//   - Retains ASCII uppercase letters ('A'-'Z') and digits ('0'-'9').
//   - Replaces contiguous sequences of non-alphanumeric characters (spaces, punctuation, symbols)
//     with a single, collapsed underscore ('_').
//   - Trims any leading and trailing underscores from the final string.
//
// Returns:
//   - string: The normalized uppercase string containing only ASCII letters, digits, and non-consecutive inner underscores.
//
// Example:
//   - NormalizeUpperCompact("  hello---world!! 123  ", '_') -> "HELLO_WORLD_123"
func NormalizeUpperCompact(s string, p rune) string {
	s = strings.ToUpper(s)
	rs := []rune(s)
	_sb := strings.Builder{}
	lastWasUnderscore := false
	for i := 0; i < len(rs); i++ {
		if rs[i] >= 'A' && rs[i] <= 'Z' {
			_sb.WriteRune(rs[i])
			lastWasUnderscore = false
		} else if rs[i] >= '0' && rs[i] <= '9' {
			_sb.WriteRune(rs[i])
			lastWasUnderscore = false
		} else {
			if !lastWasUnderscore {
				_sb.WriteRune(p)
			}
			lastWasUnderscore = true
		}
	}
	return strings.Trim(_sb.String(), string(p))
}

// StripDiacritics removes combining accent marks from Latin characters via NFD.
func StripDiacritics(s string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	res, _, err := transform.String(t, s)
	return res, err
}

// ToASCII converts any UTF-8 string into a plain ASCII approximation,
// handling accents, extended Latin, non-Latin scripts (CJK, Cyrillic, Greek, etc.), and emojis.
func ToASCII(s string) string {
	// Step 1: Strip combining diacritics
	cleaned, err := StripDiacritics(s)
	if err != nil {
		cleaned = s // fallback to raw string if transformation fails
	}

	// Step 2: Transliterate remaining non-ASCII characters and scripts to ASCII
	return anyascii.Transliterate(cleaned)
}

// RemoveAccents removes accents from the string. The resulting string only has
// the letters from English alphabet. For example, "résumé" becomes "resume".
func RemoveAccents(str string) string {
	res, _ := StripDiacritics(str)
	return res
}

const naturalSortPadWidth = 20

// NaturalSortKey returns a string key that sorts naturally when compared
// lexicographically. Numeric segments (ASCII digits) are zero-padded to 20
// digits. Case is folded to lowercase. "file2" sorts before "file10".
func NaturalSortKey(s string) string {
	if s == "" {
		return ""
	}
	var b strings.Builder
	b.Grow(len(s) * 2)
	i := 0
	runes := []rune(strings.ToLower(s))
	for i < len(runes) {
		r := runes[i]
		if r >= '0' && r <= '9' {
			// Collect the full digit run.
			start := i
			for i < len(runes) && runes[i] >= '0' && runes[i] <= '9' {
				i++
			}
			digits := string(runes[start:i])
			if len(digits) < naturalSortPadWidth {
				b.WriteString(strings.Repeat("0", naturalSortPadWidth-len(digits)))
			}
			b.WriteString(digits)
		} else {
			b.WriteRune(unicode.ToLower(r))
			i++
		}
	}
	return b.String()
}
