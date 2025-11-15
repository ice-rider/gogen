package parser

import (
	"fmt"
	"gogen/internal/util"
	"strings"
)

type TagParser struct{}

func NewTagParser() *TagParser {
	return &TagParser{}
}

type Tag struct {
	Key   string
	Value string
}

func (tp *TagParser) Parse(input string) ([]Tag, error) {
	input = strings.Trim(input, "")

	var tags []Tag
	parts := strings.Fields(input)

	for _, part := range parts {
		tag, err := tp.parseTag(part)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (tp *TagParser) parseTag(input string) (Tag, error) {
	parts := strings.SplitN(input, ":", 2)
	if len(parts) != 2 {
		return Tag{}, fmt.Errorf("invalid tag format: %s", input)
	}

	key := parts[0]
	value := strings.Trim(parts[1], "`")

	return Tag{
		Key:   key,
		Value: value,
	}, nil
}

func (tp *TagParser) BuildStructTag(tags []Tag) string {
	if len(tags) == 0 {
		return ""
	}

	var parts []string
	for _, tag := range tags {
		parts = append(parts, fmt.Sprintf(`%s:"%s"`, tag.Key, tag.Value))
	}

	return "" + strings.Join(parts, " ") + "`"
}

func (tp *TagParser) BuildFieldTags(fieldName, jsonStyle string, required, unique, index bool) string {
	var tags []Tag

	jsonTag := util.ToSnakeCase(fieldName)
	if jsonStyle == "camelCase" {
		jsonTag = util.ToCamelCase(fieldName)
	} else if jsonStyle == "PascalCase" {
		jsonTag = util.ToPascalCase(fieldName)
	}
	tags = append(tags, Tag{Key: "json", Value: jsonTag})

	dbTag := util.ToSnakeCase(fieldName)
	tags = append(tags, Tag{Key: "db", Value: dbTag})

	if required || unique {
		var validates []string
		if required {
			validates = append(validates, "required")
		}
		if unique {
			validates = append(validates, "unique")
		}
		tags = append(tags, Tag{
			Key:   "validate",
			Value: strings.Join(validates, ","),
		})
	}

	return tp.BuildStructTag(tags)
}
