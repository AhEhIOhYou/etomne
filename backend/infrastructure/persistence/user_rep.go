package persistence

import (
	"errors"
	"strings"

	"github.com/AhEhIOhYou/etomne/backend/domain/entities"
	"github.com/AhEhIOhYou/etomne/backend/domain/repository"
	"github.com/AhEhIOhYou/etomne/backend/infrastructure/security"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

func (r *UserRepo) SaveUser(user *entities.User) (*entities.User, error) {
	err := r.db.Debug().Table("user").Create(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			return nil, err
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) GetUser(id uint64) (*entities.User, error) {
	var user entities.User
	err := r.db.Debug().Table("user").Where("id = ?", id).Take(&user).Error
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
	err := r.db.Debug().Table("user").Limit(int(count)).Find(&users).Error
	if err != nil {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}

	return users, nil
}

func (r *UserRepo) GetUserByEmailAndPassword(u *entities.User) (*entities.User, error) {
	var user entities.User
	err := r.db.Debug().Table("user").Where("email = ?", u.Email).Take(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	err = security.VerifyPassword(user.Password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, err
	}

	return &user, nil
}
