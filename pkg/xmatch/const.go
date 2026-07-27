package xmatch

// and or eq lt gt lte gte (like regex contains startwith endswith wildcard)

type Matcher interface {
	Match(val any) bool
}

type AndMatcher struct {
	Matchers []Matcher
}

func (a AndMatcher) Match(val any) bool {
	if len(a.Matchers) == 0 {
		return false
	}
	for _, m := range a.Matchers {
		if !m.Match(val) {
			return false
		}
	}
	return true
}

type OrMatcher struct {
	Matchers []Matcher
}

func (a OrMatcher) Match(val any) bool {
	if len(a.Matchers) == 0 {
		return false
	}
	for _, m := range a.Matchers {
		if m.Match(val) {
			return true
		}
	}
	return false
}

type NotMatcher struct {
	Nmatcher Matcher
}

func (a NotMatcher) Match(val any) bool {
	return !a.Nmatcher.Match(val)
}

type FalseMatcher struct{}

func (a FalseMatcher) Match(val any) bool {
	return false
}

type TrueMatcher struct{}

func (a TrueMatcher) Match(val any) bool {
	return true
}
