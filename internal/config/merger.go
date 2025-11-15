package config

import (
	"gogen/pkg/models"
)

func (l *Loader) mergeConfigs(global, user *models.Config) *models.Config {
	if user == nil {
		return global
	}

	result := *global

	if user.Version != "" {
		result.Version = user.Version
	}

	result.Paths = l.mergePaths(global.Paths, user.Paths)

	result.Naming = l.mergeNaming(global.Naming, user.Naming)

	result.Templates = l.mergeTemplates(global.Templates, user.Templates)

	result.Generation = l.mergeGeneration(global.Generation, user.Generation)

	result.Imports = l.mergeImports(global.Imports, user.Imports)

	return &result
}

func (l *Loader) mergePaths(global, user models.Paths) models.Paths {
	result := global

	if user.Domain != "" {
		result.Domain = user.Domain
	}
	if user.Repository != "" {
		result.Repository = user.Repository
	}
	if user.UseCase != "" {
		result.UseCase = user.UseCase
	}
	if user.Handler != "" {
		result.Handler = user.Handler
	}
	if user.Mocks != "" {
		result.Mocks = user.Mocks
	}
	if user.Tests != "" {
		result.Tests = user.Tests
	}

	return result
}

func (l *Loader) mergeNaming(global, user models.Naming) models.Naming {
	result := global

	if user.Style != "" {
		result.Style = user.Style
	}

	if user.Suffixes != nil {
		if result.Suffixes == nil {
			result.Suffixes = make(map[string]string)
		}
		for k, v := range user.Suffixes {
			result.Suffixes[k] = v
		}
	}

	if user.Prefixes != nil {
		if result.Prefixes == nil {
			result.Prefixes = make(map[string]string)
		}
		for k, v := range user.Prefixes {
			result.Prefixes[k] = v
		}
	}

	return result
}

func (l *Loader) mergeTemplates(global, user models.Templates) models.Templates {
	result := global

	if user.Entity != "" {
		result.Entity = user.Entity
	}
	if user.RepositoryInterface != "" {
		result.RepositoryInterface = user.RepositoryInterface
	}
	if user.RepositoryImpl != "" {
		result.RepositoryImpl = user.RepositoryImpl
	}
	if user.UseCase != "" {
		result.UseCase = user.UseCase
	}
	if user.Handler != "" {
		result.Handler = user.Handler
	}
	if user.Mock != "" {
		result.Mock = user.Mock
	}
	if user.TestEntity != "" {
		result.TestEntity = user.TestEntity
	}
	if user.TestRepository != "" {
		result.TestRepository = user.TestRepository
	}
	if user.TestUseCase != "" {
		result.TestUseCase = user.TestUseCase
	}

	return result
}

func (l *Loader) mergeGeneration(global, user models.Generation) models.Generation {
	result := global

	if user.ErrorHandling != "" {
		result.ErrorHandling = user.ErrorHandling
	}

	return result
}

func (l *Loader) mergeImports(global, user models.Imports) models.Imports {
	result := global

	if len(user.Entity) > 0 {
		result.Entity = user.Entity
	}
	if len(user.Repository) > 0 {
		result.Repository = user.Repository
	}
	if len(user.UseCase) > 0 {
		result.UseCase = user.UseCase
	}
	if len(user.Test) > 0 {
		result.Test = user.Test
	}

	return result
}
