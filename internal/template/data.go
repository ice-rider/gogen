package template

import (
	"gogen/pkg/models"
)

type EntityData struct {
	Name          string
	Fields        []models.Field
	TableName     string
	ModulePath    string
	AddComments   bool
	AddValidation bool
	JSONStyle     string
}

type RepositoryData struct {
	Name             string
	Entity           string
	TableName        string
	ModulePath       string
	DBType           string
	CustomMethods    []CustomMethod
	WithTransactions bool
	AddComments      bool
	Fields           []models.Field
}

type CustomMethod struct {
	Name    string
	Comment string
	Params  []MethodParam
	Return  string
}

type MethodParam struct {
	Name string
	Type string
}

type UseCaseData struct {
	Name         string
	Description  string
	ModulePath   string
	Dependencies []Dependency
	InputFields  []models.Field
	OutputFields []models.Field
	WithLogging  bool
	WithMetrics  bool
	AddComments  bool
	Example      string
}

type Dependency struct {
	Name  string
	Found bool
}

type MockData struct {
	Name       string
	Entity     string
	ModulePath string
	Methods    []MockMethod
}

type MockMethod struct {
	Name   string
	Params []MethodParam
	Return []string
}
