package dependency

import (
	"strings"

	"gogen/pkg/models"
)

type Detector struct{}

func NewDetector() *Detector {
	return &Detector{}
}

func (d *Detector) DetectUseCaseDependencies(uc *models.UseCaseConfig, plan *models.GenerationPlan) []models.Dependency {
	var deps []models.Dependency

	entityName := d.extractEntityFromUseCaseName(uc.Name)

	if entityName != "" {

		repoName := entityName + "Repository"
		found := plan.HasRepository(entityName)

		deps = append(deps, models.Dependency{
			Name:  repoName,
			Type:  "repository",
			Found: found,
		})
	}

	for _, dep := range uc.Dependencies {

		exists := false
		for _, existing := range deps {
			if existing.Name == dep.Name {
				exists = true
				break
			}
		}

		if !exists {
			deps = append(deps, dep)
		}
	}

	return deps
}

func (d *Detector) extractEntityFromUseCaseName(name string) string {

	name = strings.TrimSuffix(name, "UseCase")

	prefixes := []string{
		"Create", "Get", "Update", "Delete", "List",
		"Find", "Search", "Fetch", "Remove", "Add",
		"Register", "Login", "Logout", "Activate",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			entityName := strings.TrimPrefix(name, prefix)
			if entityName != "" {
				return entityName
			}
		}
	}

	return ""
}

func (d *Detector) DetectMissingDependencies(plan *models.GenerationPlan) map[string][]string {
	missing := make(map[string][]string)

	for _, uc := range plan.UseCases {
		deps := d.DetectUseCaseDependencies(&uc, plan)

		var missingForUC []string
		for _, dep := range deps {
			if !dep.Found {
				missingForUC = append(missingForUC, dep.Name)
			}
		}

		if len(missingForUC) > 0 {
			missing[uc.Name] = missingForUC
		}
	}

	return missing
}
