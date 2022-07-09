package application

import (
	"etomne/domain/entities"
	"etomne/domain/repository"
)

type modelApp struct {
	md repository.ModelRepository
}

func (m *modelApp) SaveModel(model *entities.Model) (*entities.Model, map[string]string) {
	return m.md.SaveModel(model)
}

func (m *modelApp) GetAllModel() ([]entities.Model, error) {
	return m.md.GetAllModels()
}

func (m *modelApp) GetModel(modelId uint64) (*entities.Model, error) {
	return m.md.GetModel(modelId)
}

func (m *modelApp) UpdateModel(model *entities.Model) (*entities.Model, map[string]string) {
	return m.md.UpdateModel(model)
}

func (m *modelApp) DeleteModel(modelId uint64) error {
	return m.md.DeleteModel(modelId)
}

var _ ModelAppInterface = &modelApp{}

type ModelAppInterface interface {
	SaveModel(*entities.Model) (*entities.Model, map[string]string)
	GetAllModel() ([]entities.Model, error)
	GetModel(uint64) (*entities.Model, error)
	UpdateModel(*entities.Model) (*entities.Model, map[string]string)
	DeleteModel(uint64) error
}
