package repository

import "github.com/AhEhIOhYou/etomne/backend/domain/entities"

type FileRepository interface {
	SaveFile(*entities.File) (*entities.File, map[string]string)
	GetFile(uint64) (*entities.File, error)
	UpdateFile(*entities.File) (*entities.File, map[string]string)
	DeleteFile(uint64) error
	CheckAvailabilityFile(uint64, uint64) (bool, error)
}
