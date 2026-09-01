package xstrings

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// SplitWsSq splits a string on whitespace respecting single-quoted (') strings.
func SplitWsSq(cmd string) []string {
	return SplitWsWithQuotes(cmd, '\'')
}

// SplitWsDq splits a string on whitespace respecting double-quoted (") strings.
func SplitWsDq(cmd string) []string {
	return SplitWsWithQuotes(cmd, '"')
}

// SplitWsTq splits a string on whitespace respecting tilted-quoted (`) strings.
func SplitWsTq(cmd string) []string {
	return SplitWsWithQuotes(cmd, '`')
}

// SplitWsWithQuotes splits a string on whitespace respecting the given quotation character.
func SplitWsWithQuotes(cmd string, qc rune) []string {
	var fields []string
	var cur strings.Builder
	inQuote := false
	var lastchar rune = -1
	for _, r := range cmd {
		switch {
		case lastchar == '\\' && r == qc:
			fallthrough
		case lastchar == '\\' && r == '\\':
			cur.WriteRune(r)
			lastchar = -1
			continue
		case lastchar == '\\':
			cur.WriteRune('\\')
			cur.WriteRune(r)
		case r == '\\':
			lastchar = r
			continue
		case r == qc && !inQuote:
			inQuote = true
		case r == qc && inQuote:
			inQuote = false
		case unicode.IsSpace(r) && !inQuote:
			if cur.Len() > 0 {
				fields = append(fields, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteRune(r)
		}
		lastchar = r
	}
	if cur.Len() > 0 {
		fields = append(fields, cur.String())
	}
	return fields
}

// RuneInSet reports whether the rune `r` exists within the slice `set`.
//
// Parameters:
//   - r:   The target rune to search for.
//   - set: A slice of runes to search within.
//
// Returns:
//   - true if `r` matches any element in `set`; otherwise false.
func RuneInSet(r rune, set []rune) bool {
	for _, e := range set {
		if e == r {
			return true
		}
	}
	return false
}

// RuneInSetIndex searches for the first occurrence of the rune `r` in the slice `set`
// and returns its 0-based index.
//
// Parameters:
//   - r:   The target rune to locate.
//   - set: A slice of runes to search within.
//
// Returns:
//   - The zero-based slice index of the first matching rune, or -1 if `r` is not present in `set`.
func RuneInSetIndex(r rune, set []rune) int {
	for i, e := range set {
		if e == r {
			return i
		}
	}
	return -1
}

// RuneInString reports whether the rune `r` exists anywhere within the `set` string.
//
// Parameters:
//   - r:   The target rune to search for.
//   - set: A string containing characters to search within.
//
// Returns:
//   - true if `r` is found in `set`; otherwise false.
//
// Note:
//   - Converts `set` into a rune slice, performing a linear scan.
func RuneInString(r rune, set string) bool {
	return RuneInSet(r, []rune(set))
}

// RuneInStringIndex searches for the first occurrence of rune `r` in the `set` string
// and returns its rune-based index.
//
// Parameters:
//   - r:   The target rune to locate.
//   - set: A string containing characters to search within.
//
// Returns:
//   - The 0-based rune index of `r` within `set`, or -1 if `r` is not present.
//
// Note:
//   - Converts `set` into a rune slice, returning the rune index rather than byte index.
func RuneInStringIndex(r rune, set string) int {
	return RuneInSetIndex(r, []rune(set))
}

var (
	// ErrMismatchedQuoteSets is returned when opening and closing quote slices differ in length.
	ErrMismatchedQuoteSets = errors.New("qs and qe slices must have equal length")
	// ErrUnclosedQuote is returned when a quote pair is opened but not closed before end-of-string.
	ErrUnclosedQuote = errors.New("unclosed quote delimiter")
	// ErrDanglingEscape is returned when a string ends with an unescaped backslash.
	ErrDanglingEscape = errors.New("dangling backslash escape at end of string")
)

// SplitWsWithDefaultQuotes parses a command string into whitespace-delimited fields
// using standard default quote delimiters (Quotes).
//
// Parameters:
//   - cmd: The input command string to be split into fields.
//
// Behavior:
//   - Uses pre-configured quote pairs defined by standard Quotes.
//   - Delegates execution to SplitWsWithQuoteSetFull with silent set to true, suppressing
//     syntax errors caused by unclosed quotes or trailing escape backslashes and returning
//     partially parsed fields instead.
//
// Returns:
//   - []string: A slice of parsed command field strings.
func SplitWsWithDefaultQuotes(cmd string) []string {
	_ret, _ := SplitWsWithQuoteSetFull(cmd, []rune(Quotes), []rune(Quotes), true)
	return _ret
}

// SplitWsWithExtendedQuotes parses a command string into whitespace-delimited fields
// using extended quote delimiters (Quotes + OpenBrackets / CloseBrackets).
//
// Parameters:
//   - cmd: The input command string to be split into fields.
//
// Behavior:
//   - Combines standard Quotes and matching OpenBrackets/CloseBrackets as quote pairs.
//   - Delegates parsing to SplitWsWithQuoteSetFull with silent set to true.
//   - Discards any returned errors.
//
// Returns:
//   - []string: A slice of parsed command field strings.
func SplitWsWithExtendedQuotes(cmd string) []string {
	_ret, _ := SplitWsWithQuoteSetFull(cmd, []rune(Quotes+OpenBrackets), []rune(Quotes+CloseBrackets), true)
	return _ret
}

// SplitWsWithQuoteSet wraps SplitWsWithQuoteSetFull only reporting essential errors.
func SplitWsWithQuoteSet(cmd string, qs []rune, qe []rune) ([]string, error) {
	return SplitWsWithQuoteSetFull(cmd, qs, qe, true)
}

// SplitWsWithQuoteSetFull splits a command string into whitespace-delimited fields,
// accounting for paired quote delimiters, backslash escape sequences, and error reporting flags.
//
// Parameters:
//   - cmd:    The input command string to be split into fields.
//   - qs:     A slice of opening quote runes (e.g., ['"', '\”, '(']).
//   - qe:     A slice of closing quote runes corresponding by index to `qs` (e.g., ['"', '\”, ')']).
//   - silent: When true, suppresses syntax validation errors for dangling backslashes
//     and unclosed quotes, returning the partial fields parsed up to EOF instead.
//
// Behavior:
//   - Whitespace Delimitation: Spaces separate fields only when outside of active quotes.
//   - Quote Handling: Entering an opening quote from `qs` pauses whitespace splitting until
//     its corresponding closing quote from `qe` (matched by slice index) is encountered.
//     Quote characters delimiter pairs are stripped from the resulting output fields.
//   - Escape Sequences: A backslash ('\') escapes opening quotes, closing quotes, and literal
//     backslashes. Non-special characters preceded by a backslash preserve both the backslash
//     and the character.
//   - Content Retention: Empty quoted strings (e.g., `""`) are preserved as empty string fields (`""`).
//
// Returns:
//   - []string: Parsed command fields.
//   - error:    Returns ErrMismatchedQuoteSets if len(qs) > len(qe). Unless silent is true,
//     returns ErrDanglingEscape for trailing backslashes or ErrUnclosedQuote
//     for unclosed quote pairs.
func SplitWsWithQuoteSetFull(cmd string, qs []rune, qe []rune, silent bool) ([]string, error) {
	if len(qs) > len(qe) {
		return nil, fmt.Errorf("%w: len(qs)=%d, len(qe)=%d", ErrMismatchedQuoteSets, len(qs), len(qe))
	}

	var fields []string
	var cur strings.Builder
	inQuote := -1
	var lastchar rune = -1
	hasContent := false

	for _, r := range cmd {
		if lastchar == '\\' {
			switch {
			case RuneInSet(r, qs), RuneInSet(r, qe), r == '\\':
				cur.WriteRune(r)
			default:
				cur.WriteRune('\\')
				cur.WriteRune(r)
			}
			hasContent = true
			lastchar = -1
			continue
		}

		if r == '\\' {
			lastchar = r
			continue
		}

		switch {
		case inQuote == -1 && RuneInSet(r, qs):
			inQuote = RuneInSetIndex(r, qs)
			hasContent = true
		case inQuote != -1 && r == qe[inQuote]:
			inQuote = -1
		case unicode.IsSpace(r) && inQuote == -1:
			if hasContent || cur.Len() > 0 {
				fields = append(fields, cur.String())
				cur.Reset()
				hasContent = false
			}
		default:
			cur.WriteRune(r)
			hasContent = true
		}
		lastchar = r
	}

	// Validate dangling trailing backslash
	if !silent && lastchar == '\\' {
		return nil, fmt.Errorf("%w at index %d", ErrDanglingEscape, len(cmd)-1)
	}

	// Validate unclosed quote state
	if !silent && inQuote != -1 {
		return nil, fmt.Errorf("%w: missing matching %q for %q", ErrUnclosedQuote, qe[inQuote], qs[inQuote])
	}

	if hasContent || cur.Len() > 0 {
		fields = append(fields, cur.String())
	}

	return fields, nil
}

// MustSplitWsWithQuoteSet wraps SplitWsWithQuoteSet and panics if an error is returned.
// Useful for package initializations or contexts where invalid syntax is fatal.
func MustSplitWsWithQuoteSet(cmd string, qs []rune, qe []rune) []string {
	fields, err := SplitWsWithQuoteSetFull(cmd, qs, qe, true)
	if err != nil {
		panic(fmt.Sprintf("MustSplitWsWithQuoteSet: %v", err))
	}
	return fields
}

// SplitBySet splits the string `cmd` into fields around any rune present in `set`.
// It acts similarly to strings.Fields, but matches against a custom set of characters
// rather than standard whitespace.
//
// Parameters:
//   - cmd: The input string to be split into fields.
//   - set: A string containing characters to treat as delimiters.
//
// Behavior:
//   - Consecutive delimiters are collapsed, and empty fields are omitted from the result.
//   - Returns an empty slice if `cmd` contains no non-delimiter characters.
//
// Returns:
//   - []string: A slice of substrings split around the characters in `set`.
func SplitBySet(cmd string, set string) []string {
	return strings.FieldsFunc(cmd, func(r rune) bool {
		for _, _r := range []rune(set) {
			if r == _r {
				return true
			}
		}
		return false
	})
}

// SplitByDefaultSet splits the input string `cmd` into fields using characters defined
// in CommonFieldSeparators as delimiters.
//
// Parameters:
//   - cmd: The input string to be split into fields.
//
// Behavior:
//   - Delegates parsing to SplitBySet using the CommonFieldSeparators global string.
//   - Collapses consecutive delimiters and omits empty strings from the resulting slice.
//
// Returns:
//   - []string: A slice of non-empty substrings split around common field separators.
func SplitByDefaultSet(cmd string) []string {
	return SplitBySet(cmd, CommonFieldSeparators)
}

// SplitPart splits the string `cmd` into two substrings around the first occurrence of `sep`.
//
// Parameters:
//   - cmd: The input string to be split.
//   - sep: The delimiter string used to split `cmd`.
//
// Behavior:
//   - Performs a slice operation up to a maximum of two parts using `strings.SplitN`.
//   - If `sep` is not found within `cmd`, the entire `cmd` string is returned as the first part,
//     and the second part returns as an empty string.
//
// Returns:
//   - head (string): The substring before the first occurrence of `sep` (or full `cmd` if unfound).
//   - tail (string): The substring after the first occurrence of `sep` (or "" if unfound).
func SplitPart(cmd string, sep string) (string, string) {
	_parts := strings.SplitN(cmd, sep, 2)
	if len(_parts) == 1 {
		return _parts[0], ""
	}
	return _parts[0], _parts[1]
}
