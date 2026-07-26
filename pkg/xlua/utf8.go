package xlua

import (
	"errors"
	"fmt"

	rt "github.com/arnodel/golua/runtime"
)

func RegisterUtf8Functions(r *rt.Runtime) error {

	_tab := r.GlobalEnv().Get(rt.StringValue("utf8")).AsTable()
	// from Luay Extensions
	//utf8.charat (s, n) -> string
	r.SetEnvGoFunc(_tab, "charat", utf8CharAt, 2, false)
	//
	//utf8.firstchar (s) -> string
	//
	//utf8.lastchar (s) -> string
	//
	//utf8.length (s) -> int
	r.SetEnvGoFunc(_tab, "length", utf8Length, 1, false)

	//utf8.substr (s, n [, i]) -> string
	//
	//utf8.cutstr (s, n) -> str1, str2
	//
	//utf8.indexof (s, n [, i]) -> int/nil
	//
	//utf8.indexofany (s, n1 [, ..., nN]) -> int/nil
	//
	//utf8.lastindexof (s, n [, i]) -> int/nil
	//
	//utf8.lastindexofany (s, n1 [, ..., nN]) -> int/nil
	//
	//utf8.lower (s) -> string
	//
	//utf8.upper (s) -> string
	//
	//utf8.find(hay, needle [, start]) -> offset or nil
	//
	//utf8.count(hay, needle [, start]) -> offset of nil
	//
	//utf8.contains(hay, needle) -> boolean
	//
	//utf8.matches(hay, rx) -> boolean
	//
	//utf8.matchesany(hay, rx, ...) -> boolean
	//
	//utf8.startswith(hay, needle) -> boolean
	//
	//utf8.startswithany ( string, string [, ...]) -> boolean
	//
	//utf8.startswithnocase ( string, string [, ...] ) -> boolean
	//
	//utf8.endswith(hay, needle) -> boolean
	//
	//utf8.endswithany ( string, string [, ...]) -> boolean
	//
	//utf8.endswithnocase ( string, string [, ...] ) -> boolean
	//
	//utf8.rxsplit(hay, rx) -> list
	//
	//utf8.split(hay, sep) -> list
	//
	//utf8.join(sep, hay, ...) -> string
	//
	//utf8.format(fmt, ...) -> string
	//
	//utf8.repeat(part, int) -> string
	//
	//utf8.trim(hay) -> string
	//
	//utf8.ltrim(hay) -> string
	//
	//utf8.rtrim(hay) -> string
	//
	//utf8.strip(hay) -> string
	//
	//utf8.striptrailing(hay) -> string
	//
	//utf8.stripleading(hay) -> string
	//
	//utf8.replace(string, charFrom, charTo) -> string
	//
	//utf8.replace(string, table) -> string
	//
	//utf8.rxreplace(string, rx, to) -> string
	//
	//utf8.rxreplaceall(string, rx, to) -> string
	return nil
}

func utf8Length(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	_str, _err := c.StringArg(0)
	if _err != nil {
		return nil, _err
	}
	_runes := []rune(_str)
	_lr := int64(len(_runes))
	return c.PushingNext(t.Runtime, rt.IntValue(_lr)), nil
}

func utf8CharAt(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	_str, _err := c.StringArg(0)
	if _err != nil {
		return nil, _err
	}
	_idx, _err := c.IntArg(1)
	if _err != nil {
		return nil, _err
	}

	_runes := []rune(_str)
	_lr := int64(len(_runes))
	if _lr > _idx {
		return nil, errors.New(fmt.Sprintf("index out of range %d <= %d", _idx, _lr))
	}

	return c.PushingNext(t.Runtime, rt.StringValue(string([]rune{_runes[_idx-1]}))), nil
}
