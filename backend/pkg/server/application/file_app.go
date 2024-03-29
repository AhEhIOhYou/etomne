package application

import (
	"github.com/AhEhIOhYou/etomne/pkg/server/domain/entities"
	"github.com/AhEhIOhYou/etomne/pkg/server/domain/repository"
)

type fileApp struct {
	fl repository.FileRepository
}

var _ FileAppInterface = &fileApp{}

type FileAppInterface interface {
	SaveFile(*entities.File) (*entities.File, error)
	GetFile(uint64) (*entities.File, error)
	UpdateFile(file *entities.File) (*entities.File, error)
	DeleteFile(uint64) error
	CheckAvailabilityFile(uint64, uint64) (bool, error)
}

func (f *fileApp) SaveFile(file *entities.File) (*entities.File, error) {
	return f.fl.SaveFile(file)
}

func (f *fileApp) GetFile(fileId uint64) (*entities.File, error) {
	return f.fl.GetFile(fileId)
}

func (f *fileApp) UpdateFile(file *entities.File) (*entities.File, error) {
	return f.fl.UpdateFile(file)
}

func (f *fileApp) DeleteFile(fileId uint64) error {
	return f.fl.DeleteFile(fileId)
}

func (f *fileApp) CheckAvailabilityFile(fileId uint64, userId uint64) (bool, error) {
	return f.fl.CheckAvailabilityFile(fileId, userId)
}
