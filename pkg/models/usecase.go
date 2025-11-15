package models

type UseCaseConfig struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Dependencies []Dependency `json:"dependencies"`
	InputFields  []Field      `json:"input_fields"`
	OutputFields []Field      `json:"output_fields"`
	WithLogging  bool         `json:"with_logging"`
	WithMetrics  bool         `json:"with_metrics"`
	AddComments  bool         `json:"add_comments"`
	Example      string       `json:"example"`
}

func (u *UseCaseConfig) GetName() string {
	return u.Name
}

func (u *UseCaseConfig) GetType() ComponentType {
	return ComponentTypeUseCase
}

type Dependency struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Found bool   `json:"found"`
}

func (u *UseCaseConfig) GetMissingDependencies() []string {
	var missing []string
	for _, dep := range u.Dependencies {
		if !dep.Found {
			missing = append(missing, dep.Name)
		}
	}
	return missing
}

func (u *UseCaseConfig) HasDependency(name string) bool {
	for _, dep := range u.Dependencies {
		if dep.Name == name {
			return true
		}
	}
	return false
}
