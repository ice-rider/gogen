package generator

import (
	"context"
	"fmt"

	"gogen/internal/file"
	"gogen/internal/format"
	"gogen/internal/template"
	"gogen/pkg/models"
)

type Generator struct {
	renderer  *template.Renderer
	writer    *file.Writer
	formatter *format.Formatter
	imports   *format.ImportsManager
	config    *models.Config
}

func NewGenerator(
	renderer *template.Renderer,
	writer *file.Writer,
	formatter *format.Formatter,
	imports *format.ImportsManager,
	config *models.Config,
) *Generator {
	return &Generator{
		renderer:  renderer,
		writer:    writer,
		formatter: formatter,
		imports:   imports,
		config:    config,
	}
}

func (g *Generator) Generate(ctx context.Context, plan *models.GenerationPlan) error {

	for _, entity := range plan.Entities {
		if err := g.GenerateEntity(ctx, &entity, plan); err != nil {
			return fmt.Errorf("failed to generate entity %s: %w", entity.Name, err)
		}
	}

	for _, repo := range plan.Repositories {
		if err := g.GenerateRepository(ctx, &repo, plan); err != nil {
			return fmt.Errorf("failed to generate repository %s: %w", repo.Name, err)
		}
	}

	for _, uc := range plan.UseCases {
		if err := g.GenerateUseCase(ctx, &uc, plan); err != nil {
			return fmt.Errorf("failed to generate usecase %s: %w", uc.Name, err)
		}
	}

	if plan.WithMocks {
		for _, repo := range plan.Repositories {
			if err := g.GenerateMock(ctx, &repo, plan); err != nil {
				return fmt.Errorf("failed to generate mock for %s: %w", repo.Name, err)
			}
		}
	}

	if plan.WithTests {
		if err := g.GenerateTests(ctx, plan); err != nil {
			return fmt.Errorf("failed to generate tests: %w", err)
		}
	}

	return nil
}
