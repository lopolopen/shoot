package shoot

import (
	"strings"
	"unicode"
)

const Cmd = "shoot"

type Meta struct {
	Cmd         string
	PackageName string
	TypeName    string
}

func ToCamelCase(s string) string {
	if s == "" {
		return ""
	}

	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '_' || r == '-'
	})

	if len(parts) == 0 {
		return ""
	}

	// Convert the first word to lowercase
	camelCaseStr := strings.ToLower(parts[0])

	// Capitalize the first letter of subsequent words and append them
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			camelCaseStr += string(unicode.ToUpper(rune(parts[i][0]))) + strings.ToLower(parts[i][1:])
		}
	}

	return camelCaseStr
}
