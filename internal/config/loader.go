package config

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"gogen/pkg/models"

	"gopkg.in/yaml.v3"
)

var globalConfigFS embed.FS

type Loader struct {
	projectRoot string
}

func NewLoader(projectRoot string) *Loader {
	return &Loader{
		projectRoot: projectRoot,
	}
}

func (l *Loader) Load() (*models.Config, error) {

	globalConfig, err := l.loadGlobalConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load global config: %w", err)
	}

	userConfig, err := l.loadUserConfig()
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load user config: %w", err)
	}

	finalConfig := l.mergeConfigs(globalConfig, userConfig)

	if err := l.validateConfig(finalConfig); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return finalConfig, nil
}

func (l *Loader) loadGlobalConfig() (*models.Config, error) {
	data, err := globalConfigFS.ReadFile("global.yaml")
	if err != nil {
		return nil, err
	}

	var config models.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (l *Loader) loadUserConfig() (*models.Config, error) {
	configPath := filepath.Join(l.projectRoot, "gogen.yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config models.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
