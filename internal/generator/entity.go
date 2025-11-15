package generator

import (
	"context"
	"fmt"
	"path/filepath"

	"gogen/internal/template"
	"gogen/internal/util"
	"gogen/pkg/models"
)

func (g *Generator) GenerateEntity(ctx context.Context, entity *models.EntityConfig, plan *models.GenerationPlan) error {

	data := template.EntityData{
		Name:          entity.Name,
		Fields:        entity.Fields,
		TableName:     entity.TableName,
		ModulePath:    plan.ModulePath,
		AddComments:   entity.AddComments || g.config.Generation.AddComments,
		AddValidation: entity.AddValidation,
		JSONStyle:     entity.JSONStyle,
	}

	if data.TableName == "" {
		data.TableName = util.ToSnakeCase(util.Pluralize(entity.Name))
	}

	content, err := g.renderer.Render("entity", data)
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

	fileName := util.ToSnakeCase(entity.Name) + ".go"
	filePath := filepath.Join(g.config.Paths.Domain, fileName)

	if err := g.writer.Write(filePath, withImports, false); err != nil {
		return err
	}

	return nil
}
