package transfer

import (
	"strings"
	"unicode"
)

func FirtLower(s string) string {
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
	// 分割字符串（假设输入是下划线分隔的格式，如"user_name"）
	parts := strings.Split(s, "_")
	for i, part := range parts {
		// 首字母小写，其他部分首字母大写
		if i == 0 {
			// 第一个单词全小写
			parts[i] = strings.ToLower(part)
		} else {
			// 后续单词首字母大写，其余小写
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
	// 分割字符串
	parts := strings.Split(s, "_")
	for i, part := range parts {
		// 所有单词首字母大写
		if len(part) > 0 {
			runes := []rune(part)
			runes[0] = unicode.ToUpper(runes[0])
			parts[i] = string(runes)
		}
	}
	return strings.Join(parts, "")
}
