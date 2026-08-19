package xstrings

import (
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

func SplitByDefaultSet(cmd string) []string {
	return SplitBySet(cmd, CommonFieldSeparators)
}

func SplitPart(cmd string, sep string) (string, string) {
	_parts := strings.SplitN(cmd, sep, 2)
	if len(_parts) == 1 {
		return _parts[0], ""
	}
	return _parts[0], _parts[1]
}
