package xmatch

import "errors"

// FilterOperator is a filter operator.
type FilterOperator int

const (
	// FilterAnd is the & operator.
	FilterAnd FilterOperator = iota

	// FilterOr is the | operator.
	FilterOr

	// FilterNot is the ! operator.
	FilterNot

	// FilterPresent is the =* operator.
	FilterPresent

	// FilterEqualityMatch is the = operator.
	FilterEqualityMatch

	// FilterSubstrings is the = operator in conjunction with a string that
	// has leading and trailing * characters.
	FilterSubstrings

	// FilterSubstringsPrefix is the = operator in conjunction with a string
	// that has a leading * character.
	FilterSubstringsPrefix

	// FilterSubstringsPostfix is the = operator in conjunction with a string
	// that has a trailing * character.
	FilterSubstringsPostfix

	// FilterGreaterOrEqual is the >= operator.
	FilterGreaterOrEqual

	// FilterLessOrEqual is the <= operator.
	FilterLessOrEqual

	// FilterApproxMatch is the ~= operator.
	FilterApproxMatch

	// FilterRxMatch is the ~~ operator.
	FilterRxMatch
)

// Filter is an LDAP-style filter string.
type Filter struct {

	// Op is the operation.
	Op FilterOperator

	// Children is a list of any sub-filters if this filter is a compound
	// filter.
	Children []*Filter

	// Left is the left operand.
	Left string

	// Right is the right operand.
	Right string
}

const (
	filterAnd               = FilterAnd
	filterOr                = FilterOr
	filterNot               = FilterNot
	filterPresent           = FilterPresent
	filterEqualityMatch     = FilterEqualityMatch
	filterSubstrings        = FilterSubstrings
	filterSubstringsPrefix  = FilterSubstringsPrefix
	filterSubstringsPostfix = FilterSubstringsPostfix
	filterGreaterOrEqual    = FilterGreaterOrEqual
	filterLessOrEqual       = FilterLessOrEqual
	filterApproxMatch       = FilterApproxMatch
	filterRxMatch           = FilterRxMatch
)

var filterMap = map[FilterOperator]string{
	filterAnd:            "And",
	filterOr:             "Or",
	filterNot:            "Not",
	filterEqualityMatch:  "Equality Match",
	filterSubstrings:     "Substrings",
	filterGreaterOrEqual: "Greater Or Equal",
	filterLessOrEqual:    "Less Or Equal",
	filterPresent:        "Present",
	filterApproxMatch:    "Approx Match",
	filterRxMatch:        "Regex Match",
}

var (
	errCharZeroNotLParen = errors.New("filter does not start with an '()'")
	errUnexpectedEOF     = errors.New("unexpected end of filter")
	errCompile           = errors.New("error compiling filter")
	errParse             = errors.New("error parsing filter")
)

type FilterAttrMapFuncValue func(key string) (string, bool)
type FilterAttrMapFuncValues func(key string) ([]string, bool)
type FilterAttrMapValue map[string]string
type FilterAttrMapValues map[string][]string

type FilterAttrMap interface {
	map[string]string | map[string][]string | FilterAttrMapValue | FilterAttrMapValues | FilterAttrMapFuncValue | FilterAttrMapFuncValues
}

type FilterAttrMapSimple interface {
	map[string]string | map[string][]string | FilterAttrMapValue | FilterAttrMapValues
}
