package template

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"gogen/pkg/models"
)

var templatesFS embed.FS

type Loader struct {
	projectRoot string
	config      *models.Config
	cache       map[string]*template.Template
}

func NewLoader(projectRoot string, config *models.Config) *Loader {
	return &Loader{
		projectRoot: projectRoot,
		config:      config,
		cache:       make(map[string]*template.Template),
	}
}

func (l *Loader) Load(templateName string) (*template.Template, error) {

	if tmpl, ok := l.cache[templateName]; ok {
		return tmpl, nil
	}

	templatePath := l.getTemplatePath(templateName)

	var tmpl *template.Template
	var err error

	if l.isCustomTemplate(templatePath) {
		tmpl, err = l.loadFromFile(templatePath)
	} else {

		tmpl, err = l.loadFromEmbed(templatePath)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load template %s: %w", templateName, err)
	}

	l.cache[templateName] = tmpl

	return tmpl, nil
}

func (l *Loader) getTemplatePath(name string) string {
	switch name {
	case "entity":
		return l.config.Templates.Entity
	case "repository_interface":
		return l.config.Templates.RepositoryInterface
	case "repository_impl":
		return l.config.Templates.RepositoryImpl
	case "usecase":
		return l.config.Templates.UseCase
	case "handler":
		return l.config.Templates.Handler
	case "mock":
		return l.config.Templates.Mock
	case "test_entity":
		return l.config.Templates.TestEntity
	case "test_repository":
		return l.config.Templates.TestRepository
	case "test_usecase":
		return l.config.Templates.TestUseCase
	default:
		return ""
	}
}

func (l *Loader) isCustomTemplate(path string) bool {

	if filepath.IsAbs(path) {
		return true
	}

	fullPath := filepath.Join(l.projectRoot, path)
	_, err := os.Stat(fullPath)
	return err == nil
}

func (l *Loader) loadFromFile(path string) (*template.Template, error) {
	fullPath := path
	if !filepath.IsAbs(path) {
		fullPath = filepath.Join(l.projectRoot, path)
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(filepath.Base(path)).
		Funcs(l.getFuncMap()).
		Parse(string(data))

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func (l *Loader) loadFromEmbed(path string) (*template.Template, error) {
	data, err := templatesFS.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(filepath.Base(path)).
		Funcs(l.getFuncMap()).
		Parse(string(data))

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func (l *Loader) ClearCache() {
	l.cache = make(map[string]*template.Template)
}
