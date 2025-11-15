package project

import (
	"fmt"
	"os"
	"path/filepath"

	"gogen/pkg/models"
)

func (f *Finder) EnsureStructure(config *models.Config) error {
	root, err := f.FindRoot()
	if err != nil {
		return err
	}

	paths := []string{
		config.Paths.Domain,
		config.Paths.Repository,
		config.Paths.UseCase,
		config.Paths.Handler,
		config.Paths.Mocks,
		config.Paths.Tests,
	}

	for _, p := range paths {
		fullPath := filepath.Join(root, p)

		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", fullPath, err)
		}
	}

	return nil
}

func (f *Finder) CheckStructure(config *models.Config) ([]string, error) {
	root, err := f.FindRoot()
	if err != nil {
		return nil, err
	}

	var missing []string

	paths := map[string]string{
		"domain":     config.Paths.Domain,
		"repository": config.Paths.Repository,
		"usecase":    config.Paths.UseCase,
		"handler":    config.Paths.Handler,
		"mocks":      config.Paths.Mocks,
		"tests":      config.Paths.Tests,
	}

	for name, p := range paths {
		fullPath := filepath.Join(root, p)

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			missing = append(missing, name)
		}
	}

	return missing, nil
}
