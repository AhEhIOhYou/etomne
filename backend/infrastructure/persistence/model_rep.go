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
	//The images are uploaded to digital ocean spaces. So we need to prepend the url. This might not be your use case, if you are not uploading image to Digital Ocean.

	err := r.db.Debug().Create(&model).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "model title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return model, nil
}

func (r *ModelRepo) GetModel(id uint64) (*entities.Model, error) {
	var model entities.Model
	err := r.db.Debug().Where("id = ?", id).Take(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("model not found")
	}
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	return &model, nil
}

func (r *ModelRepo) GetAllModel() ([]entities.Model, error) {
	var models []entities.Model
	err := r.db.Debug().Limit(100).Order("created_at desc").Find(&models).Error
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
	err := r.db.Debug().Save(&model).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return model, nil
}

func (r *ModelRepo) DeleteModel(id uint64) error {
	var model entities.Model
	err := r.db.Debug().Where("id = ?", id).Delete(&model).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

func (r *ModelRepo) GetFilesByModel(modelId uint64) ([]entities.File, error) {
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

func (r *ModelRepo) AddModelFile(mf *entities.ModelFile) (*entities.ModelFile, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Create(&mf).Error
	if err != nil {
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return mf, nil
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

func (r *ModelRepo) CheckAvailability(fileId uint64, userId uint64) (bool, error) {
	var result int
	rows := r.db.Table("models as m").
		Select("COUNT(f.id)").
		Joins("join model_files as mf on mf.model_id = m.id").
		Joins("join files as f on f.id = mf.file_id").
		Where("m.user_id = ? AND f.id = ?", userId, fileId).Limit(1).Row()

	if err := rows.Scan(&result); err != nil {
		return false, err
	}

	return result == 1, nil
}
