package xstrings

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cast"
)

func Stringify(obj interface{}) (string, error) {
	return StringifyWithLevel(obj, 3)
}

func StringifyWithLevel(obj interface{}, level int) (string, error) {
	_sb := new(strings.Builder)
	_safe := make(map[uintptr]bool)
	_err := doStringify(obj, _sb, &_safe, level)
	if _err != nil {
		return "nil", _err
	}
	return _sb.String(), nil
}

func StringifyNoError(obj interface{}) string {
	_str, _err := Stringify(obj)
	if _err == nil {
		return _str
	}
	return _err.Error()
}

func StringifyNoErrorWithLevel(obj interface{}, level int) string {
	_str, _err := StringifyWithLevel(obj, level)
	if _err == nil {
		return _str
	}
	return _err.Error()
}

func doStringify(obj interface{}, sb *strings.Builder, safe *map[uintptr]bool, level int) error {
	v := reflect.ValueOf(obj)
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
		_ref := uintptr(reflect.ValueOf(obj).Pointer())
		_bool, _ok := (*safe)[_ref]
		if _ok {
			if _bool {
				sb.WriteString(fmt.Sprintf("(ref=0x%08x)", _ref))
				return nil
			}
		}
		(*safe)[_ref] = true
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		if level == 0 {
			sb.WriteString("<obj>")
			return nil
		}
		return doStringifyStruct(obj, v, sb, safe, level)
	case reflect.Map:
		if level == 0 {
			sb.WriteString("{map}")
			return nil
		}
		return doStringifyMap(obj, v, sb, safe, level)
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		if level == 0 {
			sb.WriteString("[arr]")
			return nil
		}
		return doStringifyArray(obj, v, sb, safe, level)
	case reflect.String:
		if level == 0 {
			sb.WriteString("(str)")
			return nil
		}
		sb.WriteRune('`')
		str := v.Interface().(string)
		sb.WriteString(str)
		sb.WriteRune('`')
		return nil
	default:
		//if level == 0 {
		//			sb.WriteString("(val)")
		//			return nil
		//		}
		sb.WriteString(cast.ToString(v.Interface()))
		return nil
	}
}

func doStringifyArray(obj interface{}, v reflect.Value, sb *strings.Builder, safe *map[uintptr]bool, level int) error {
	level--
	// we already checked for safe in upstream#
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	len := v.Len()
	sb.WriteRune('[')
	first := true
	for i := 0; i < len; i++ {
		value := v.Index(i)
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		if !first {
			sb.WriteRune(',')
		}
		first = false
		doStringify(value.Interface(), sb, safe, level)
	}
	sb.WriteRune(']')
	return nil
}

func doStringifyMap(obj interface{}, v reflect.Value, sb *strings.Builder, safe *map[uintptr]bool, level int) error {
	level--
	// we already checked for safe in upstream
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	first := true
	sb.WriteRune('{')
	for _, index := range v.MapKeys() {
		value := v.MapIndex(index)
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		if !first {
			sb.WriteRune(',')
		}
		first = false
		key := cast.ToString(index.Interface())
		sb.WriteString(`"`)
		sb.WriteString(key)
		sb.WriteString(`"=`)
		doStringify(value.Interface(), sb, safe, level)
	}
	sb.WriteRune('}')
	return nil
}

func doStringifyStruct(obj interface{}, v reflect.Value, sb *strings.Builder, safe *map[uintptr]bool, level int) error {
	level--
	// we already checked for safe in upstream
	t := v.Type()
	len := t.NumField()
	sb.WriteRune('<')
	first := true
	for i := 0; i < len; i++ {

		field := t.Field(i)
		// we can't access the value of unexported fields
		if field.PkgPath != "" {
			continue
		}
		key := field.Name
		value := v.FieldByIndex(field.Index)

		// If it is a pointer, get the true value
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}

		if !first {
			sb.WriteRune(',')
		}
		//sb.WriteString(`"`)
		sb.WriteString(key)
		//sb.WriteString(`":`)
		sb.WriteRune(':')
		doStringify(value.Interface(), sb, safe, level)
		first = false
	}
	sb.WriteRune('>')
	return nil
}
