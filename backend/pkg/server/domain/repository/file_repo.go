package repository

import "github.com/AhEhIOhYou/etomne/pkg/server/domain/entities"

type FileRepository interface {
	SaveFile(*entities.File) (*entities.File, error)
	GetFile(uint64) (*entities.File, error)
	UpdateFile(*entities.File) (*entities.File, error)
	DeleteFile(uint64) error
	CheckAvailabilityFile(uint64, uint64) (bool, error)
}
