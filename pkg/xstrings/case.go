package xstrings

import (
	"strings"
	"unicode"
)

// ToTitle converts the first letter of each word to uppercase.
func ToTitle(s string) string {
	return strings.Title(strings.ToLower(s))
}

// ToSnakeCase converts camelCase or PascalCase to snake_case.
func ToSnakeCase(s string) string {
	return CamelCaseToUnderscore(s)
}

// ToKebabCase converts camelCase, PascalCase, or snake_case to kebab-case.
func ToKebabCase(s string) string {
	if s == "" {
		return s
	}

	// Handle kebab-case input
	if strings.Contains(s, "-") && !strings.Contains(s, "_") {
		return strings.ToLower(s)
	}

	// Handle snake_case input
	if strings.Contains(s, "_") {
		return strings.ReplaceAll(strings.ToLower(s), "_", "-")
	}

	// Handle camelCase or PascalCase input
	snake := CamelCaseToUnderscore(s)
	return strings.ReplaceAll(snake, "_", "-")
}

// ToPascalCase converts snake_case, kebab-case, or camelCase to PascalCase.
func ToPascalCase(s string) string {
	if s == "" {
		return s
	}

	// Handle already PascalCase input
	if !strings.Contains(s, "_") && !strings.Contains(s, "-") {
		// If it's already PascalCase, return as is
		if len(s) > 0 && IsUpperChar(rune(s[0])) {
			return s
		}
		// If it's camelCase, convert to PascalCase
		if len(s) > 0 {
			return strings.ToUpper(string(s[0])) + s[1:]
		}
		return s
	}

	// Handle kebab-case input by converting to snake_case first
	if strings.Contains(s, "-") && !strings.Contains(s, "_") {
		s = strings.ReplaceAll(s, "-", "_")
	}

	// Handle snake_case input
	return UnderscoreToCamelCase(s)
}

// ToCamelCase converts snake_case, kebab-case, or PascalCase to camelCase.
func ToCamelCase(s string) string {
	if s == "" {
		return s
	}

	// Handle already camelCase or PascalCase input
	if !strings.Contains(s, "_") && !strings.Contains(s, "-") {
		// If it's already PascalCase, convert to camelCase
		if len(s) > 0 && IsUpperChar(rune(s[0])) {
			return strings.ToLower(string(s[0])) + s[1:]
		}
		// If it's already camelCase, return as is
		return s
	}

	// Handle kebab-case input by converting to snake_case first
	if strings.Contains(s, "-") && !strings.Contains(s, "_") {
		s = strings.ReplaceAll(s, "-", "_")
	}

	// Handle snake_case input
	pascal := UnderscoreToCamelCase(s)
	if pascal == "" || len(pascal) == 0 {
		return pascal
	}
	return strings.ToLower(string(pascal[0])) + pascal[1:]
}

// SplitCamelCase splits and returns words in camelCase format.
//
// Example:
//
//	SplitCamelCase("loremIpsum") //Output []string{"lorem", "Ipsum"}
func SplitCamelCase(str string) []string {
	str = strings.TrimSpace(str)
	if Len(str) < 2 {
		return []string{str}
	}
	var prev rune
	var start int
	words := []string{}
	runes := []rune(str)

	for i, r := range runes {
		if i != 0 {
			switch {
			case unicode.IsDigit(r) && unicode.IsLetter(prev):
				fallthrough
			case unicode.IsUpper(r) && unicode.IsLower(prev):
				words = append(words, string(runes[start:i]))
				start = i
			case unicode.IsLower(r) && unicode.IsUpper(prev) && start != i-1:
				words = append(words, string(runes[start:i-1]))
				start = i - 1
			}
		}
		prev = r
	}

	if start < len(runes) {
		words = append(words, string(runes[start:]))
	}
	return words
}
