package parser

import (
	"fmt"
	"strings"

	"gogen/internal/util"
	"gogen/pkg/models"
)

type FieldParser struct{}

func NewFieldParser() *FieldParser {
	return &FieldParser{}
}

func (fp *FieldParser) Parse(input string) ([]models.Field, error) {
	if input == "" {
		return []models.Field{}, nil
	}

	fieldStrings := strings.Split(input, ",")
	fields := make([]models.Field, 0, len(fieldStrings))

	for _, fieldStr := range fieldStrings {
		field, err := fp.parseField(strings.TrimSpace(fieldStr))
		if err != nil {
			return nil, fmt.Errorf("failed to parse field '%s': %w", fieldStr, err)
		}
		fields = append(fields, field)
	}

	return fields, nil
}

func (fp *FieldParser) parseField(input string) (models.Field, error) {
	parts := strings.Split(input, ":")

	if len(parts) < 2 {
		return models.Field{}, fmt.Errorf("invalid format, expected Name:Type[:tags]")
	}

	field := models.Field{
		Name: strings.TrimSpace(parts[0]),
		Type: strings.TrimSpace(parts[1]),
	}

	if err := util.ValidatePascalCase(field.Name); err != nil {
		return models.Field{}, fmt.Errorf("invalid field name: %w", err)
	}

	if err := util.ValidateType(field.Type); err != nil {
		return models.Field{}, fmt.Errorf("invalid field type: %w", err)
	}

	if len(parts) >= 3 {
		tags := strings.Split(parts[2], ",")
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
			field.Tags = append(field.Tags, tag)

			switch tag {
			case "required":
				field.Required = true
			case "unique":
				field.Unique = true
			case "index":
				field.Index = true
			}
		}
	}

	field.JSONTag = util.ToSnakeCase(field.Name)
	field.DBTag = util.ToSnakeCase(field.Name)

	return field, nil
}

func (fp *FieldParser) ParseJSON(input string) ([]models.Field, error) {

	return nil, fmt.Errorf("JSON parsing not implemented yet")
}
