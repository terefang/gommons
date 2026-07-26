package xlua

import (
	"reflect"

	rt "github.com/arnodel/golua/runtime"
)

func ToLuaValue(r rt.Runtime, obj interface{}) (rt.Value, error) {
	return ToLuaValueWithLevel(r, obj, 3)
}

func ToLuaValueWithLevel(r rt.Runtime, obj interface{}, level int) (rt.Value, error) {

	_safe := make(map[uintptr]*rt.Value)
	_val, _err := doToLuaValue(r, obj, &_safe, level)
	if _err != nil {
		return rt.NilValue, _err
	}
	return _val, nil
}

func ToLuaValueNoError(r rt.Runtime, obj interface{}) rt.Value {
	_val, _err := ToLuaValue(r, obj)
	if _err == nil {
		return _val
	}
	return rt.NilValue
}

func ToLuaValueNoErrorWithLevel(r rt.Runtime, obj interface{}, level int) rt.Value {
	_val, _err := ToLuaValueWithLevel(r, obj, level)
	if _err == nil {
		return _val
	}
	return rt.NilValue
}

func doToLuaValue(r rt.Runtime, obj interface{}, safe *map[uintptr]*rt.Value, level int) (rt.Value, error) {
	if obj == nil {
		return rt.NilValue, nil
	}

	v := reflect.ValueOf(obj)
	var _ref uintptr = 0
	// check for recursive cycles
	switch v.Kind() {
	case reflect.Ptr:
		fallthrough
	case reflect.Struct:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		_ref = reflect.ValueOf(obj).Pointer()
		_val, _ok := (*safe)[_ref]
		if _ok {
			return *_val, nil
		}
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var _lval rt.Value
	var _err error = nil

	switch v.Kind() {
	case reflect.Struct:
		_lval, _err = doToLuaValueStruct(r, obj, v, safe, level)
	case reflect.Map:
		_lval, _err = doToLuaValueMap(r, obj, v, safe, level)
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		_lval, _err = doToLuaValueArray(r, obj, v, safe, level)
	default:
		_lval = rt.AsValue(obj)
	}
	(*safe)[_ref] = &_lval
	return _lval, _err
}

func doToLuaValueArray(r rt.Runtime, obj interface{}, v reflect.Value, safe *map[uintptr]*rt.Value, level int) (rt.Value, error) {
	if level == 0 {
		return rt.AsValue(rt.NewTable()), nil
	}
	level--

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	len := v.Len()
	_t := rt.NewTable()
	for i := 0; i < len; i++ {
		value := v.Index(i)
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		_val, _err := doToLuaValue(r, value.Interface(), safe, level)
		if _err != nil {
			_t.Set(rt.AsValue(i+1), rt.NilValue)
			continue
		}
		_t.Set(rt.AsValue(i+1), _val)
	}
	return rt.AsValue(_t), nil
}

func doToLuaValueMap(r rt.Runtime, obj interface{}, v reflect.Value, safe *map[uintptr]*rt.Value, level int) (rt.Value, error) {
	if level == 0 {
		return rt.AsValue(rt.NewTable()), nil
	}
	level--

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	_t := rt.NewTable()
	for _, index := range v.MapKeys() {
		value := v.MapIndex(index)
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		_key := index.String()
		_val, _err := doToLuaValue(r, value.Interface(), safe, level)
		if _err != nil {
			//_t.Set(rt.AsValue(index), rt.AsValue(_err.Error()))
			_t.Set(rt.StringValue(_key), rt.NilValue)
			continue
		}
		_t.Set(rt.StringValue(_key), _val)
	}
	return rt.AsValue(_t), nil
}

func doToLuaValueStruct(r rt.Runtime, obj interface{}, v reflect.Value, safe *map[uintptr]*rt.Value, level int) (rt.Value, error) {
	//if level == 0 {
	//	return rt.AsValue(rt.NewTable()), nil
	//}
	//level--

	_lval := r.NewUserDataValue(obj, rt.NewTable())

	return _lval, nil
}
