package models

import "strings"

type Field struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Tags     []string `json:"tags"`
	JSONTag  string   `json:"json_tag"`
	DBTag    string   `json:"db_tag"`
	Comment  string   `json:"comment"`
	Required bool     `json:"required"`
	Unique   bool     `json:"unique"`
	Index    bool     `json:"index"`
}

func (f *Field) ZeroValue() string {
	switch f.Type {
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
		if strings.HasPrefix(f.Type, "*") ||
			strings.HasPrefix(f.Type, "[]") ||
			strings.HasPrefix(f.Type, "map[") {
			return "nil"
		}
		return "nil"
	}
}

func (f *Field) IsPointer() bool {
	return strings.HasPrefix(f.Type, "*")
}

func (f *Field) BaseType() string {
	t := strings.TrimPrefix(f.Type, "*")
	t = strings.TrimPrefix(t, "[]")
	return t
}
