package transfer

import (
	"strings"
	"unicode"
)

func FirstLower(s string) string {
	if s == "" {
		return s
	}
	first := s[:1]
	return strings.ToLower(first)
}

func ToCamelCase(s string) string {
	if s == "" {
		return s
	}
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(part)
		} else {
			if len(part) > 0 {
				runes := []rune(part)
				runes[0] = unicode.ToUpper(runes[0])
				parts[i] = string(runes)
			}
		}
	}
	return strings.Join(parts, "")
}

func ToPascalCase(s string) string {
	if s == "" {
		return s
	}
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			runes := []rune(part)
			runes[0] = unicode.ToUpper(runes[0])
			parts[i] = string(runes)
		}
	}
	return strings.Join(parts, "")
}
