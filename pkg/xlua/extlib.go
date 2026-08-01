// Copyright 2024 The zb Authors
// SPDX-License-Identifier: MIT

package xlua

import (
    "errors"
    "fmt"
    "os"
    "strings"

    rt "github.com/arnodel/golua/runtime"
    "github.com/terefang/gommons/pkg"
)

func CreateTable(r *rt.Runtime, target *rt.Table, n string) *rt.Table {
    tab := rt.NewTable()
    tv := rt.TableValue(tab)
    r.SetEnv(target, n, tv)
    return tab
}

func CreateTableStringList(r *rt.Runtime, target *rt.Table, n string, s ...string) *rt.Table {
    tab := CreateTable(r, target, n)
    //_l := len(s)
    for _i, _v := range s {
        tab.Set(rt.IntValue(int64(_i+1)), rt.StringValue(_v))
    }
    return tab
}

func TableAddStringList(r *rt.Runtime, tab *rt.Table, args []string) {
    _l := tab.Len()
    for _i, _v := range args {
        tab.Set(rt.IntValue(int64(_i+1)+_l), rt.StringValue(_v))
    }
}

func SetOsArgs(r *rt.Runtime) error {
    CreateTableStringList(r, r.GlobalEnv(), "arg", os.Args...)
    return nil
}

func SetArgs(r *rt.Runtime, args []string) error {
    CreateTableStringList(r, r.GlobalEnv(), "arg", args...)
    return nil
}

func utilStringifyEx(t *rt.Thread, _arg *rt.Value, _depth int64, _sb *strings.Builder) error {
    if _arg.IsNil() {
        _sb.WriteString("nil")
        return nil
    }

    _depth--

    if _tab, _ok := _arg.TryTable(); _ok && (_depth >= 0) {
        _sb.WriteString("{ ")
        var _k rt.Value = rt.NilValue
        var _v rt.Value = rt.NilValue

        if true {
            for _ok {
                _k, _v, _ok = _tab.Next(_k)
                if _k.IsNil() {
                    break
                }
                utilStringifyEx(t, &_k, _depth, _sb)
                _sb.WriteString(" -> ")
                utilStringifyEx(t, &_v, _depth, _sb)
                _sb.WriteString(", ")
            }
        } else {
            _l, _ := rt.IntLen(t, *_arg)
            for _i := _l; _i > 0; _i-- {
                if _depth > 0 {
                    _sb.WriteString(fmt.Sprintf("%d -> ", _i))
                    _v := _tab.Get(rt.IntValue(_i))
                    utilStringifyEx(t, &_v, _depth, _sb)
                    _sb.WriteString(", ")
                } else {
                    _str, _ := _tab.Get(rt.IntValue(_i)).ToString()
                    _sb.WriteString(fmt.Sprintf("%d -> %s, ", _i, _str))
                }
            }
        }
        _sb.WriteString("} ")
    } else if _i, _ok := _arg.TryInt(); _ok {
        _sb.WriteString(fmt.Sprintf("%d", _i))
    } else if _b, _ok := _arg.TryBool(); _ok {
        if _b {
            _sb.WriteString("true")
        } else {
            _sb.WriteString("false")
        }
    } else if _str, _ok := _arg.TryString(); _ok {
        _sb.WriteString("\"" + _str + "\"")
    } else {
        _s, _ := _arg.ToString()
        _sb.WriteString(_s)
    }
    return nil
}

func utilStringify(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
    if err := c.Check1Arg(); err != nil {
        return nil, err
    }

    var _depth int64 = 1

    if err := c.CheckNArgs(2); err == nil {
        _depth, _ = c.Arg(1).TryInt()
    }

    _sb := new(strings.Builder)
    _arg := c.Arg(0)
    err := utilStringifyEx(t, &_arg, _depth, _sb)
    if err != nil {
        return nil, err
    }

    return c.PushingNext(t.Runtime, rt.StringValue(_sb.String())), nil
}

//--- Returns an iterator that behaves opposite to @{ipairs}, i.e. it iterates
//-- over all key-value pairs that @{ipairs} does not. In other words, it skips
//-- integer keys from 1 up to the first integer key absent from the table. The
//-- order in which the keys are enumerated is not specified!

func apairs(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
    if err := c.Check1Arg(); err != nil {
        return nil, err
    }
    next := c.Next()
    _arg := c.Arg(0)
    if _tab, _ok := _arg.TryTable(); _ok {
        var _k rt.Value = rt.NilValue
        var _v rt.Value = rt.NilValue
        var _lk rt.Value = rt.NilValue

        for _ok {
            _k, _v, _ok = _tab.Next(_k)
            if _k.IsNil() || _v.IsNil() {
                break
            }
            if _k.TypeName() == "string" {
                break
            }
            _lk = _k
        }
        _next := t.GlobalEnv().Get(rt.StringValue("next"))
        t.Push1(next, _next)
        t.Push1(next, c.Arg(0))
        t.Push1(next, _lk)
        return next, nil
    }
    return nil, errors.New("not a table")
}

func RegisterGlobalFunctions(r *rt.Runtime) error {

    r.SetEnvGoFunc(r.GlobalEnv(), "stringify", utilStringify, 2, false)
    r.SetEnvGoFunc(r.GlobalEnv(), "apairs", apairs, 2, false)

    _v := r.GlobalEnv().Get(rt.StringValue("_VERSION")).AsString()
    r.SetEnv(r.GlobalEnv(), "_VERSION", rt.StringValue(_v+" extlib "+pkg.PkgVersion))

    _err := RegisterUtf8Functions(r)
    if _err != nil {
        return _err
    }

    _err = RegisterTableFunctions(r)
    if _err != nil {
        return _err
    }

    return nil
}
