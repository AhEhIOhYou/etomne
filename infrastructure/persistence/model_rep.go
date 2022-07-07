package persistence

import (
	"errors"
	"etomne/domain/entities"
	"etomne/domain/repository"
	"gorm.io/gorm"
	"os"
	"strings"
)

type ModelRepo struct {
	db *gorm.DB
}

func (r *ModelRepo) SaveModel(model *entities.Model) (*entities.Model, map[string]string) {
	dbErr := map[string]string{}
	//The images are uploaded to digital ocean spaces. So we need to prepend the url. This might not be your use case, if you are not uploading image to Digital Ocean.
	model.File = os.Getenv("DO_SPACES_URL") + model.File

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
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("model not found")
	}
	return &model, nil
}

func (r *ModelRepo) GetAllModels() ([]entities.Model, error) {
	var models []entities.Model
	err := r.db.Debug().Limit(100).Order("created_at desc").Find(&models).Error
	if err != nil {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("model not found")
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

func NewModelRepo(db *gorm.DB) *ModelRepo {
	return &ModelRepo{
		db: db,
	}
}

var _ repository.ModelRepository = &ModelRepo{}
