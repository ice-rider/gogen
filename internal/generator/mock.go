package generator

import (
	"context"
	"fmt"
	"path/filepath"

	"gogen/internal/template"
	"gogen/internal/util"
	"gogen/pkg/models"
)

func (g *Generator) GenerateMock(ctx context.Context, repo *models.RepositoryConfig, plan *models.GenerationPlan) error {

	methods := g.collectRepositoryMethods(repo)

	data := template.MockData{
		Name:       repo.Name,
		Entity:     repo.Entity,
		ModulePath: plan.ModulePath,
		Methods:    methods,
	}

	content, err := g.renderer.Render("mock", data)
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

	fileName := util.ToSnakeCase(repo.Name) + "_repository_mock.go"
	filePath := filepath.Join(g.config.Paths.Mocks, fileName)

	if err := g.writer.Write(filePath, withImports, false); err != nil {
		return err
	}

	return nil
}

func (g *Generator) collectRepositoryMethods(repo *models.RepositoryConfig) []template.MockMethod {
	methods := []template.MockMethod{

		{
			Name: "Create",
			Params: []template.MethodParam{
				{Name: "ctx", Type: "context.Context"},
				{Name: "entity", Type: fmt.Sprintf("*domain.%s", repo.Entity)},
			},
			Return: []string{"error"},
		},
		{
			Name: "GetByID",
			Params: []template.MethodParam{
				{Name: "ctx", Type: "context.Context"},
				{Name: "id", Type: "string"},
			},
			Return: []string{fmt.Sprintf("*domain.%s", repo.Entity), "error"},
		},
		{
			Name: "Update",
			Params: []template.MethodParam{
				{Name: "ctx", Type: "context.Context"},
				{Name: "entity", Type: fmt.Sprintf("*domain.%s", repo.Entity)},
			},
			Return: []string{"error"},
		},
		{
			Name: "Delete",
			Params: []template.MethodParam{
				{Name: "ctx", Type: "context.Context"},
				{Name: "id", Type: "string"},
			},
			Return: []string{"error"},
		},
		{
			Name: "List",
			Params: []template.MethodParam{
				{Name: "ctx", Type: "context.Context"},
				{Name: "limit", Type: "int"},
				{Name: "offset", Type: "int"},
			},
			Return: []string{fmt.Sprintf("[]*domain.%s", repo.Entity), "error"},
		},
	}

	for _, cm := range repo.CustomMethods {
		method := template.MockMethod{
			Name:   cm.Name,
			Params: make([]template.MethodParam, 0, len(cm.Params)+1),
			Return: cm.Returns,
		}

		method.Params = append(method.Params, template.MethodParam{
			Name: "ctx",
			Type: "context.Context",
		})

		for _, p := range cm.Params {
			method.Params = append(method.Params, template.MethodParam{
				Name: p.Name,
				Type: p.Type,
			})
		}

		methods = append(methods, method)
	}

	return methods
}
