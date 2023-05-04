package persistence

import (
	"errors"

	"github.com/AhEhIOhYou/etomne/pkg/server/domain/entities"
	"github.com/AhEhIOhYou/etomne/pkg/server/domain/repository"
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

func (r *FileRepo) SaveFile(file *entities.File) (*entities.File, error) {
	err := r.db.Debug().Table("file").Create(&file).Error
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (r *FileRepo) GetFile(id uint64) (*entities.File, error) {
	var file entities.File
	err := r.db.Debug().Table("file").Where("id = ?", id).Take(&file).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("file not found")
	}

	return &file, nil
}

func (r *FileRepo) UpdateFile(file *entities.File) (*entities.File, error) {
	err := r.db.Debug().Table("file").Save(&file).Error
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (r *FileRepo) DeleteFile(id uint64) error {
	var file entities.File
	err := r.db.Debug().Table("file").Where("id = ?", id).Delete(&file).Error
	if err != nil {
		return errors.New("database error, please try again")
	}

	return nil
}

func (r *FileRepo) CheckAvailabilityFile(fileId uint64, userId uint64) (bool, error) {
	var result int
	rows := r.db.Table("file").
		Select("COUNT(id)").
		Where("owner_id = ?", userId, fileId).Limit(1).Row()
	if err := rows.Scan(&result); err != nil {
		return false, err
	}

	return result == 1, nil
}
