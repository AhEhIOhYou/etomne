package repository

import (
	"github.com/AhEhIOhYou/etomne/backend/domain/entities"
)

type ModelRepository interface {
	SaveModel(*entities.Model) (*entities.Model, map[string]string)
	GetModel(uint642 uint64) (*entities.Model, error)
	GetAllModel() ([]entities.Model, error)
	UpdateModel(*entities.Model) (*entities.Model, map[string]string)
	DeleteModel(uint642 uint64) error
}
