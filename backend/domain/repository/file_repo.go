package repository

import "etomne/backend/domain/entities"

type FileRepository interface {
	SaveFile(files *entities.File) (*entities.File, map[string]string)
	GetFile(uint64) (*entities.File, error)
	GetFilesByModel(uint64) ([]entities.File, error)
	AddModelFile(*entities.ModelFile) (*entities.ModelFile, map[string]string)
}
