package repository

import "github.com/AhEhIOhYou/etomne/backend/domain/entities"

type FileRepository interface {
	SaveFile(files *entities.File) (*entities.File, map[string]string)
	GetFile(uint64) (*entities.File, error)
	UpdateFile(file *entities.File) (*entities.File, map[string]string)
	DeleteFile(uint64) error
}
