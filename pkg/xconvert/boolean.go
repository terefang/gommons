package xconvert

import "strings"

func ToBooleanDefault(boolean any, def bool) bool {
	if boolean == nil {
		return def
	}
	return ToBoolean(boolean)
}

func ToBoolean(boolean any) bool {
	if boolean == nil {
		return false
	}
	_bool, ok := boolean.(bool)
	if ok {
		return _bool
	}
	_str, ok := boolean.(string)
	if ok {
		return ToBooleanString(_str)
	}
	_int, ok := boolean.(int64)
	if ok {
		return _int != 0
	}
	return false
}

func ToBooleanString(str string) bool {
	// any defined string is "true" unless it is a false indicator:
	// false, f, off, none, no, n, null, nul, nil, 0, <blank>
	for _, it := range AsArray("false", "f", "off", "none", "n", "no", "null", "nul", "nil", "0", "") {
		if BooleanCheckString(str, it) {
			return false
		}
	}
	return true
}

func BooleanCheckString(str string, check string) bool {
	str = strings.ToLower(strings.TrimSpace(str))
	check = strings.ToLower(strings.TrimSpace(check))
	if strings.Compare(str, check) == 0 {
		return true
	}
	return false
}
