package generator

import (
	"context"
	"fmt"
	"path/filepath"

	"gogen/internal/template"
	"gogen/internal/util"
	"gogen/pkg/models"
)

func (g *Generator) GenerateTests(ctx context.Context, plan *models.GenerationPlan) error {

	for _, entity := range plan.Entities {
		if err := g.generateEntityTest(ctx, &entity, plan); err != nil {
			return fmt.Errorf("failed to generate test for entity %s: %w", entity.Name, err)
		}
	}

	for _, repo := range plan.Repositories {
		if err := g.generateRepositoryTest(ctx, &repo, plan); err != nil {
			return fmt.Errorf("failed to generate test for repository %s: %w", repo.Name, err)
		}
	}

	for _, uc := range plan.UseCases {
		if err := g.generateUseCaseTest(ctx, &uc, plan); err != nil {
			return fmt.Errorf("failed to generate test for usecase %s: %w", uc.Name, err)
		}
	}

	return nil
}

func (g *Generator) generateEntityTest(ctx context.Context, entity *models.EntityConfig, plan *models.GenerationPlan) error {
	data := struct {
		Name       string
		Fields     []models.Field
		ModulePath string
	}{
		Name:       entity.Name,
		Fields:     entity.Fields,
		ModulePath: plan.ModulePath,
	}

	content, err := g.renderer.Render("test_entity", data)
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

	fileName := util.ToSnakeCase(entity.Name) + "_test.go"
	filePath := filepath.Join(g.config.Paths.Domain, fileName)

	if err := g.writer.Write(filePath, withImports, false); err != nil {
		return err
	}

	return nil
}

func (g *Generator) generateRepositoryTest(ctx context.Context, repo *models.RepositoryConfig, plan *models.GenerationPlan) error {
	data := struct {
		Name       string
		Entity     string
		ModulePath string
		Methods    []template.MockMethod
	}{
		Name:       repo.Name,
		Entity:     repo.Entity,
		ModulePath: plan.ModulePath,
		Methods:    g.collectRepositoryMethods(repo),
	}

	content, err := g.renderer.Render("test_repository", data)
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

	fileName := util.ToSnakeCase(repo.Name) + "_repository_test.go"
	filePath := filepath.Join(g.config.Paths.Repository, fileName)

	if err := g.writer.Write(filePath, withImports, false); err != nil {
		return err
	}

	return nil
}

func (g *Generator) generateUseCaseTest(ctx context.Context, uc *models.UseCaseConfig, plan *models.GenerationPlan) error {
	data := struct {
		Name         string
		ModulePath   string
		Dependencies []models.Dependency
		InputFields  []models.Field
		OutputFields []models.Field
	}{
		Name:         uc.Name,
		ModulePath:   plan.ModulePath,
		Dependencies: uc.Dependencies,
		InputFields:  uc.InputFields,
		OutputFields: uc.OutputFields,
	}

	content, err := g.renderer.Render("test_usecase", data)
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

	fileName := util.ToSnakeCase(uc.Name) + "_usecase_test.go"
	filePath := filepath.Join(g.config.Paths.UseCase, fileName)

	if err := g.writer.Write(filePath, withImports, false); err != nil {
		return err
	}

	return nil
}
