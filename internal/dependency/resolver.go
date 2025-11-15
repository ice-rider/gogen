package dependency

import (
	"fmt"
	"strings"

	"gogen/pkg/models"
)

type Resolver struct {
	detector *Detector
}

func NewResolver(detector *Detector) *Resolver {
	return &Resolver{
		detector: detector,
	}
}

func (r *Resolver) Resolve(plan *models.GenerationPlan) error {

	for i := range plan.UseCases {
		uc := &plan.UseCases[i]

		deps := r.detector.DetectUseCaseDependencies(uc, plan)
		uc.Dependencies = deps

		for _, dep := range deps {
			if !dep.Found {

				if dep.Type == "repository" {
					if err := r.autoCreateRepository(dep.Name, plan); err != nil {
						return fmt.Errorf("cannot resolve dependency %s for %s: %w",
							dep.Name, uc.Name, err)
					}

					dep.Found = true
				}
			}
		}
	}

	return nil
}

func (r *Resolver) autoCreateRepository(repoName string, plan *models.GenerationPlan) error {

	entityName := strings.TrimSuffix(repoName, "Repository")

	if !plan.HasEntity(entityName) {
		return fmt.Errorf("entity %s not found, cannot create repository", entityName)
	}

	entity := plan.GetEntityByName(entityName)

	repo := models.RepositoryConfig{
		Name:          entityName,
		Entity:        entityName,
		TableName:     entity.TableName,
		DBType:        "postgres",
		CustomMethods: []models.CustomMethod{},
		AddComments:   true,
		Fields:        entity.Fields,
	}

	plan.Repositories = append(plan.Repositories, repo)

	return nil
}

func (r *Resolver) ValidatePlan(plan *models.GenerationPlan) error {
	missing := r.detector.DetectMissingDependencies(plan)

	if len(missing) > 0 {
		var errors []string
		for ucName, deps := range missing {
			errors = append(errors, fmt.Sprintf("UseCase %s: missing %v", ucName, deps))
		}
		return fmt.Errorf("unresolved dependencies:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}
