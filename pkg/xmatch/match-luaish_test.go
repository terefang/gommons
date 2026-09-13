package xmatch

import (
	"bytes"
	"testing"
)

func TestLLFind(t *testing.T) {
	tests := []struct {
		name          string
		pattern       string
		input         string
		offset        int
		wantMatched   bool
		wantStart     int
		wantEnd       int
		wantCaptures  []string // String representation or positions
		isPosCaptures []bool   // Identifies if capture is a position capture
	}{
		// 1. Literal & Anchors
		{
			name:        "Literal match",
			pattern:     "hello",
			input:       "say hello world",
			offset:      0,
			wantMatched: true,
			wantStart:   4,
			wantEnd:     9,
		},
		{
			name:        "Start Anchor ^ success",
			pattern:     "^say",
			input:       "say hello",
			offset:      0,
			wantMatched: true,
			wantStart:   0,
			wantEnd:     3,
		},
		{
			name:        "Start Anchor ^ failure",
			pattern:     "^hello",
			input:       "say hello",
			offset:      0,
			wantMatched: false,
		},
		{
			name:        "End Anchor $ success",
			pattern:     "world$",
			input:       "hello world",
			offset:      0,
			wantMatched: true,
			wantStart:   6,
			wantEnd:     11,
		},

		// 2. Character Classes
		{
			name:        "Digit %d and Alpha %a",
			pattern:     "%d+%s+%a+",
			input:       "item 1234  abcd",
			offset:      0,
			wantMatched: true,
			wantStart:   5,
			wantEnd:     15,
		},
		{
			name:        "Negated Digit %D",
			pattern:     "%D+",
			input:       "123abc456",
			offset:      0,
			wantMatched: true,
			wantStart:   3,
			wantEnd:     6,
		},
		{
			name:        "Hex Digit %x",
			pattern:     "%x+",
			input:       "Code: 0xFF12",
			offset:      0,
			wantMatched: true,
			wantStart:   0,
			wantEnd:     1,
		},
		{
			name:        "Hex Digit %x 2",
			pattern:     "0x%x+",
			input:       "Code: 0xFF12",
			offset:      0,
			wantMatched: true,
			wantStart:   6,
			wantEnd:     12,
		},

		// 3. Custom Sets & Ranges
		{
			name:        "Custom set range [a-z]",
			pattern:     "[a-z]+",
			input:       "123 abc 456",
			offset:      0,
			wantMatched: true,
			wantStart:   4,
			wantEnd:     7,
		},
		{
			name:        "Negated custom set [^0-9]",
			pattern:     "[^0-9]+",
			input:       "123abc456",
			offset:      0,
			wantMatched: true,
			wantStart:   3,
			wantEnd:     6,
		},

		// 4. Quantifiers (*, +, -, ?)
		{
			name:        "Greedy quantifier *",
			pattern:     "a*",
			input:       "baaaac",
			offset:      0,
			wantMatched: true,
			wantStart:   0,
			wantEnd:     0, // matches empty string at index 0 before 'b'
		},
		{
			name:        "Greedy quantifier +",
			pattern:     "a+",
			input:       "baaaac",
			offset:      0,
			wantMatched: true,
			wantStart:   1,
			wantEnd:     5,
		},
		{
			name:        "Lazy quantifier -",
			pattern:     "a-",
			input:       "baaaac",
			offset:      1,
			wantMatched: true,
			wantStart:   1,
			wantEnd:     1, // lazy match consumes minimum (0) 'a's
		},
		{
			name:        "Lazy quantifier - in context",
			pattern:     "<.->",
			input:       "<div>content</div>",
			offset:      0,
			wantMatched: true,
			wantStart:   0,
			wantEnd:     5, // matches "<div>" rather than the entire string
		},
		{
			name:        "Optional quantifier ?",
			pattern:     "https?://",
			input:       "visit http://example.com",
			offset:      0,
			wantMatched: true,
			wantStart:   6,
			wantEnd:     13,
		},

		// 5. Captures
		{
			name:          "Group capture",
			pattern:       "(%a+)%s+(%d+)",
			input:         "age 30",
			offset:        0,
			wantMatched:   true,
			wantStart:     0,
			wantEnd:       6,
			wantCaptures:  []string{"age", "30"},
			isPosCaptures: []bool{false, false},
		},
		{
			name:          "Position capture ()",
			pattern:       "()hello()",
			input:         "abc hello def",
			offset:        0,
			wantMatched:   true,
			wantStart:     4,
			wantEnd:       9,
			wantCaptures:  []string{"5", "10"}, // 1-based indices
			isPosCaptures: []bool{true, true},
		},
		{
			name:          "Backreference %1",
			pattern:       "([\"'])(.-)%1",
			input:         `he said "hello world" to me`,
			offset:        0,
			wantMatched:   true,
			wantStart:     8,
			wantEnd:       21,
			wantCaptures:  []string{`"`, `hello world`},
			isPosCaptures: []bool{false, false},
		},

		// 6. Balanced Pair Matching (%b)
		{
			name:        "Balanced pair %b()",
			pattern:     "%b()",
			input:       "a + (b * (c + d)) + e",
			offset:      0,
			wantMatched: true,
			wantStart:   4,
			wantEnd:     17, // "(b * (c + d))"
		},
		{
			name:        "Balanced pair %b<>",
			pattern:     "%b<>",
			input:       "foo <bar <baz>> qux",
			offset:      0,
			wantMatched: true,
			wantStart:   4,
			wantEnd:     15, // "<bar <baz>>"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, matched := LLFind(tt.pattern, []byte(tt.input), tt.offset)

			if matched != tt.wantMatched {
				t.Fatalf("LLFind() matched = %v, wantMatched %v", matched, tt.wantMatched)
			}

			if !matched {
				return
			}

			if res.Start != tt.wantStart || res.End != tt.wantEnd {
				t.Errorf("LLFind() range = [%d, %d), want [%d, %d)", res.Start, res.End, tt.wantStart, tt.wantEnd)
			}

			if len(tt.wantCaptures) > 0 {
				if res.NumCaptures() != len(tt.wantCaptures) {
					t.Fatalf("LLFind() NumCaptures = %d, want %d", res.NumCaptures(), len(tt.wantCaptures))
				}

				for i := 1; i <= res.NumCaptures(); i++ {
					capBytes, pos, ok := res.GetCapture(i)
					if !ok {
						t.Errorf("GetCapture(%d) failed", i)
						continue
					}

					if tt.isPosCaptures[i-1] {
						if pos != parsePos(tt.wantCaptures[i-1]) {
							t.Errorf("GetCapture(%d) pos = %d, want %s", i, pos, tt.wantCaptures[i-1])
						}
					} else {
						if !bytes.Equal(capBytes, []byte(tt.wantCaptures[i-1])) {
							t.Errorf("GetCapture(%d) bytes = %s, want %s", i, string(capBytes), tt.wantCaptures[i-1])
						}
					}
				}
			}
		})
	}
}

func parsePos(s string) int {
	var pos int
	for _, c := range s {
		pos = pos*10 + int(c-'0')
	}
	return pos
}
