package models

type Config struct {
	Version    string     `yaml:"version"`
	Paths      Paths      `yaml:"paths"`
	Naming     Naming     `yaml:"naming"`
	Templates  Templates  `yaml:"templates"`
	Generation Generation `yaml:"generation"`
	Imports    Imports    `yaml:"imports"`
}

type Paths struct {
	Domain     string `yaml:"domain"`
	Repository string `yaml:"repository"`
	UseCase    string `yaml:"usecase"`
	Handler    string `yaml:"handler"`
	Mocks      string `yaml:"mocks"`
	Tests      string `yaml:"tests"`
}

type Naming struct {
	Style    string            `yaml:"style"`
	Suffixes map[string]string `yaml:"suffixes"`
	Prefixes map[string]string `yaml:"prefixes"`
}

type Templates struct {
	Entity              string `yaml:"entity"`
	RepositoryInterface string `yaml:"repository_interface"`
	RepositoryImpl      string `yaml:"repository_impl"`
	UseCase             string `yaml:"usecase"`
	Handler             string `yaml:"handler"`
	Mock                string `yaml:"mock"`
	TestEntity          string `yaml:"test_entity"`
	TestRepository      string `yaml:"test_repository"`
	TestUseCase         string `yaml:"test_usecase"`
}

type Generation struct {
	AddComments        bool   `yaml:"add_comments"`
	AddExamples        bool   `yaml:"add_examples"`
	SeparateInterfaces bool   `yaml:"separate_interfaces"`
	UsePointers        bool   `yaml:"use_pointers"`
	ErrorHandling      string `yaml:"error_handling"`
}

type Imports struct {
	Entity     []string `yaml:"entity"`
	Repository []string `yaml:"repository"`
	UseCase    []string `yaml:"usecase"`
	Test       []string `yaml:"test"`
}
