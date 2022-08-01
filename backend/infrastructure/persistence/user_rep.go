package persistence

import (
	"errors"
	"github.com/AhEhIOhYou/etomne/backend/domain/entities"
	"github.com/AhEhIOhYou/etomne/backend/domain/repository"
	"github.com/AhEhIOhYou/etomne/backend/infrastructure/security"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
)

type UserRepo struct {
	db *gorm.DB
}

var _ repository.UserRepository = &UserRepo{}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) SaveUser(user *entities.User) (*entities.User, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Create(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["email_taken"] = "email already taken"
			return nil, dbErr
		}
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return user, nil
}

func (r *UserRepo) GetUser(id uint64) (*entities.User, error) {
	var user entities.User
	err := r.db.Debug().Where("id = ?", id).Take(&user).Error
	if err != nil {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *UserRepo) GetUsers(count uint64) ([]entities.User, error) {
	var users []entities.User
	err := r.db.Debug().Table("users").Limit(int(count)).Find(&users).Error
	if err != nil {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return users, nil
}

func (r *UserRepo) GetUserByEmailAndPassword(u *entities.User) (*entities.User, map[string]string) {
	var user entities.User
	dbErr := map[string]string{}
	err := r.db.Debug().Where("email = ?", u.Email).Take(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		dbErr["no_user"] = "user not found"
		return nil, dbErr
	}
	if err != nil {
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	err = security.VerifyPassword(user.Password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		dbErr["incorrect_password"] = "incorrect password"
		return nil, dbErr
	}
	return &user, nil
}

func (r *UserRepo) GetPhotosByUser(id uint64) ([]entities.File, error) {
	var photos []entities.File
	err := r.db.Debug().Joins("JOIN user_files on user_files.file_id=files.id").
		Where("user_files.user_id = ?", id).Find(&photos).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("photos not found")
	}
	if err != nil {
		return nil, err
	}
	return photos, nil
}

func (r *UserRepo) SaveUserPhoto(file *entities.File, userId, size uint64) (*entities.UserPhoto, map[string]string) {

	dbErr := map[string]string{}
	err := r.db.Debug().Create(&file).Error
	if err != nil {
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}

	userPhoto := entities.UserPhoto{
		UserId: userId,
		FileId: file.ID,
		Size:   size,
	}

	err = r.db.Debug().Create(&userPhoto).Error
	if err != nil {
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}

	return &userPhoto, nil
}

func (r *UserRepo) DeleteUserPhoto(fileId uint64) error {
	var pUser entities.UserPhoto
	err := r.db.Debug().Table("model_files").Where("file_id = ?", fileId).Delete(&pUser).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

func (r *UserRepo) DeleteAllUserPhotos(userId uint64) error {
	var pUser entities.UserPhoto
	err := r.db.Debug().Table("user_files").Where("user_id = ?", userId).Delete(&pUser).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}
