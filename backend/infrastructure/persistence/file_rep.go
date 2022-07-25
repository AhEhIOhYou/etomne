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

var _ repository.FileRepository = &FileRepo{}

func NewFileRepo(db *gorm.DB) *FileRepo {
	return &FileRepo{
		db: db,
	}
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

func (r *FileRepo) UpdateFile(file *entities.File) (*entities.File, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Save(&file).Error
	if err != nil {
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return file, nil
}

func (r *FileRepo) DeleteFile(id uint64) error {
	var file entities.File
	err := r.db.Debug().Where("id = ?", id).Delete(&file).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}
