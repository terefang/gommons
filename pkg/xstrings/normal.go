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
//   - NormalizeUpper("  hello-world! 123  ") -> "HELLO_WORLD__123"
func NormalizeUpper(s string) string {
	s = strings.ToUpper(s)
	rs := []rune(s)
	_sb := strings.Builder{}
	for i := 0; i < len(rs); i++ {
		if rs[i] >= 'A' && rs[i] <= 'Z' {
			_sb.WriteRune(rs[i])
		} else if rs[i] >= '0' && rs[i] <= '9' {
			_sb.WriteRune(rs[i])
		} else {
			_sb.WriteRune('_')
		}
	}
	return strings.Trim(_sb.String(), "_")
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
//   - NormalizeUpperCompact("  hello---world!! 123  ") -> "HELLO_WORLD_123"
func NormalizeUpperCompact(s string) string {
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
				_sb.WriteRune('_')
			}
			lastWasUnderscore = true
		}
	}
	return strings.Trim(_sb.String(), "_")
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
