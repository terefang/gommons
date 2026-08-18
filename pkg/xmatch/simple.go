package xmatch

import (
	"errors"
	"maps"
	"regexp"
	"slices"
	"strings"
)

func ConvertAttrSimpleToKeys[V FilterAttrMapSimple](kv V) ([]string, bool) {
	var _kvi interface{} = kv
	switch _kv := _kvi.(type) {
	case nil:
		return nil, false
	case map[string][]string:
		return slices.Collect(maps.Keys(_kv)), true
	case FilterAttrMapValues:
		return slices.Collect(maps.Keys(_kv)), true
	case map[string]string:
		return slices.Collect(maps.Keys(_kv)), true
	case FilterAttrMapValue:
		return slices.Collect(maps.Keys(_kv)), true
	default:
		return nil, false
	}
}

func MatchSimpleOp(value string, ops FilterOperator, matcher string) bool {
	switch ops {
	case filterGreaterOrEqual:
		fallthrough
	case filterLessOrEqual:
		fallthrough
	case filterEqualityMatch:
		return strings.EqualFold(value, matcher)
	case filterApproxMatch:
		fallthrough
	case filterSubstrings:
		return strings.Contains(strings.ToLower(value), strings.ToLower(matcher))
	case filterPresent:
		return true
	case filterRxMatch:
		_ok, err := regexp.MatchString(matcher, strings.ToLower(value))
		if err != nil {
			return false
		}
		return _ok
	}
	return false
}

func MatchSimpleOpVforK[V FilterAttrMapSimple](kv V, ops FilterOperator, matcher string) ([]string, error) {
	_keys, _ok := ConvertAttrSimpleToKeys[V](kv)
	if !_ok {
		return nil, errors.New("error keys")
	}
	var _ret []string = make([]string, 0)
	for _, _key := range _keys {
		_values, _ok := ConvertAttributeFromKey[V](_key, kv)
		if !_ok {
			continue
		}
		for _, _v := range _values {
			if MatchSimpleOp(_v, ops, matcher) {
				_ret = append(_ret, _key)
				break
			}
		}
	}
	return _ret, nil
}

func MatchSimpleEqVforK[V FilterAttrMapSimple](kv V, matcher string) ([]string, error) {
	return MatchSimpleOpVforK[V](kv, filterEqualityMatch, matcher)
}

func MatchSimpleContainsVforK[V FilterAttrMapSimple](kv V, matcher string) ([]string, error) {
	return MatchSimpleOpVforK[V](kv, filterSubstrings, matcher)
}

func MatchSimpleRxVforK[V FilterAttrMapSimple](kv V, matcher string) ([]string, error) {
	return MatchSimpleOpVforK[V](kv, filterRxMatch, matcher)
}
