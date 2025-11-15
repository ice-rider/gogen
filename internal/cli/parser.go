package cli

import (
	"fmt"
	"strings"

	"gogen/internal/parser"
	"gogen/internal/util"
	"gogen/pkg/models"
)

type Parser struct {
	fieldParser *parser.FieldParser
}

func NewParser() *Parser {
	return &Parser{
		fieldParser: parser.NewFieldParser(),
	}
}

func (p *Parser) BuildPlan(flags *Flags) (*models.GenerationPlan, error) {
	plan := &models.GenerationPlan{
		WithTests: flags.WithTests,
		WithMocks: flags.WithMocks,
	}

	for _, entityName := range flags.Entities {
		entity, err := p.parseEntity(entityName)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга сущности %s: %w", entityName, err)
		}
		plan.Entities = append(plan.Entities, entity)
	}

	for _, repoName := range flags.Repositories {
		repo, err := p.parseRepository(repoName)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга репозитория %s: %w", repoName, err)
		}
		plan.Repositories = append(plan.Repositories, repo)
	}

	for _, ucName := range flags.UseCases {
		uc, err := p.parseUseCase(ucName)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга use case %s: %w", ucName, err)
		}
		plan.UseCases = append(plan.UseCases, uc)
	}

	return plan, nil
}

func (p *Parser) parseEntity(input string) (models.EntityConfig, error) {

	parts := strings.SplitN(input, ":", 2)

	name := strings.TrimSpace(parts[0])

	if err := util.ValidatePascalCase(name); err != nil {
		return models.EntityConfig{}, err
	}

	entity := models.EntityConfig{
		Name:          name,
		TableName:     util.ToSnakeCase(util.Pluralize(name)),
		AddValidation: true,
		AddComments:   true,
		JSONStyle:     "snake_case",
	}

	if len(parts) > 1 {
		fields, err := p.fieldParser.Parse(parts[1])
		if err != nil {
			return models.EntityConfig{}, err
		}
		entity.Fields = fields
	}

	return entity, nil
}

func (p *Parser) parseRepository(input string) (models.RepositoryConfig, error) {
	name := strings.TrimSpace(input)

	name = strings.TrimSuffix(name, "Repository")

	if err := util.ValidatePascalCase(name); err != nil {
		return models.RepositoryConfig{}, err
	}

	repo := models.RepositoryConfig{
		Name:             name,
		Entity:           name,
		TableName:        util.ToSnakeCase(util.Pluralize(name)),
		DBType:           "postgres",
		WithTransactions: true,
		AddComments:      true,
	}

	return repo, nil
}

func (p *Parser) parseUseCase(input string) (models.UseCaseConfig, error) {
	name := strings.TrimSpace(input)

	name = strings.TrimSuffix(name, "UseCase")

	if err := util.ValidatePascalCase(name); err != nil {
		return models.UseCaseConfig{}, err
	}

	uc := models.UseCaseConfig{
		Name:         name,
		Description:  fmt.Sprintf("операцию %s", name),
		Dependencies: []models.Dependency{},
		InputFields:  []models.Field{},
		OutputFields: []models.Field{},
		WithLogging:  false,
		WithMetrics:  false,
		AddComments:  true,
	}

	return uc, nil
}
