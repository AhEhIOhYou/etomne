package repository

import (
	"github.com/AhEhIOhYou/etomne/pkg/server/domain/entities"
)

type ModelRepository interface {
	SaveModel(*entities.Model) (*entities.Model, error)
	GetModel(uint64) (*entities.Model, error)
	GetAllModel(int, int) ([]entities.Model, error)
	UpdateModel(*entities.Model) (*entities.Model, error)
	DeleteModel(uint64) error

	GetFilesByModel(uint64) ([]entities.File, error)
	SaveModelFile(*entities.File, uint64) (*entities.File, error)
	DeleteModelFile(uint64) error
	DeleteAllModelFiles(uint64) error
}
