package util

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	identifierRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

	reservedKeywords = map[string]bool{
		"break": true, "case": true, "chan": true, "const": true,
		"continue": true, "default": true, "defer": true, "else": true,
		"fallthrough": true, "for": true, "func": true, "go": true,
		"goto": true, "if": true, "import": true, "interface": true,
		"map": true, "package": true, "range": true, "return": true,
		"select": true, "struct": true, "switch": true, "type": true,
		"var": true,
	}
)

func ValidateIdentifier(name string) error {
	if name == "" {
		return fmt.Errorf("identifier cannot be empty")
	}

	if !identifierRegex.MatchString(name) {
		return fmt.Errorf("invalid identifier: %s (must start with letter or underscore)", name)
	}

	if reservedKeywords[name] {
		return fmt.Errorf("identifier cannot be a reserved keyword: %s", name)
	}

	return nil
}

func ValidatePascalCase(name string) error {
	if err := ValidateIdentifier(name); err != nil {
		return err
	}

	if len(name) > 0 && !unicode.IsUpper(rune(name[0])) {
		return fmt.Errorf("PascalCase identifier must start with uppercase letter: %s", name)
	}

	return nil
}

func ValidateType(typeName string) error {
	if typeName == "" {
		return fmt.Errorf("type cannot be empty")
	}

	basicTypes := map[string]bool{
		"string": true, "bool": true,
		"int": true, "int8": true, "int16": true, "int32": true, "int64": true,
		"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
		"float32": true, "float64": true,
		"byte": true, "rune": true,
		"complex64": true, "complex128": true,
	}

	baseType := typeName
	baseType = strings.TrimPrefix(baseType, "*")
	baseType = strings.TrimPrefix(baseType, "[]")

	if basicTypes[baseType] {
		return nil
	}

	if strings.HasPrefix(baseType, "map[") || strings.HasPrefix(baseType, "chan ") {
		return nil
	}

	if strings.Contains(baseType, ".") {
		parts := strings.Split(baseType, ".")
		if len(parts) != 2 {
			return fmt.Errorf("invalid qualified type: %s", typeName)
		}
		return ValidateIdentifier(parts[1])
	}

	return ValidateIdentifier(baseType)
}
