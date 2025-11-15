package interactive

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"gogen/pkg/models"
)

type MethodsPrompter struct{}

func NewMethodsPrompter() *MethodsPrompter {
	return &MethodsPrompter{}
}

func (mp *MethodsPrompter) PromptMethods() ([]models.CustomMethod, error) {
	fmt.Println("\nДобавление кастомных методов репозитория")
	fmt.Println("Формат: MethodName(param1 Type, param2 Type) (ReturnType, error)")
	fmt.Println("Пример: FindByEmail(email string) (*User, error)")
	fmt.Println("Пустая строка для завершения\n")

	var methods []models.CustomMethod
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

		method, err := mp.parseMethod(input)
		if err != nil {
			fmt.Printf("  ⚠️  Ошибка: %v, попробуйте снова\n", err)
			continue
		}

		methods = append(methods, method)
		i++
	}

	return methods, nil
}

func (mp *MethodsPrompter) parseMethod(input string) (models.CustomMethod, error) {

	parenIdx := strings.Index(input, "(")
	if parenIdx == -1 {
		return models.CustomMethod{}, fmt.Errorf("неверный формат метода")
	}

	name := strings.TrimSpace(input[:parenIdx])

	method := models.CustomMethod{
		Name:    name,
		Comment: fmt.Sprintf("%s кастомный метод репозитория", name),
		Params:  []models.MethodParam{},
		Returns: []string{"error"},
	}

	return method, nil
}
