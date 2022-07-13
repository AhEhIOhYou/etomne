package application

import (
	"etomne/backend/domain/entities"
	"etomne/backend/domain/repository"
)

type fileApp struct {
	fl repository.FileRepository
}

func (f *fileApp) SaveFile(file *entities.File) (*entities.File, map[string]string) {
	return f.fl.SaveFile(file)
}

func (f *fileApp) GetFile(fileId uint64) (*entities.File, error) {
	return f.fl.GetFile(fileId)
}

func (f *fileApp) GetFilesByModel(modelId uint64) ([]entities.File, error) {
	return f.fl.GetFilesByModel(modelId)
}

func (f *fileApp) AddModelFile(modelFile *entities.ModelFile) (*entities.ModelFile, map[string]string) {
	return f.fl.AddModelFile(modelFile)
}

var _ FileAppInterface = &fileApp{}

type FileAppInterface interface {
	SaveFile(*entities.File) (*entities.File, map[string]string)
	GetFile(uint64) (*entities.File, error)
	GetFilesByModel(uint64) ([]entities.File, error)
	AddModelFile(*entities.ModelFile) (*entities.ModelFile, map[string]string)
}
