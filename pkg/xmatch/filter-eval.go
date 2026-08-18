package xmatch

import (
	"errors"
	"regexp"
	"strings"
)

func FilterEvaluateKV[V FilterAttrMap](f *Filter, kv V) (bool, error) {
	// no filter no match
	if f == nil {
		return false, nil
	}
	// no values no match
	if kv == nil {
		return false, nil
	}

	switch f.Op {
	case filterAnd:
		return FilterEvaluateKvAnd(f, kv)
	case filterOr:
		return FilterEvaluateKvOr(f, kv)
	case filterNot:
		return FilterEvaluateKvNot(f, kv)
	case filterEqualityMatch:
		return FilterEvaluateKvEq(f, kv)
	case filterSubstrings:
		return FilterEvaluateKvContains(f, kv)
	case filterGreaterOrEqual:
		return FilterEvaluateKvGte(f, kv)
	case filterLessOrEqual:
		return FilterEvaluateKvLte(f, kv)
	case filterPresent:
		return FilterEvaluateKvExists(f, kv)
	case filterApproxMatch:
		//return FilterEvaluateKvApprox(f, kv)
		return FilterEvaluateKvContains(f, kv)
	case filterRxMatch:
		return FilterEvaluateKvRx(f, kv)
	}
	return false, errors.New("unknown filter op")
}

func FilterEvaluateKvAnd[V FilterAttrMap](f *Filter, kv V) (bool, error) {
	for _, _f := range f.Children {
		_ok, err := FilterEvaluateKV(_f, kv)
		if err != nil {
			return false, err
		}
		if !_ok {
			return false, nil
		}
	}
	return true, nil
}

func FilterEvaluateKvOr[V FilterAttrMap](f *Filter, kv V) (bool, error) {
	for _, _f := range f.Children {
		_ok, err := FilterEvaluateKV(_f, kv)
		if err != nil {
			return false, err
		}
		if _ok {
			return true, nil
		}
	}
	return false, nil
}

func FilterEvaluateKvNot[V FilterAttrMap](f *Filter, kv V) (bool, error) {
	_ok, err := FilterEvaluateKV(f.Children[0], kv)
	return !_ok, err
}

func FilterEvaluateKvExists[V FilterAttrMap](f *Filter, kv V) (bool, error) {
	var _kvi interface{} = kv
	switch _kv := _kvi.(type) {
	case nil:
		return false, nil
	case FilterAttrMapValues:
		_, _ok := _kv[f.Left]
		if _ok {
			return true, nil
		}
		return false, nil
	case FilterAttrMapValue:
		_, _ok := _kv[f.Left]
		if _ok {
			return true, nil
		}
		return false, nil
	case FilterAttrMapFuncValues:
		_, _ok := _kv(f.Left)
		if _ok {
			return true, nil
		}
		return false, nil
	case FilterAttrMapFuncValue:
		_, _ok := _kv(f.Left)
		if _ok {
			return true, nil
		}
		return false, nil
	}
	return false, nil
}

func ConvertAttributeFromKey[V FilterAttrMap](key string, kv V) ([]string, bool) {
	var _values []string
	var _ok bool

	var _kvi interface{} = kv
	switch _kv := _kvi.(type) {
	case nil:
		return nil, false
	case FilterAttrMapValues:
		_values, _ok = _kv[key]
		return _values, _ok
	case FilterAttrMapValue:
		_v, _ok := _kv[key]
		if !_ok {
			return nil, _ok
		}
		return []string{_v}, _ok
	case FilterAttrMapFuncValues:
		_values, _ok = _kv(key)
		return _values, _ok
	case FilterAttrMapFuncValue:
		_v, _ok := _kv(key)
		if !_ok {
			return nil, _ok
		}
		return []string{_v}, _ok
	default:
		return nil, false
	}
}
func FilterEvaluateKvEq[V FilterAttrMap](f *Filter, kv V) (bool, error) {
	_values, _ok := ConvertAttributeFromKey[V](f.Left, kv)
	if !_ok {
		return _ok, nil
	}

	for _, _v := range _values {
		if strings.EqualFold(_v, f.Right) {
			return true, nil
		}
	}
	return false, nil
}

func FilterEvaluateKvContains[V FilterAttrMap](f *Filter, kv V) (bool, error) {
	_values, _ok := ConvertAttributeFromKey[V](f.Left, kv)
	if !_ok {
		return _ok, nil
	}

	for _, _v := range _values {
		if strings.Contains(strings.ToLower(_v), strings.ToLower(f.Right)) {
			return true, nil
		}
	}
	return false, nil
}

func FilterEvaluateKvRx[V FilterAttrMap](f *Filter, kv V) (bool, error) {
	_values, _ok := ConvertAttributeFromKey[V](f.Left, kv)
	if !_ok {
		return _ok, nil
	}

	for _, _v := range _values {
		_ok, err := regexp.MatchString(f.Right, strings.ToLower(_v))
		if err != nil {
			return false, err
		}

		if _ok {
			return true, nil
		}
	}
	return false, nil
}

func FilterEvaluateKvGte[V FilterAttrMap](f *Filter, kv V) (bool, error) {
	_values, _ok := ConvertAttributeFromKey[V](f.Left, kv)
	if !_ok {
		return _ok, nil
	}

	for _, _v := range _values {
		if strings.Compare(strings.ToLower(_v), strings.ToLower(f.Right)) >= 0 {
			return true, nil
		}
	}
	return false, nil
}

func FilterEvaluateKvLte[V FilterAttrMap](f *Filter, kv V) (bool, error) {
	_values, _ok := ConvertAttributeFromKey[V](f.Left, kv)
	if !_ok {
		return _ok, nil
	}

	for _, _v := range _values {
		if strings.Compare(strings.ToLower(_v), strings.ToLower(f.Right)) <= 0 {
			return true, nil
		}
	}
	return false, nil
}
