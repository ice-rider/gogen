package models

type EntityConfig struct {
	Name          string  `json:"name"`
	Fields        []Field `json:"fields"`
	TableName     string  `json:"table_name"`
	AddValidation bool    `json:"add_validation"`
	AddComments   bool    `json:"add_comments"`
	JSONStyle     string  `json:"json_style"`
}

func (e *EntityConfig) GetName() string {
	return e.Name
}

func (e *EntityConfig) GetType() ComponentType {
	return ComponentTypeEntity
}

func (e *EntityConfig) HasTimeField() bool {
	for _, field := range e.Fields {
		if field.Type == "time.Time" || field.BaseType() == "time.Time" {
			return true
		}
	}
	return false
}

func (e *EntityConfig) HasUUIDField() bool {
	for _, field := range e.Fields {
		if field.Type == "uuid.UUID" || field.BaseType() == "uuid.UUID" {
			return true
		}
	}
	return false
}

func (e *EntityConfig) GetRequiredFields() []Field {
	var required []Field
	for _, field := range e.Fields {
		if field.Required {
			required = append(required, field)
		}
	}
	return required
}
