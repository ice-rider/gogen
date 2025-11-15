package util

import (
	"strings"
	"unicode"
)

func ToSnakeCase(s string) string {
	var result strings.Builder

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

func ToPascalCase(s string) string {

	if len(s) > 0 && unicode.IsUpper(rune(s[0])) && !strings.Contains(s, "_") {
		return s
	}

	parts := strings.Split(s, "_")
	var result strings.Builder

	for _, part := range parts {
		if len(part) > 0 {
			result.WriteRune(unicode.ToUpper(rune(part[0])))
			result.WriteString(part[1:])
		}
	}

	return result.String()
}

func ToCamelCase(s string) string {
	pascal := ToPascalCase(s)

	if len(pascal) == 0 {
		return pascal
	}

	return strings.ToLower(pascal[:1]) + pascal[1:]
}
