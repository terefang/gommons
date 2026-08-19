package xstrings

import (
	"reflect"
	"strconv"
	"strings"
)

// Cut truncates a string to the specified maximum length and appends "..." if truncated.
// If the string length is less than or equal to max, it returns the original string.
func Cut(max int, s string) string {
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}

// RemSpace removes all spaces from the string.
// Returns the original string if it's empty.
func RemSpace(s string) string {
	if s == "" {
		return s
	}
	return strings.ReplaceAll(s, " ", "")
}

// IsSpace checks if the string contains only spaces (after removing all spaces).
// Returns true if the string is empty or contains only spaces.
func IsSpace(s string) bool {
	return RemSpace(s) == ""
}

// HasLen checks if the string has meaningful content (non-space characters).
// Returns true if the string contains non-space characters.
func HasLen(s string) bool {
	return RemSpace(s) != ""
}

// IsEmpty checks if the string is empty or contains only whitespace.
func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

// BlurEmail blurs the local part of an email address while keeping the domain visible.
// Replaces characters in the local part with asterisks for privacy.
func BlurEmail(email string) string {
	p := strings.IndexByte(email, '@')
	if p < 0 {
		return email
	}
	name := email[0:p]
	domain := email[p:]
	b := strings.Builder{}
	blur := Blur(name, 1, len(name)-1, "*", 4)
	b.WriteString(blur)
	b.WriteString(domain)
	return b.String()
}

// Blur replaces characters in a string with a separator for privacy.
// Parameters:
//   - str: the input string
//   - start: start position for blurring
//   - end: end position for blurring
//   - sep: separator to use for blurring
//   - num: number of separator characters to use
func Blur(str string, start int, end int, sep string, num int) string {
	l := len(str)
	if start >= l {
		return str
	}
	b := strings.Builder{}
	prev := str[0:start]
	suf := str[end:l]
	b.WriteString(prev)
	for i := 0; i < num; i++ {
		b.WriteString(sep)
	}
	b.WriteString(suf)
	return b.String()
}

// EndsWith checks if the string ends with the specified suffix.
// Returns true if the suffix is empty or if the string ends with the suffix.
func EndsWith(str string, sep string) bool {
	l := len(str)
	if sep == "" {
		return true
	}
	if l == 0 {
		return false
	}
	i := strings.LastIndex(str, sep)
	return l-i-len(sep) == 0
}

// StartsWith checks if the string starts with the specified prefix.
// Returns true if the prefix is empty or if the string starts with the prefix.
func StartsWith(str string, sep string) bool {
	l := len(str)
	if sep == "" {
		return true
	}
	if l == 0 {
		return false
	}
	i := strings.Index(str, sep)
	return i == 0
}

// String converts any value to its string representation.
// Handles various types including integers, floats, strings, booleans, and pointers.
// For pointer types, it returns the type name and memory address.
func String(a any) string {
	v := reflect.ValueOf(a)
	switch v.Kind() {
	case reflect.Invalid: // Zero value
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10) // Format as decimal integer
	case reflect.Uint8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', 10, 64)
	case reflect.Complex64, reflect.Complex128:
		return strconv.FormatComplex(v.Complex(), 'f', 10, 64)
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Chan, reflect.Map, reflect.Func, reflect.Pointer, reflect.Slice: // Reference types and functions, pointers
		return v.Type().String() + " 0x" + strconv.FormatUint(uint64(v.Pointer()), 16) // Memory address
	default:
		return v.Type().String() + " value" // Default output: type name + " value"
	}
}

// IsUpperChar checks if a rune is an uppercase letter (A-Z).
func IsUpperChar(c rune) bool {
	return c >= 'A' && c <= 'Z'
}

// CamelCaseToUnderscore converts camelCase or PascalCase strings to underscore_case.
// Handles various patterns like "HelloWorld" -> "hello_world", "helloWorld" -> "hello_world".
func CamelCaseToUnderscore(s string) string {
	if s == "" {
		return s
	}
	rs := []rune(s)
	sb := strings.Builder{}
	sb.WriteRune(rs[0])
	added := false // Flag to track consecutive underscore additions
	for i, r := range rs {
		if i == 0 {
			continue
		}
		// Uppercase followed by lowercase, add underscore first
		if !added && IsUpperChar(r) && (i < len(rs)-1 && !IsUpperChar(rs[i+1])) {
			sb.WriteString("_")
			sb.WriteRune(r)
			added = true
			// Uppercase followed by uppercase, add underscore after
		} else if !added && !IsUpperChar(r) && (i < len(rs)-1 && IsUpperChar(rs[i+1])) {
			sb.WriteRune(r)
			sb.WriteString("_")
			added = true
		} else {
			sb.WriteRune(r)
			added = false
		}
	}
	return strings.ToLower(sb.String())
}

// UnderscoreToCamelCase converts underscore_case strings to camelCase.
// Converts "hello_world" to "HelloWorld" (PascalCase).
func UnderscoreToCamelCase(s string) string {
	rs := []rune(s)
	sb := strings.Builder{}
	lastUnderscore := false
	for i, r := range rs {
		if i == 0 {
			sb.WriteString(strings.ToUpper(string(r)))
			continue
		}
		if string(r) != "_" {
			if !IsUpperChar(r) && lastUnderscore {
				sb.WriteString(strings.ToUpper(string(r)))
			} else {
				sb.WriteRune(r)
			}
			lastUnderscore = false
		} else {
			lastUnderscore = true
		}
	}
	return sb.String()
}

// IsNumeric checks if the string contains only numeric characters.
func IsNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// IsAlpha checks if the string contains only alphabetic characters.
func IsAlpha(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}

// IsAlphaNumeric checks if the string contains only alphanumeric characters.
func IsAlphaNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

// IsEmail checks if the string is a valid email format (basic validation).
func IsEmail(s string) bool {
	if s == "" {
		return false
	}
	parts := strings.Split(s, "@")
	if len(parts) != 2 {
		return false
	}
	if parts[0] == "" || parts[1] == "" {
		return false
	}
	return strings.Contains(parts[1], ".")
}

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

// Reverse reverses the string.
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// RemoveDuplicates removes consecutive duplicate characters from the string.
func RemoveDuplicates(s string) string {
	if s == "" {
		return s
	}
	var result strings.Builder
	runes := []rune(s)
	result.WriteRune(runes[0])
	for i := 1; i < len(runes); i++ {
		if runes[i] != runes[i-1] {
			result.WriteRune(runes[i])
		}
	}
	return result.String()
}

// PadLeft pads the string to the left with the specified character to reach the target length.
func PadLeft(s string, length int, padChar rune) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(string(padChar), length-len(s))
	return padding + s
}

// PadRight pads the string to the right with the specified character to reach the target length.
func PadRight(s string, length int, padChar rune) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(string(padChar), length-len(s))
	return s + padding
}

// PadCenter pads the string on both sides with the specified character to reach the target length.
func PadCenter(s string, length int, padChar rune) string {
	if len(s) >= length {
		return s
	}
	totalPadding := length - len(s)
	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding
	return strings.Repeat(string(padChar), leftPadding) + s + strings.Repeat(string(padChar), rightPadding)
}

// SplitAndTrim splits the string by separator and trims whitespace from each part.
func SplitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// JoinNonEmpty joins non-empty strings with the specified separator.
func JoinNonEmpty(sep string, strs ...string) string {
	var nonEmpty []string
	for _, s := range strs {
		if strings.TrimSpace(s) != "" {
			nonEmpty = append(nonEmpty, s)
		}
	}
	return strings.Join(nonEmpty, sep)
}

// Chunk splits the string into chunks of specified size.
func Chunk(s string, size int) []string {
	if size <= 0 || s == "" {
		return []string{s}
	}
	runes := []rune(s)
	var chunks []string
	for i := 0; i < len(runes); i += size {
		end := i + size
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
	}
	return chunks
}

// Wrap wraps the string to the specified width, breaking at word boundaries.
func Wrap(s string, width int) []string {
	if width <= 0 {
		return []string{s}
	}
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine strings.Builder
	currentLine.WriteString(words[0])

	for i := 1; i < len(words); i++ {
		word := words[i]
		if currentLine.Len()+1+len(word) <= width {
			currentLine.WriteString(" " + word)
		} else {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// CountWords counts the number of words in the string.
func CountWords(s string) int {
	words := strings.Fields(s)
	return len(words)
}

// CountLines counts the number of lines in the string.
func CountLines(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n") + 1
}

// CountOccurrences counts the number of occurrences of a substring in the string.
func CountOccurrences(s, substr string) int {
	if substr == "" {
		return 0
	}
	return strings.Count(s, substr)
}

// FindAll finds all occurrences of a substring and returns their positions.
func FindAll(s, substr string) []int {
	if substr == "" {
		return nil
	}
	var positions []int
	start := 0
	for {
		pos := strings.Index(s[start:], substr)
		if pos == -1 {
			break
		}
		positions = append(positions, start+pos)
		start += pos + len(substr)
	}
	return positions
}

// ContainsAny checks if the string contains any of the specified substrings.
func ContainsAny(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// ContainsAll checks if the string contains all of the specified substrings.
func ContainsAll(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if !strings.Contains(s, substr) {
			return false
		}
	}
	return true
}
