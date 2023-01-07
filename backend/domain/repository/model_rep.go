package repository

import (
	"github.com/AhEhIOhYou/etomne/backend/domain/entities"
)

type ModelRepository interface {
	SaveModel(*entities.Model) (*entities.Model, map[string]string)
	GetModel(uint64) (*entities.Model, error)
	GetAllModel(int, int) ([]entities.Model, error)
	UpdateModel(*entities.Model) (*entities.Model, map[string]string)
	DeleteModel(uint64) error

	GetFilesByModel(uint64) ([]entities.File, error)
	SaveModelFile(*entities.File, uint64) (*entities.ModelFile, map[string]string)
	DeleteModelFile(uint64) error
	DeleteAllModelFiles(uint64) error

	CheckAvailabilityModel(uint64, uint64) (bool, error)
}
