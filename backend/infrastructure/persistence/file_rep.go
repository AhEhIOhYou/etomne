package persistence

import (
	"errors"
	"github.com/AhEhIOhYou/etomne/backend/domain/entities"
	"github.com/AhEhIOhYou/etomne/backend/domain/repository"
	"gorm.io/gorm"
)

type FileRepo struct {
	db *gorm.DB
}

func (r *FileRepo) SaveFile(file *entities.File) (*entities.File, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Create(&file).Error
	if err != nil {
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return file, nil
}

func (r *FileRepo) GetFile(id uint64) (*entities.File, error) {
	var file entities.File
	err := r.db.Debug().Where("id = ?", id).Take(&file).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("file not found")
	}
	return &file, nil
}

func (r *FileRepo) GetFilesByModel(modelId uint64) ([]entities.File, error) {
	var files []entities.File
	err := r.db.Debug().Joins("JOIN model_files on model_files.file_id=files.id").Where("model_files.model_id = ?", modelId).Find(&files).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("files not found")
	}
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (r *FileRepo) AddModelFile(mf *entities.ModelFile) (*entities.ModelFile, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Create(&mf).Error
	if err != nil {
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return mf, nil
}

func NewFileRepo(db *gorm.DB) *FileRepo {
	return &FileRepo{
		db: db,
	}
}

var _ repository.FileRepository = &FileRepo{}
