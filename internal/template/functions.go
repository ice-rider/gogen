package template

import (
	"strings"
	"text/template"

	"gogen/internal/util"
)

func (l *Loader) getFuncMap() template.FuncMap {
	return template.FuncMap{

		"ToSnakeCase":  util.ToSnakeCase,
		"ToPascalCase": util.ToPascalCase,
		"ToCamelCase":  util.ToCamelCase,
		"ToLower":      strings.ToLower,
		"ToUpper":      strings.ToUpper,
		"ToTitle":      strings.Title,

		"Pluralize":   util.Pluralize,
		"Singularize": util.Singularize,

		"TrimSuffix": strings.TrimSuffix,
		"TrimPrefix": strings.TrimPrefix,
		"Contains":   strings.Contains,
		"Replace":    strings.ReplaceAll,
		"Split":      strings.Split,
		"Join":       strings.Join,

		"Add": func(a, b int) int { return a + b },
		"Sub": func(a, b int) int { return a - b },
		"Mul": func(a, b int) int { return a * b },
		"Div": func(a, b int) int { return a / b },

		"IsEmpty": func(s string) bool { return s == "" },
		"IsZero":  func(v interface{}) bool { return v == nil },

		"GetZeroValue": getZeroValue,
		"IsPointer":    isPointerType,
		"GetBaseType":  getBaseType,

		"NeedTimeImport": needTimeImport,
		"NeedUUIDImport": needUUIDImport,
	}
}

func getZeroValue(typeName string) string {
	switch typeName {
	case "string":
		return ""
	case "int", "int8", "int16", "int32", "int64":
		return "0"
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return "0"
	case "float32", "float64":
		return "0.0"
	case "bool":
		return "false"
	case "time.Time":
		return "time.Time{}"
	case "uuid.UUID":
		return "uuid.Nil"
	default:
		if strings.HasPrefix(typeName, "*") {
			return "nil"
		}
		if strings.HasPrefix(typeName, "[]") {
			return "nil"
		}
		if strings.HasPrefix(typeName, "map[") {
			return "nil"
		}
		return "nil"
	}
}

func isPointerType(typeName string) bool {
	return strings.HasPrefix(typeName, "*")
}

func getBaseType(typeName string) string {
	t := strings.TrimPrefix(typeName, "*")
	t = strings.TrimPrefix(t, "[]")
	return t
}

func needTimeImport(fields interface{}) bool {

	return true
}

func needUUIDImport(fields interface{}) bool {
	return true
}
