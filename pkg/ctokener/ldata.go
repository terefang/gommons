package ctokener

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
)

func LdataParseFile(f string) (any, error) {
	file, err := os.Open(f)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	return LdataParse(file)
}

func LdataParse(rdr io.Reader) (any, error) {
	var rootNode any
	_ct := NewTokener(rdr)
	_ct.ConfigureLdataSyntax()
	for true {
		_tk, err := _ct.NextToken()
		if err != nil {
			return nil, err
		}
		// the first token type make the root node
		if _tk == TOKEN_TYPE_UNKNOWN && _ct.CardinalValue == '{' {
			rootNode, err = LdataParseMap(_ct)
			return rootNode, nil
		} else if _tk == TOKEN_TYPE_UNKNOWN && _ct.CardinalValue == '[' {
			rootNode, err = LdataParseList(_ct)
			return rootNode, nil
		} else if _tk == TOKEN_TYPE_IDENTIFIER || _tk == TOKEN_TYPE_STRING || _tk == TOKEN_TYPE_STRING_MULTILINE {
			_ct.PushBack()
			rootNode, err = LdataParseMap(_ct)
			return rootNode, nil
		} else if _tk == TOKEN_TYPE_EOF {
			return rootNode, nil
		}
		return nil, errors.New(fmt.Sprintf("(%d:%d) [%s] FATAL INVALID-TOKEN (%s)", _ct.Pos.Line(), _ct.Pos.Col(), file_line(), string(_ct.CharacterValue)))
	}
	return rootNode, nil
}

func LdataParseMap(_ct *CustomTokener) (map[string]any, error) {
	nmap := make(map[string]any)
	err := LdataParseIntoMap(_ct, nmap)
	if err != nil {
		return nil, err
	}
	return nmap, nil
}

func LdataParseIntoMap(_ct *CustomTokener, nmap map[string]any) error {
	for nkey, err := LdataParseValue(_ct); (err == nil) && (nkey != nil); nkey, err = LdataParseValue(_ct) {
		_key, _ok := nkey.(string)
		if !_ok {
			_ikey, _ok := nkey.(int)
			if !_ok {
				errors.New(fmt.Sprintf("(%d:%d) [%s] FATAL INVALID-TOKEN-TYPE FOR KEY (%s)", _ct.Pos.Line(), _ct.Pos.Col(), file_line(), nkey))
			}
			_key = fmt.Sprintf("%d", _ikey)
		}

		if _key == "}" {
			return nil
		}

		nvalue, _err := LdataParseValue(_ct)
		if _err != nil {
			return _err
		}
		nmap[_key] = nvalue
	}
	// check if last toke was EOF
	if _ct.TokenType == TOKEN_TYPE_EOF {
		return nil
	}
	return errors.New(fmt.Sprintf("(%d:%d) [%s] FATAL ERROR", _ct.Pos.Line(), _ct.Pos.Col(), file_line()))
}

func LdataParseList(_ct *CustomTokener) ([]any, error) {
	nlist, err := LdataParseIntoList(_ct)
	if err != nil {
		return nil, err
	}

	return nlist, nil
}

func LdataParseIntoList(_ct *CustomTokener) ([]any, error) {
	nlist := make([]any, 0)
	for nvalue, err := LdataParseValue(_ct); err == nil; nvalue, err = LdataParseValue(_ct) {
		// EOF
		if nvalue == nil {
			return nlist, nil
		}

		_key, _ok := nvalue.(string)
		if _ok {
			if _key == ")" {
				return nlist, nil
			}
			if _key == "]" {
				return nlist, nil
			}
		}
		nlist = append(nlist, nvalue)
	}

	return nil, errors.New(fmt.Sprintf("(%d:%d) [%s] FATAL ERROR", _ct.Pos.Line(), _ct.Pos.Col(), file_line()))
}

func LdataParseValue(_ct *CustomTokener) (any, error) {
	_tk, err := _ct.NextToken()
	if err != nil {
		return nil, err
	}
	if _tk == TOKEN_TYPE_STRING {
		return _ct.StringValue, nil
	} else if _tk == TOKEN_TYPE_STRING_MULTILINE {
		return _ct.StringValue, nil
	} else if _tk == TOKEN_TYPE_IDENTIFIER {
		return _ct.StringValue, nil
	} else if _tk == TOKEN_TYPE_CARDINAL {
		return _ct.CardinalValue, nil
	} else if _tk == TOKEN_TYPE_NUMBER {
		return _ct.NumericValue, nil
	} else if _tk == TOKEN_TYPE_BYTES {
		return _ct.ByteValue, nil
	} else if _tk == TOKEN_TYPE_BOOL {
		return _ct.BooleanValue, nil
	} else if _tk == TOKEN_TYPE_UNKNOWN {
		if _ct.CharacterValue == '{' {
			return LdataParseMap(_ct)
		} else if _ct.CharacterValue == '(' {
			return LdataParseList(_ct)
		} else if _ct.CharacterValue == '[' {
			return LdataParseList(_ct)
		} else {
			return string(_ct.CharacterValue), nil
		}
	}

	if _tk == TOKEN_TYPE_EOF {
		return nil, nil
	}

	// TODO
	return nil, errors.New(fmt.Sprintf("(%d:%d) [%s] FATAL INVALID-TOKEN-TYPE (%d)", _ct.Pos.Line(), _ct.Pos.Col(), file_line(), _tk))
}

func file_line() string {
	_, fileName, fileLine, ok := runtime.Caller(1)
	var s string
	if ok {
		s = fmt.Sprintf("%s:%d", fileName, fileLine)
	} else {
		s = ""
	}
	return s
}
