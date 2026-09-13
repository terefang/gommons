/*
Package provides Lua-style pattern matching for Go byte slices.

PATTERN SYNTAX REFERENCE

1. Character Classes

A character class matches any single byte belonging to the designated set:

	.    Matches any character.
	%a   Alphabetic letters (a-z, A-Z).
	%c   Control characters (ASCII 0x00-0x1F and 0x7F).
	%d   Digits (0-9).
	%l   Lowercase letters (a-z).
	%p   Punctuation characters (!, ", #, $, %, &, ', (, ), etc.).
	%s   Space/whitespace characters (' ', '\f', '\n', '\r', '\t', '\v').
	%u   Uppercase letters (A-Z).
	%w   Alphanumeric characters (0-9, a-z, A-Z).
	%x   Hexadecimal digits (0-9, a-f, A-F).
	%z   The null byte (0x00).
	%x   (Non-alphanumeric 'x') Escapes special characters to match them
	     literally (e.g., %% matches %, %. matches .).

Capitalizing the class letter negates the set (e.g., %D matches any non-digit
character, %S matches any non-whitespace character).

2. Custom Sets & Ranges ([...])

Custom sets match single characters against user-defined lists or ranges:

	[abc]   Matches any character inside the set.
	[^abc]  Matches any character NOT inside the set.
	[a-z]   Matches any character between two values inclusive.
	[%d%s]  Character classes can be embedded inside sets.

3. Quantifiers (Repetition)

Quantifiers apply to the single character, class, or custom set directly
preceding them:

	*   0 or more repetitions (Greedy: matches as many as possible).
	+   1 or more repetitions (Greedy: matches as many as possible).
	-   0 or more repetitions (Lazy: matches as few as possible).
	?   0 or 1 optional occurrence (Greedy).

4. Anchors

	^   Constrains match to start at the beginning of the input string
	    (must be the first character in the pattern).
	$   Constrains match to end at the end of the input string
	    (must be the last character in the pattern).

5. Captures & Backreferences

	(...)   Group Capture: Captures the matched substring inside parentheses
	        into the MatchData output.
	()      Position Capture: Captures the current 1-based index position
	        in the input string without consuming characters.
	%1-%9   Backreference: Matches a substring identical to the string
	        captured by the Nth capture group earlier in the pattern.

6. Balanced Pair Matching (%b)

	%bxy    Matches balanced string pairs starting with byte 'x' and ending
	        with byte 'y'. Correctly handles nested balanced pairs
	        (e.g., %b() matches "(a + (b * c))").

LIMITATIONS

The following standard regex features are NOT supported by design:
  - Alternation / OR operations (|)
  - Explicit repetition counts ({n,m})
  - Lookaround assertions ((?=...), (?<=...))
*/

package xmatch

const maxCaptures = 32

// LLMatchData holds the result of a successful pattern match.
type LLMatchData struct {
	Subject  []byte
	Start    int // 0-based start index in Subject
	End      int // 0-based end index (exclusive) in Subject
	captures [maxCaptures]struct {
		start int // 0-based start index (-1 if unset)
		len   int // length of capture (-1 for position capture)
	}
	numCaptures int
}

// NumCaptures returns the number of explicit captures found.
func (m *LLMatchData) NumCaptures() int {
	return m.numCaptures
}

// GetCapture returns the matched slice or 1-based position for index (1-based).
func (m *LLMatchData) GetCapture(i int) ([]byte, int, bool) {
	if i < 1 || i > m.numCaptures {
		return nil, 0, false
	}
	cap := m.captures[i-1]
	if cap.len == -1 {
		// Position capture (1-based index)
		return nil, cap.start + 1, true
	}
	return m.Subject[cap.start : cap.start+cap.len], 0, true
}

type LLMatchState struct {
	src         []byte
	pat         []byte
	level       int // recursion depth
	numCaptures int
	capture     [maxCaptures]struct {
		start int
		len   int
	}
}

// LLFind searches for the first occurrence of pattern in src starting at offset.
func LLFind(pattern string, src []byte, offset int) (*LLMatchData, bool) {
	if offset < 0 || offset > len(src) {
		return nil, false
	}

	pat := []byte(pattern)
	anchor := len(pat) > 0 && pat[0] == '^'
	patIdx := 0
	if anchor {
		patIdx = 1
	}

	ms := &LLMatchState{
		src: src,
		pat: pat,
	}

	s := offset
	for {
		ms.numCaptures = 0
		end := ms.match(s, patIdx)
		if end != -1 {
			res := &LLMatchData{
				Subject:     src,
				Start:       s,
				End:         end,
				numCaptures: ms.numCaptures,
			}
			copy(res.captures[:], ms.capture[:])
			return res, true
		}
		s++
		if s > len(src) || anchor {
			break
		}
	}

	return nil, false
}

// match is the core recursive matching engine.
func (ms *LLMatchState) match(s, p int) int {
	ms.level++
	if ms.level > 2000 {
		return -1 // prevent stack overflow
	}
	defer func() { ms.level-- }()

	for p < len(ms.pat) {
		switch ms.pat[p] {
		case '(':
			if p+1 < len(ms.pat) && ms.pat[p+1] == ')' {
				return ms.matchPosCapture(s, p+2)
			}
			return ms.matchStartCapture(s, p+1)

		case ')':
			return ms.matchEndCapture(s, p+1)

		case '$':
			if p+1 == len(ms.pat) { // end anchor
				if s == len(ms.src) {
					return s
				}
				return -1
			}
			// Treat $ as literal if not at end of pattern
			goto checkQuantifier

		case '%':
			if p+1 < len(ms.pat) {
				ch := ms.pat[p+1]
				if ch >= '1' && ch <= '9' {
					s = ms.matchBackref(s, int(ch-'1'))
					if s == -1 {
						return -1
					}
					p += 2
					continue
				}
				if ch == 'b' {
					if p+3 < len(ms.pat) {
						s = ms.matchBalanced(s, ms.pat[p+2], ms.pat[p+3])
						if s == -1 {
							return -1
						}
						p += 4
						continue
					}
					return -1
				}
			}
			goto checkQuantifier

		default:
			goto checkQuantifier
		}

	checkQuantifier:
		pNext := ms.classEnd(p)
		var q byte = 0
		if pNext < len(ms.pat) {
			q = ms.pat[pNext]
		}

		switch q {
		case '*': // Greedy: 0 or more
			return ms.maxExpand(s, p, pNext+1)
		case '+': // Greedy: 1 or more
			if s < len(ms.src) && ms.matchClass(ms.src[s], p) {
				return ms.maxExpand(s+1, p, pNext+1)
			}
			return -1
		case '-': // Lazy: 0 or more
			return ms.minExpand(s, p, pNext+1)
		case '?': // Optional: 0 or 1
			if s < len(ms.src) && ms.matchClass(ms.src[s], p) {
				if res := ms.match(s+1, pNext+1); res != -1 {
					return res
				}
			}
			p = pNext + 1 // try matching 0 occurrences
		default:
			if s < len(ms.src) && ms.matchClass(ms.src[s], p) {
				s++
				p = pNext
			} else {
				return -1
			}
		}
	}
	return s
}

// maxExpand handles greedy quantifiers (* and +).
func (ms *LLMatchState) maxExpand(s, p, pNext int) int {
	i := 0
	for s+i < len(ms.src) && ms.matchClass(ms.src[s+i], p) {
		i++
	}
	// Backtrack from maximum matches to minimum
	for i >= 0 {
		if res := ms.match(s+i, pNext); res != -1 {
			return res
		}
		i--
	}
	return -1
}

// minExpand handles lazy quantifier (-).
func (ms *LLMatchState) minExpand(s, p, pNext int) int {
	for {
		if res := ms.match(s, pNext); res != -1 {
			return res
		}
		if s < len(ms.src) && ms.matchClass(ms.src[s], p) {
			s++
		} else {
			break
		}
	}
	return -1
}

// matchClass evaluates character class matching at pat[p].
func (ms *LLMatchState) matchClass(c byte, p int) bool {
	switch ms.pat[p] {
	case '.':
		return true
	case '%':
		if p+1 >= len(ms.pat) {
			return false
		}
		return LLMatchSingleClass(c, ms.pat[p+1])
	case '[':
		return ms.matchSetClass(c, p)
	default:
		return ms.pat[p] == c
	}
}

// LLMatchSingleClass tests a character against a single class specifier.
func LLMatchSingleClass(c, cl byte) bool {
	low := cl
	if cl >= 'A' && cl <= 'Z' {
		low = cl + ('a' - 'A')
	}
	res := false
	switch low {
	case 'a':
		res = (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
	case 'c':
		res = (c <= 0x1F) || c == 0x7F
	case 'd':
		res = c >= '0' && c <= '9'
	case 'l':
		res = c >= 'a' && c <= 'z'
	case 'p':
		res = (c >= '!' && c <= '/') || (c >= ':' && c <= '@') || (c >= '[' && c <= '`') || (c >= '{' && c <= '~')
	case 's':
		res = c == ' ' || c == '\f' || c == '\n' || c == '\r' || c == '\t' || c == '\v'
	case 'u':
		res = c >= 'A' && c <= 'Z'
	case 'w':
		res = (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
	case 'x':
		res = (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
	case 'z':
		res = c == 00
	default:
		return c == cl
	}
	if cl >= 'A' && cl <= 'Z' {
		return !res
	}
	return res
}

func (ms *LLMatchState) matchSetClass(c byte, p int) bool {
	sig := true
	p++
	if p < len(ms.pat) && ms.pat[p] == '^' {
		sig = false
		p++
	}
	matched := false
	for p < len(ms.pat) && ms.pat[p] != ']' {
		if ms.pat[p] == '%' && p+1 < len(ms.pat) {
			if LLMatchSingleClass(c, ms.pat[p+1]) {
				matched = true
			}
			p += 2
			continue
		}
		if p+2 < len(ms.pat) && ms.pat[p+1] == '-' && ms.pat[p+2] != ']' {
			if c >= ms.pat[p] && c <= ms.pat[p+2] {
				matched = true
			}
			p += 3
			continue
		}
		if ms.pat[p] == c {
			matched = true
		}
		p++
	}
	return matched == sig
}

func (ms *LLMatchState) classEnd(p int) int {
	switch ms.pat[p] {
	case '%':
		if p+1 < len(ms.pat) {
			return p + 2
		}
		return p + 1
	case '[':
		p++
		if p < len(ms.pat) && ms.pat[p] == '^' {
			p++
		}
		for p < len(ms.pat) {
			if ms.pat[p] == '%' && p+1 < len(ms.pat) {
				p += 2
				continue
			}
			if ms.pat[p] == ']' {
				return p + 1
			}
			p++
		}
		return p
	default:
		return p + 1
	}
}

func (ms *LLMatchState) matchStartCapture(s, p int) int {
	if ms.numCaptures >= maxCaptures {
		return -1
	}
	capIdx := ms.numCaptures
	ms.numCaptures++
	ms.capture[capIdx].start = s
	ms.capture[capIdx].len = -1

	res := ms.match(s, p)
	if res == -1 {
		ms.numCaptures--
	}
	return res
}

func (ms *LLMatchState) matchEndCapture(s, p int) int {
	for i := ms.numCaptures - 1; i >= 0; i-- {
		if ms.capture[i].len == -1 {
			ms.capture[i].len = s - ms.capture[i].start
			res := ms.match(s, p)
			if res == -1 {
				ms.capture[i].len = -1
			}
			return res
		}
	}
	return -1
}

func (ms *LLMatchState) matchPosCapture(s, p int) int {
	if ms.numCaptures >= maxCaptures {
		return -1
	}
	capIdx := ms.numCaptures
	ms.numCaptures++
	ms.capture[capIdx].start = s
	ms.capture[capIdx].len = -1 // -1 signifies position capture

	res := ms.match(s, p)
	if res == -1 {
		ms.numCaptures--
	}
	return res
}

func (ms *LLMatchState) matchBackref(s, capIdx int) int {
	if capIdx >= ms.numCaptures || ms.capture[capIdx].len < 0 {
		return -1
	}
	cLen := ms.capture[capIdx].len
	cStart := ms.capture[capIdx].start

	if len(ms.src)-s < cLen {
		return -1
	}
	for i := 0; i < cLen; i++ {
		if ms.src[s+i] != ms.src[cStart+i] {
			return -1
		}
	}
	return s + cLen
}

func (ms *LLMatchState) matchBalanced(s int, open, close byte) int {
	if s >= len(ms.src) || ms.src[s] != open {
		return -1
	}
	count := 1
	s++
	for s < len(ms.src) {
		if ms.src[s] == close {
			count--
			if count == 0 {
				return s + 1
			}
		} else if ms.src[s] == open {
			count++
		}
		s++
	}
	return -1
}
