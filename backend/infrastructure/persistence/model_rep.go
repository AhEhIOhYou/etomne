package persistence

import (
	"errors"
	"github.com/AhEhIOhYou/etomne/backend/domain/entities"
	"github.com/AhEhIOhYou/etomne/backend/domain/repository"
	"gorm.io/gorm"
	"strings"
)

type ModelRepo struct {
	db *gorm.DB
}

var _ repository.ModelRepository = &ModelRepo{}

func NewModelRepo(db *gorm.DB) *ModelRepo {
	return &ModelRepo{
		db: db,
	}
}

func (r *ModelRepo) SaveModel(model *entities.Model) (*entities.Model, error) {
	err := r.db.Debug().Table("model").Create(&model).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			return nil, err
		}
		return nil, err
	}

	return model, nil
}

func (r *ModelRepo) GetModel(id uint64) (*entities.Model, error) {
	var model *entities.Model
	err := r.db.Debug().Table("model").Where("id = ?", id).Take(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("model not found")
	}
	if err != nil {
		return nil, errors.New("database error, please try again")
	}

	return model, nil
}

func (r *ModelRepo) GetAllModel(page, limit int) ([]entities.Model, error) {
	var models []entities.Model
	offset := (page - 1) * limit
	err := r.db.Debug().Table("model").Limit(limit).Offset(offset).Order("created_at desc").Find(&models).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("model not found")
	}
	if err != nil {
		return nil, err
	}

	return models, nil
}

func (r *ModelRepo) UpdateModel(model *entities.Model) (*entities.Model, error) {
	err := r.db.Debug().Table("model").Save(&model).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			return nil, err
		}
		return nil, err
	}

	return model, nil
}

func (r *ModelRepo) DeleteModel(id uint64) error {
	var model entities.Model
	err := r.db.Debug().Table("model").Where("id = ?", id).Delete(&model).Error
	if err != nil {
		return errors.New("database error, please try again")
	}

	return nil
}

func (r *ModelRepo) GetFilesByModel(modelId uint64) ([]entities.File, error) {
	var files []entities.File
	err := r.db.
		Table("file").
		Joins("JOIN model_files on model_files.file_id=file.id").
		Where("model_files.model_id = ?", modelId).
		Find(&files).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (r *ModelRepo) SaveModelFile(file *entities.File, modelId uint64) (*entities.File, error) {
	err := r.db.Debug().Table("file").Create(&file).Error
	if err != nil {
		return nil, err
	}
	modelFile := entities.ModelFile{
		ModelId: modelId,
		FileId:  file.ID,
	}
	err = r.db.Debug().Create(&modelFile).Error
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (r *ModelRepo) DeleteModelFile(fileId uint64) error {
	var fModel entities.ModelFile
	err := r.db.Debug().Table("model_files").Where("file_id = ?", fileId).Delete(&fModel).Error
	if err != nil {
		return errors.New("database error, please try again")
	}

	return nil
}

func (r *ModelRepo) DeleteAllModelFiles(modelId uint64) error {
	var fModel entities.ModelFile
	err := r.db.Debug().Table("model_files").Where("model_id = ?", modelId).Delete(&fModel).Error
	if err != nil {
		return errors.New("database error, please try again")
	}

	return nil
}

func (r *ModelRepo) CheckAvailabilityModel(modelId uint64, userId uint64) (bool, error) {
	var result int
	rows := r.db.
		Table("model").
		Select("COUNT(model.id)").
		Where("model.id = ? AND model.user_id = ?", modelId, userId).Limit(1).Row()

	if err := rows.Scan(&result); err != nil {
		return false, err
	}

	return result == 1, nil
}
