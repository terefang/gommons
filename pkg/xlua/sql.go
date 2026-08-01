package xlua

import (
	"database/sql"
	"fmt"

	"github.com/arnodel/golua/lib/packagelib"
	rt "github.com/arnodel/golua/runtime"
)

type sqlKeyType struct{}

var sqlKey = rt.AsValue(sqlKeyType{})

type sqlData struct {
	metatable *rt.Table
}

// LibLoader can load the utf8 lib.
var SqlLibLoader = packagelib.Loader{
	Load: loadSql,
	Name: "sql",
}

func loadSql(r *rt.Runtime) (rt.Value, func()) {
	pkg := rt.NewTable()
	r.SetEnvGoFunc(pkg, "drivers", sqlDrivers, 0, true)
	r.SetEnvGoFunc(pkg, "connect", sqlConnect, 2, true)

	methods := rt.NewTable()
	meta := rt.NewTable()
	r.SetEnv(meta, "__name", rt.StringValue("sql"))
	r.SetEnv(meta, "__index", rt.TableValue(methods))

	//	r.SetEnvGoFunc(methods, "setvbuf", filesetvbuf, 3, false)
	r.SetEnvGoFunc(meta, "exec", sql_exec, 2, true)
	r.SetEnvGoFunc(meta, "close", sql__close, 1, false)
	r.SetEnvGoFunc(meta, "__close", sql__close, 1, false)
	r.SetEnvGoFunc(meta, "__tostring", sql__tostring, 1, false)

	r.SetRegistry(sqlKey, rt.AsValue(&sqlData{
		metatable: meta,
	}))

	return rt.TableValue(pkg), nil
}

func sql_exec(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.CheckNArgs(2); err != nil {
		return nil, err
	}
	_sql, err := SqlArg(c, 0)
	if err == nil {
		_qu, _ := c.StringArg(1)
		c.Etc()
		_r, err := _sql.Exec(_qu)
		if err == nil {
			return c.PushingNext(t.Runtime, rt.IntValue(_r)), nil
		}
	}
	return c.Next(), nil
}

func sql__tostring(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	_sql, err := SqlArg(c, 0)
	if err != nil {
		return nil, err
	}
	return c.PushingNext(t.Runtime, rt.StringValue(fmt.Sprintf("sql(d=%s,`%s`)", _sql.driver, _sql.uri))), nil
}

func sql__close(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	_sql, err := SqlArg(c, 0)
	if err != nil {
		return nil, err
	}
	err = _sql.Close()
	if err != nil {
		return nil, err
	}
	return c.Next(), nil
}

func getSqlData(r *rt.Runtime) *sqlData {
	return r.Registry(sqlKey).Interface().(*sqlData)
}

func sqlDrivers(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	_list := rt.NewTable()
	TableAddStringList(t.Runtime, _list, sql.Drivers())
	return c.PushingNext1(t.Runtime, rt.AsValue(_list)), nil
}

func sqlConnect(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	_driver, _err := c.StringArg(0)
	if _err == nil {
		_uri, _err := c.StringArg(1)
		if _err == nil {
			_db, _err := sql.Open(_driver, _uri)
			if _err == nil {
				_sv := &sqlValue{db: _db, tx: nil, driver: _driver, uri: _uri}
				_ldb := t.NewUserDataValue(_sv, getSqlData(t.Runtime).metatable)
				return c.PushingNext1(t.Runtime, _ldb), nil
			}
			return nil, _err
		}
		return nil, _err
	}
	return nil, _err
}

func SqlArg(c *rt.GoCont, n int) (*sqlValue, error) {
	f, ok := ValueToSql(c.Arg(n))
	if ok {
		return f, nil
	}
	return nil, fmt.Errorf("#%d must be a sql", n+1)
}

func ValueToSql(v rt.Value) (*sqlValue, bool) {
	u, ok := v.TryUserData()
	if ok {
		return u.Value().(*sqlValue), true
	}
	return nil, false
}
