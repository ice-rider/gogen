package models

type GenerationPlan struct {
	Entities     []EntityConfig     `json:"entities"`
	Repositories []RepositoryConfig `json:"repositories"`
	UseCases     []UseCaseConfig    `json:"usecases"`
	Handlers     []HandlerConfig    `json:"handlers"`
	WithTests    bool               `json:"with_tests"`
	WithMocks    bool               `json:"with_mocks"`
	ModulePath   string             `json:"module_path"`
	ProjectRoot  string             `json:"project_root"`
}

type HandlerConfig struct {
	Name        string `json:"name"`
	UseCase     string `json:"usecase"`
	Route       string `json:"route"`
	Method      string `json:"method"`
	AddComments bool   `json:"add_comments"`
}

func (h *HandlerConfig) GetName() string {
	return h.Name
}

func (h *HandlerConfig) GetType() ComponentType {
	return ComponentTypeHandler
}

func (p *GenerationPlan) IsEmpty() bool {
	return len(p.Entities) == 0 &&
		len(p.Repositories) == 0 &&
		len(p.UseCases) == 0 &&
		len(p.Handlers) == 0
}

func (p *GenerationPlan) ComponentCount() int {
	return len(p.Entities) + len(p.Repositories) + len(p.UseCases) + len(p.Handlers)
}

func (p *GenerationPlan) GetEntityByName(name string) *EntityConfig {
	for i := range p.Entities {
		if p.Entities[i].Name == name {
			return &p.Entities[i]
		}
	}
	return nil
}

func (p *GenerationPlan) GetRepositoryByName(name string) *RepositoryConfig {
	for i := range p.Repositories {
		if p.Repositories[i].Name == name {
			return &p.Repositories[i]
		}
	}
	return nil
}

func (p *GenerationPlan) HasEntity(name string) bool {
	return p.GetEntityByName(name) != nil
}

func (p *GenerationPlan) HasRepository(name string) bool {
	return p.GetRepositoryByName(name) != nil
}
