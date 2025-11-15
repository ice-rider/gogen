package models

type Component interface {
	GetName() string
	GetType() ComponentType
}

type ComponentType string

const (
	ComponentTypeEntity     ComponentType = "entity"
	ComponentTypeRepository ComponentType = "repository"
	ComponentTypeUseCase    ComponentType = "usecase"
	ComponentTypeHandler    ComponentType = "handler"
	ComponentTypeMock       ComponentType = "mock"
	ComponentTypeTest       ComponentType = "test"
)
