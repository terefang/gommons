package xmatch

import (
	"bytes"
	"errors"
	"net/url"
	"strings"
)

// CompileLdapishFilter compiles a filter string.
func CompileLdapishFilter(s string) (*Filter, error) {

	es, err := url.QueryUnescape(s)
	if err != nil {
		return nil, err
	}
	s = es

	if len(s) == 0 || s[0] != '(' {
		return nil, errCharZeroNotLParen
	}

	f, pos, err := compileLdapishFilter(s, 1)
	if err != nil {
		return nil, err
	}

	if pos != len(s) {
		return nil, errors.New("finished compiling filter with extra at end: " + s[pos:])
	}

	return f, nil
}

func compileLdapishFilterSet(
	s string, pos int, parent *Filter) (int, error) {

	for pos < len(s) && s[pos] == '(' {
		child, newPos, err := compileLdapishFilter(s, pos+1)
		if err != nil {
			return pos, err
		}
		pos = newPos
		parent.Children = append(parent.Children, child)
	}

	if pos == len(s) {
		return pos, errUnexpectedEOF
	}

	return pos + 1, nil
}

func compileLdapishFilter(s string, pos int) (*Filter, int, error) {

	switch s[pos] {
	case '(':
		f, newPos, err := compileLdapishFilter(s, pos+1)
		newPos++
		return f, newPos, err

	case '&':
		f := &Filter{Op: filterAnd}
		newPos, err := compileLdapishFilterSet(s, pos+1, f)
		return f, newPos, err

	case '|':
		f := &Filter{Op: filterOr}
		newPos, err := compileLdapishFilterSet(s, pos+1, f)
		return f, newPos, err

	case '!':
		f := &Filter{Op: filterNot}
		child, newPos, err := compileLdapishFilter(s, pos+1)
		f.Children = append(f.Children, child)
		return f, newPos, err

	default:

		var (
			f          *Filter
			abuf, cbuf bytes.Buffer
			newPos     = pos
		)

		for newPos < len(s) && s[newPos] != ')' {
			switch {
			case f != nil:
				if err := cbuf.WriteByte(s[newPos]); err != nil {
					return nil, 0, err
				}

			case s[newPos] == '=':
				f = &Filter{Op: filterEqualityMatch}

			case s[newPos] == '>' && s[newPos+1] == '=':
				f = &Filter{Op: filterGreaterOrEqual}
				newPos++

			case s[newPos] == '<' && s[newPos+1] == '=':
				f = &Filter{Op: filterLessOrEqual}
				newPos++

			case s[newPos] == '~' && s[newPos+1] == '~':
				f = &Filter{Op: filterRxMatch}
				newPos++

			case s[newPos] == '~' && s[newPos+1] == '=':
				f = &Filter{Op: filterApproxMatch}
				newPos++

			case s[newPos] == '^' && s[newPos+1] == '=':
				f = &Filter{Op: filterSubstringsPostfix}
				newPos++

			case s[newPos] == '$' && s[newPos+1] == '=':
				f = &Filter{Op: filterSubstringsPrefix}
				newPos++

			case s[newPos] == '*' && s[newPos+1] == '=':
				f = &Filter{Op: filterSubstrings}
				newPos++

			case f == nil:
				if err := abuf.WriteByte(s[newPos]); err != nil {
					return nil, 0, err
				}
			}
			newPos++
		}

		if newPos == len(s) {
			return f, newPos, errUnexpectedEOF
		}

		if f == nil {
			return nil, 0, errParse
		}

		// attributes are always case-insensitive
		f.Left = strings.ToLower(abuf.String())

		var (
			cbyt = cbuf.Bytes()
			cstr = cbuf.String()
			clen = len(cbyt)
			cfch = cbyt[clen-1]
		)

		switch {
		case f.Op == filterEqualityMatch && cstr == "*":
			f.Op = filterPresent

		case f.Op == filterEqualityMatch &&
			cbyt[0] == '*' && cfch == '*' && clen > 2:
			f.Op = filterSubstrings
			f.Right = cstr[1 : clen-1]

		case f.Op == filterEqualityMatch && cbyt[0] == '*' && clen > 1:
			f.Op = filterSubstringsPrefix
			f.Right = cstr[1:]

		case f.Op == filterEqualityMatch && cfch == '*' && clen > 1:
			f.Op = filterSubstringsPostfix
			f.Right = cstr[:clen-1]

		default:
			f.Right = cstr

		}
		newPos++
		return f, newPos, nil
	}
}
