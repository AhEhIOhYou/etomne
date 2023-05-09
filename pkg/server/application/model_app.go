package application

import (
	"github.com/AhEhIOhYou/etomne/pkg/server/domain/entities"
	"github.com/AhEhIOhYou/etomne/pkg/server/domain/repository"
)

type modelApp struct {
	md repository.ModelRepository
}

var _ ModelAppInterface = &modelApp{}

type ModelAppInterface interface {
	SaveModel(*entities.Model) (*entities.Model, error)
	GetAllModel(int, int) ([]entities.Model, error)
	GetModel(uint64) (*entities.Model, error)
	UpdateModel(*entities.Model) (*entities.Model, error)
	DeleteModel(uint64) error

	GetFilesByModel(uint64) ([]entities.File, error)
	SaveModelFile(*entities.File, uint64) (*entities.File, error)
	AddFileToModel(uint64, uint64) error
	DeleteModelFile(uint64) error
	DeleteAllModelFiles(uint64) error
}

func (m *modelApp) SaveModel(model *entities.Model) (*entities.Model, error) {
	return m.md.SaveModel(model)
}

func (m *modelApp) GetAllModel(page, limit int) ([]entities.Model, error) {
	return m.md.GetAllModel(page, limit)
}

func (m *modelApp) GetModel(modelId uint64) (*entities.Model, error) {
	return m.md.GetModel(modelId)
}

func (m *modelApp) UpdateModel(model *entities.Model) (*entities.Model, error) {
	return m.md.UpdateModel(model)
}

func (m *modelApp) DeleteModel(modelId uint64) error {
	return m.md.DeleteModel(modelId)
}

func (m *modelApp) GetFilesByModel(modelId uint64) ([]entities.File, error) {
	return m.md.GetFilesByModel(modelId)
}

func (m *modelApp) SaveModelFile(file *entities.File, modelId uint64) (*entities.File, error) {
	return m.md.SaveModelFile(file, modelId)
}

func (m *modelApp) AddFileToModel(fileID uint64, modelID uint64) error {
	return m.md.AddFileToModel(fileID, modelID)
}

func (m *modelApp) DeleteModelFile(fileId uint64) error {
	return m.md.DeleteModelFile(fileId)
}

func (m *modelApp) DeleteAllModelFiles(modelId uint64) error {
	return m.md.DeleteAllModelFiles(modelId)
}
