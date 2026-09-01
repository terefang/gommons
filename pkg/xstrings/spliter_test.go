package xstrings

import (
	"errors"
	"reflect"
	"testing"
)

func TestSplitWsWithQuoteSet_ComprehensiveQuotes(t *testing.T) {
	// Standard opening and closing quote pairs covering ", ', `, (), [], {}, <>
	qs := []rune{'"', '\'', '`', '(', '[', '{', '<'}
	qe := []rune{'"', '\'', '`', ')', ']', '}', '>'}

	t.Run("Valid Input Cases", func(t *testing.T) {
		tests := []struct {
			name     string
			cmd      string
			expected []string
		}{
			{
				name:     "standard quotes: double, single, and backticks",
				cmd:      `cmd "double quote" 'single quote' ` + "`backtick string`",
				expected: []string{"cmd", "double quote", "single quote", "backtick string"},
			},
			{
				name:     "bracket delimiters: (), [], {}, <>",
				cmd:      `func (arg1 arg2) array[item 1] struct{field 1} tag<element 1>`,
				expected: []string{"func", "arg1 arg2", "arrayitem 1", "structfield 1", "tagelement 1"},
			},
			{
				name:     "empty quotes preserved as valid fields",
				cmd:      `cmd "" '' ` + "``" + ` () [] {} <>`,
				expected: []string{"cmd", "", "", "", "", "", "", ""},
			},
			{
				name:     "nesting asymmetric brackets inside symmetric quotes",
				cmd:      `select "user(id, name)" 'data[key]' ` + "`config{host}`",
				expected: []string{"select", "user(id, name)", "data[key]", "config{host}"},
			},
			{
				name:     "escaped opening and closing brackets",
				cmd:      `cmd \(not_a_group\) \[not_an_array\] \{not_a_struct\} \<not_a_tag\>`,
				expected: []string{"cmd", "(not_a_group)", "[not_an_array]", "{not_a_struct}", "<not_a_tag>"},
			},
			{
				name:     "escaped quotes within their own quote bounds",
				cmd:      `cmd "escaped \" double" 'escaped \' single'`,
				expected: []string{"cmd", `escaped " double`, `escaped ' single`},
			},
			{
				name:     "multiple spaces within quotes are preserved",
				cmd:      `cmd "  spaces   inside  " [   spaced   array   ]`,
				expected: []string{"cmd", "  spaces   inside  ", "   spaced   array   "},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := SplitWsWithQuoteSet(tt.cmd, qs, qe)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("got %q, want %q", got, tt.expected)
				}
			})
		}
	})

	t.Run("Error Cases", func(t *testing.T) {
		errorTests := []struct {
			name        string
			cmd         string
			qs          []rune
			qe          []rune
			expectedErr error
		}{
			{
				name:        "mismatched slice lengths",
				cmd:         `cmd "hello"`,
				qs:          []rune{'"', '\''},
				qe:          []rune{'"'},
				expectedErr: ErrMismatchedQuoteSets,
			},
			{
				name:        "unclosed double quote",
				cmd:         `cmd "unclosed string`,
				qs:          qs,
				qe:          qe,
				expectedErr: ErrUnclosedQuote,
			},
			{
				name:        "unclosed parenthesis",
				cmd:         `func (arg1 arg2`,
				qs:          qs,
				qe:          qe,
				expectedErr: ErrUnclosedQuote,
			},
			{
				name:        "unclosed angle bracket",
				cmd:         `tag <element`,
				qs:          qs,
				qe:          qe,
				expectedErr: ErrUnclosedQuote,
			},
			{
				name:        "dangling trailing backslash",
				cmd:         `cmd arg\`,
				qs:          qs,
				qe:          qe,
				expectedErr: ErrDanglingEscape,
			},
		}

		for _, tt := range errorTests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := SplitWsWithQuoteSet(tt.cmd, tt.qs, tt.qe)
				if err == nil {
					t.Fatalf("expected error containing %v, got nil", tt.expectedErr)
				}
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("got error %v, want %v", err, tt.expectedErr)
				}
			})
		}
	})
}

func TestMustSplitWsWithQuoteSet(t *testing.T) {
	qs := []rune{'"', '(', '['}
	qe := []rune{'"', ')', ']'}

	t.Run("valid execution does not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("unexpected panic: %v", r)
			}
		}()

		got := MustSplitWsWithQuoteSet(`echo "hello" (world)`, qs, qe)
		expected := []string{"echo", "hello", "world"}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got %v, want %v", got, expected)
		}
	})

	t.Run("invalid execution panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic for unclosed quote, but code executed safely")
			}
		}()

		// Should trigger panic due to unclosed quote
		_ = MustSplitWsWithQuoteSet(`echo "unclosed`, qs, qe)
	})
}
