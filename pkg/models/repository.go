package models

type RepositoryConfig struct {
	Name             string         `json:"name"`
	Entity           string         `json:"entity"`
	TableName        string         `json:"table_name"`
	DBType           string         `json:"db_type"`
	CustomMethods    []CustomMethod `json:"custom_methods"`
	WithTransactions bool           `json:"with_transactions"`
	AddComments      bool           `json:"add_comments"`
	Fields           []Field        `json:"fields"`
}

func (r *RepositoryConfig) GetName() string {
	return r.Name
}

func (r *RepositoryConfig) GetType() ComponentType {
	return ComponentTypeRepository
}

type CustomMethod struct {
	Name    string        `json:"name"`
	Comment string        `json:"comment"`
	Params  []MethodParam `json:"params"`
	Returns []string      `json:"returns"`
	Body    string        `json:"body"`
}

type MethodParam struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
