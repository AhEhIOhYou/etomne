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

func (r *ModelRepo) SaveModel(model *entities.Model) (*entities.Model, map[string]string) {

	dbErr := map[string]string{}

	err := r.db.Debug().Table("model").Create(&model).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "model title already taken"
			return nil, dbErr
		}
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}

	return model, nil
}

func (r *ModelRepo) GetModel(id uint64) (*entities.Model, error) {

	var model entities.Model

	err := r.db.Debug().Table("model").Where("id = ?", id).Take(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("model not found")
	}
	if err != nil {
		return nil, errors.New("database error, please try again")
	}

	return &model, nil
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

func (r *ModelRepo) UpdateModel(model *entities.Model) (*entities.Model, map[string]string) {

	dbErr := map[string]string{}

	err := r.db.Debug().Table("model").Save(&model).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "title already taken"
			return nil, dbErr
		}
		dbErr["db_error"] = "database error"
		return nil, dbErr
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

func (r *ModelRepo) GetFilesByModel(modelId uint64) ([]entities.File,  map[string]string) {
	
	var files []entities.File
	dbErr := map[string]string{}
	
	err := r.db.Debug().Table("file").Joins("JOIN model_files on model_files.file_id=file.id").Where("model_files.model_id = ?", modelId).Find(&files).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		dbErr["db_error"] = "files not found"
		return nil, dbErr
	}
	if err != nil {
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	
	return files, nil
}

func (r *ModelRepo) SaveModelFile(file *entities.File, modelId uint64) (*entities.ModelFile, map[string]string) {

	dbErr := map[string]string{}
	err := r.db.Debug().Table("file").Create(&file).Error
	if err != nil {
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}

	modelFile := entities.ModelFile{
		ModelId: modelId,
		FileId:  file.ID,
	}

	err = r.db.Debug().Create(&modelFile).Error
	if err != nil {
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}

	return &modelFile, nil
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
