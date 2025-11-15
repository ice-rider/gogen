package interactive

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"

	"gogen/internal/logger"
	"gogen/pkg/models"
)

type Interactor struct {
	logger *logger.Logger
}

func NewInteractor(logger *logger.Logger) *Interactor {
	return &Interactor{
		logger: logger,
	}
}

func (i *Interactor) EnhancePlan(plan *models.GenerationPlan, cfg *models.Config) error {
	i.logger.Section("🎯 Интерактивная настройка компонентов")

	for idx := range plan.Entities {
		entity := &plan.Entities[idx]

		if err := i.configureEntity(entity, cfg); err != nil {
			return err
		}
	}

	for idx := range plan.Repositories {
		repo := &plan.Repositories[idx]

		if err := i.configureRepository(repo, plan, cfg); err != nil {
			return err
		}
	}

	for idx := range plan.UseCases {
		uc := &plan.UseCases[idx]

		if err := i.configureUseCase(uc, plan, cfg); err != nil {
			return err
		}
	}

	return i.confirmGeneration(plan)
}

func (i *Interactor) configureEntity(entity *models.EntityConfig, cfg *models.Config) error {
	i.logger.Info("Настройка сущности: %s", entity.Name)

	if len(entity.Fields) == 0 {
		addFields := false
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Добавить поля для %s?", entity.Name),
			Default: true,
		}
		if err := survey.AskOne(prompt, &addFields); err != nil {
			return err
		}

		if addFields {
			fieldsPrompter := NewFieldsPrompter()
			fields, err := fieldsPrompter.PromptFields()
			if err != nil {
				return err
			}
			entity.Fields = fields
		}
	}

	tableName := entity.TableName
	survey.AskOne(&survey.Input{
		Message: "Название таблицы в БД:",
		Default: tableName,
	}, &tableName)
	entity.TableName = tableName

	survey.AskOne(&survey.Confirm{
		Message: "Добавить методы валидации?",
		Default: true,
	}, &entity.AddValidation)

	i.logger.Success("Сущность %s настроена", entity.Name)
	fmt.Println()

	return nil
}

func (i *Interactor) configureRepository(repo *models.RepositoryConfig, plan *models.GenerationPlan, cfg *models.Config) error {
	i.logger.Info("Настройка репозитория: %sRepository", repo.Name)

	dbType := repo.DBType
	survey.AskOne(&survey.Select{
		Message: "Тип базы данных:",
		Options: []string{"postgres", "mysql", "sqlite", "mongodb"},
		Default: dbType,
	}, &dbType)
	repo.DBType = dbType

	addCustom := false
	survey.AskOne(&survey.Confirm{
		Message: "Добавить кастомные методы репозитория?",
		Default: false,
	}, &addCustom)

	if addCustom {
		methodsPrompter := NewMethodsPrompter()
		methods, err := methodsPrompter.PromptMethods()
		if err != nil {
			return err
		}
		repo.CustomMethods = methods
	}

	survey.AskOne(&survey.Confirm{
		Message: "Поддержка транзакций?",
		Default: true,
	}, &repo.WithTransactions)

	i.logger.Success("Репозиторий %s настроен", repo.Name)
	fmt.Println()

	return nil
}

func (i *Interactor) configureUseCase(uc *models.UseCaseConfig, plan *models.GenerationPlan, cfg *models.Config) error {
	i.logger.Info("Настройка use case: %sUseCase", uc.Name)

	description := uc.Description
	survey.AskOne(&survey.Input{
		Message: "Описание use case:",
		Default: description,
	}, &description)
	uc.Description = description

	addInput := false
	survey.AskOne(&survey.Confirm{
		Message: "Определить входные параметры (Input)?",
		Default: len(uc.InputFields) == 0,
	}, &addInput)

	if addInput {
		fieldsPrompter := NewFieldsPrompter()
		fields, err := fieldsPrompter.PromptFields()
		if err != nil {
			return err
		}
		uc.InputFields = fields
	}

	addOutput := false
	survey.AskOne(&survey.Confirm{
		Message: "Определить выходные параметры (Output)?",
		Default: len(uc.OutputFields) == 0,
	}, &addOutput)

	if addOutput {
		fieldsPrompter := NewFieldsPrompter()
		fields, err := fieldsPrompter.PromptFields()
		if err != nil {
			return err
		}
		uc.OutputFields = fields
	}

	survey.AskOne(&survey.Confirm{
		Message: "Добавить логирование?",
		Default: false,
	}, &uc.WithLogging)

	survey.AskOne(&survey.Confirm{
		Message: "Добавить метрики?",
		Default: false,
	}, &uc.WithMetrics)

	i.logger.Success("Use case %s настроен", uc.Name)
	fmt.Println()

	return nil
}

func (i *Interactor) confirmGeneration(plan *models.GenerationPlan) error {
	i.logger.Section("📋 Итоговый план генерации")

	fmt.Printf("Будет создано:\n")
	fmt.Printf("  • Сущностей: %d\n", len(plan.Entities))
	fmt.Printf("  • Репозиториев: %d\n", len(plan.Repositories))
	fmt.Printf("  • Use Cases: %d\n", len(plan.UseCases))

	if plan.WithTests {
		fmt.Println("  • Тесты для всех компонентов")
	}
	if plan.WithMocks {
		fmt.Println("  • Моки для всех репозиториев")
	}

	fmt.Println()

	confirm := false
	prompt := &survey.Confirm{
		Message: "Начать генерацию?",
		Default: true,
	}

	if err := survey.AskOne(prompt, &confirm); err != nil {
		return err
	}

	if !confirm {
		return fmt.Errorf("генерация отменена пользователем")
	}

	return nil
}
