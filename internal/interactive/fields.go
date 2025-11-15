package interactive

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"gogen/internal/parser"
	"gogen/internal/util"
	"gogen/pkg/models"
)

type FieldsPrompter struct {
	fieldParser *parser.FieldParser
}

func NewFieldsPrompter() *FieldsPrompter {
	return &FieldsPrompter{
		fieldParser: parser.NewFieldParser(),
	}
}

func (fp *FieldsPrompter) PromptFields() ([]models.Field, error) {
	fmt.Println("\nВведите поля (формат: Name:Type или Name:Type:tags)")
	fmt.Println("Доступные теги: required, unique, index")
	fmt.Println("Пример: Email:string:required,unique")
	fmt.Println("Пустая строка для завершения\n")

	var fields []models.Field
	i := 1

	for {
		input := ""
		prompt := &survey.Input{
			Message: fmt.Sprintf("[%d]", i),
		}

		if err := survey.AskOne(prompt, &input); err != nil {
			return nil, err
		}

		if strings.TrimSpace(input) == "" {
			break
		}

		field, err := fp.fieldParser.Parse(input)
		if err != nil {
			fmt.Printf("  ⚠️  Ошибка: %v, попробуйте снова\n", err)
			continue
		}

		fields = append(fields, field...)
		i++
	}

	return fields, nil
}

func (fp *FieldsPrompter) PromptField() (models.Field, error) {
	field := models.Field{}

	survey.AskOne(&survey.Input{
		Message: "Название поля:",
	}, &field.Name, survey.WithValidator(func(ans interface{}) error {
		name := ans.(string)
		return util.ValidatePascalCase(name)
	}))

	typeOptions := []string{
		"string", "int", "int64", "float64", "bool",
		"time.Time", "uuid.UUID",
		"*string (nullable)", "*int (nullable)",
		"[]string (slice)", "[]int (slice)",
		"Другой...",
	}

	var selectedType string
	survey.AskOne(&survey.Select{
		Message: "Тип поля:",
		Options: typeOptions,
	}, &selectedType)

	if selectedType == "Другой..." {
		survey.AskOne(&survey.Input{
			Message: "Введите тип:",
		}, &field.Type, survey.WithValidator(func(ans interface{}) error {
			return util.ValidateType(ans.(string))
		}))
	} else {

		field.Type = strings.Split(selectedType, " ")[0]
	}

	tagOptions := []string{"required", "unique", "index"}
	selectedTags := []string{}

	survey.AskOne(&survey.MultiSelect{
		Message: "Выберите теги:",
		Options: tagOptions,
	}, &selectedTags)

	field.Tags = selectedTags
	for _, tag := range selectedTags {
		switch tag {
		case "required":
			field.Required = true
		case "unique":
			field.Unique = true
		case "index":
			field.Index = true
		}
	}

	survey.AskOne(&survey.Input{
		Message: "Комментарий (опционально):",
	}, &field.Comment)

	field.JSONTag = util.ToSnakeCase(field.Name)
	field.DBTag = util.ToSnakeCase(field.Name)

	return field, nil
}
