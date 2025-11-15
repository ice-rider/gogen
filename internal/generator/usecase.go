package generator

import (
	"context"
	"fmt"
	"path/filepath"

	"gogen/internal/template"
	"gogen/internal/util"
	"gogen/pkg/models"
)

func (g *Generator) GenerateUseCase(ctx context.Context, uc *models.UseCaseConfig, plan *models.GenerationPlan) error {

	for i := range uc.Dependencies {
		dep := &uc.Dependencies[i]

		if dep.Type == "repository" || dep.Type == "" {
			repoName := util.TrimSuffix(dep.Name, "Repository")
			dep.Found = plan.HasRepository(repoName)
		}
	}

	missing := uc.GetMissingDependencies()
	if len(missing) > 0 {
		return fmt.Errorf("use case %s has missing dependencies: %v", uc.Name, missing)
	}

	data := template.UseCaseData{
		Name:         uc.Name,
		Description:  uc.Description,
		ModulePath:   plan.ModulePath,
		Dependencies: uc.Dependencies,
		InputFields:  uc.InputFields,
		OutputFields: uc.OutputFields,
		WithLogging:  uc.WithLogging,
		WithMetrics:  uc.WithMetrics,
		AddComments:  uc.AddComments || g.config.Generation.AddComments,
		Example:      uc.Example,
	}

	if data.Description == "" {
		data.Description = fmt.Sprintf("операцию %s", uc.Name)
	}

	content, err := g.renderer.Render("usecase", data)
	if err != nil {
		return err
	}

	formatted, err := g.formatter.Format(content)
	if err != nil {
		return fmt.Errorf("generated code has syntax errors: %w", err)
	}

	withImports, err := g.imports.OrganizeImports(formatted)
	if err != nil {
		withImports = formatted
	}

	fileName := util.ToSnakeCase(uc.Name) + "_usecase.go"
	filePath := filepath.Join(g.config.Paths.UseCase, fileName)

	if err := g.writer.Write(filePath, withImports, false); err != nil {
		return err
	}

	return nil
}
