package parser

import (
	"fmt"
	"strings"
)

type TypeParser struct{}

func NewTypeParser() *TypeParser {
	return &TypeParser{}
}

type TypeInfo struct {
	BaseType    string
	IsPointer   bool
	IsSlice     bool
	IsMap       bool
	KeyType     string
	ValueType   string
	Package     string
	NeedsImport bool
}

func (tp *TypeParser) Parse(typeName string) (*TypeInfo, error) {
	info := &TypeInfo{}

	typeName = strings.TrimSpace(typeName)

	if strings.HasPrefix(typeName, "*") {
		info.IsPointer = true
		typeName = strings.TrimPrefix(typeName, "*")
	}

	if strings.HasPrefix(typeName, "[]") {
		info.IsSlice = true
		typeName = strings.TrimPrefix(typeName, "[]")
		info.ValueType = typeName
	}

	if strings.HasPrefix(typeName, "map[") {
		info.IsMap = true

		typeName = strings.TrimPrefix(typeName, "map[")
		typeName = strings.TrimSuffix(typeName, "]")

		parts := strings.SplitN(typeName, "]", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid map type format")
		}

		info.KeyType = strings.TrimSpace(parts[0])
		info.ValueType = strings.TrimSpace(parts[1])
	}

	if strings.Contains(typeName, ".") {
		parts := strings.Split(typeName, ".")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid qualified type: %s", typeName)
		}

		info.Package = parts[0]
		info.BaseType = parts[1]
		info.NeedsImport = true
	} else {
		info.BaseType = typeName
		info.NeedsImport = tp.needsImport(typeName)
	}

	return info, nil
}

func (tp *TypeParser) needsImport(typeName string) bool {

	basicTypes := map[string]bool{
		"string": true, "bool": true,
		"int": true, "int8": true, "int16": true, "int32": true, "int64": true,
		"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
		"float32": true, "float64": true,
		"byte": true, "rune": true,
		"complex64": true, "complex128": true,
		"error": true,
	}

	return !basicTypes[typeName]
}

func (tp *TypeParser) GetImportPath(typeName string) string {

	imports := map[string]string{
		"time":    "time",
		"uuid":    "github.com/google/uuid",
		"sql":     "database/sql",
		"context": "context",
	}

	info, err := tp.Parse(typeName)
	if err != nil || !info.NeedsImport {
		return ""
	}

	if info.Package != "" {
		return imports[info.Package]
	}

	return ""
}
