package generator

import (
	"context"
	"fmt"
	"path/filepath"

	"gogen/internal/template"
	"gogen/internal/util"
	"gogen/pkg/models"
)

func (g *Generator) GenerateRepository(ctx context.Context, repo *models.RepositoryConfig, plan *models.GenerationPlan) error {

	entity := plan.GetEntityByName(repo.Entity)
	if entity == nil {
		return fmt.Errorf("entity %s not found for repository %s", repo.Entity, repo.Name)
	}

	if len(repo.Fields) == 0 {
		repo.Fields = entity.Fields
	}

	if repo.TableName == "" {
		repo.TableName = entity.TableName
	}

	if g.config.Generation.SeparateInterfaces {
		if err := g.generateRepositoryInterface(ctx, repo, plan); err != nil {
			return fmt.Errorf("failed to generate repository interface: %w", err)
		}
	}

	if err := g.generateRepositoryImpl(ctx, repo, plan); err != nil {
		return fmt.Errorf("failed to generate repository implementation: %w", err)
	}

	return nil
}

func (g *Generator) generateRepositoryInterface(ctx context.Context, repo *models.RepositoryConfig, plan *models.GenerationPlan) error {
	data := template.RepositoryData{
		Name:          repo.Name,
		Entity:        repo.Entity,
		TableName:     repo.TableName,
		ModulePath:    plan.ModulePath,
		CustomMethods: repo.CustomMethods,
		AddComments:   repo.AddComments || g.config.Generation.AddComments,
		Fields:        repo.Fields,
	}

	content, err := g.renderer.Render("repository_interface", data)
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

	fileName := util.ToSnakeCase(repo.Name) + "_repository.go"
	filePath := filepath.Join(g.config.Paths.Domain, fileName)

	if err := g.writer.Write(filePath, withImports, false); err != nil {
		return err
	}

	return nil
}

func (g *Generator) generateRepositoryImpl(ctx context.Context, repo *models.RepositoryConfig, plan *models.GenerationPlan) error {

	dbType := repo.DBType
	if dbType == "" {
		dbType = "postgres"
	}

	data := template.RepositoryData{
		Name:             repo.Name,
		Entity:           repo.Entity,
		TableName:        repo.TableName,
		ModulePath:       plan.ModulePath,
		DBType:           dbType,
		CustomMethods:    repo.CustomMethods,
		WithTransactions: repo.WithTransactions,
		AddComments:      repo.AddComments || g.config.Generation.AddComments,
		Fields:           repo.Fields,
	}

	content, err := g.renderer.Render("repository_impl", data)
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

	fileName := util.ToSnakeCase(repo.Name) + "_repository.go"
	filePath := filepath.Join(g.config.Paths.Repository, fileName)

	if err := g.writer.Write(filePath, withImports, false); err != nil {
		return err
	}

	return nil
}
