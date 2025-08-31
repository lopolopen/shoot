package transfer

import (
	"strings"
	"unicode"
)

func ID(x string) string { return x }

func FirstLower(s string) string {
	if s == "" {
		return s
	}
	first := s[:1]
	return strings.ToLower(first)
}

func ToCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func ToPascalCase(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
