package repository

import "etomne/domain/entities"

type ModelRepository interface {
	SaveModel(*entities.Model) (*entities.Model, map[string]string)
	GetModel(uint642 uint64) (*entities.Model, error)
	GetAllModels() ([]entities.Model, error)
	UpdateModel(*entities.Model) (*entities.Model, map[string]string)
	DeleteModel(uint642 uint64) error
}
