package application

import (
	"github.com/AhEhIOhYou/etomne/backend/domain/entities"
	"github.com/AhEhIOhYou/etomne/backend/domain/repository"
)

type userApp struct {
	us repository.UserRepository
}

var _ UserAppInterface = &userApp{}

type UserAppInterface interface {
	SaveUser(*entities.User) (*entities.User, map[string]string)
	GetUser(uint64) (*entities.User, error)
	GetUsers(uint64) ([]entities.User, error)
	GetUserByEmailAndPassword(*entities.User) (*entities.User, map[string]string)

	GetPhotosByUser(uint64) ([]entities.File, error)
	AddUserPhoto(*entities.UserPhoto) (*entities.UserPhoto, map[string]string)
	DeleteUserPhoto(uint64) error
	DeleteAllUserPhotos(uint64) error
}

func (u *userApp) SaveUser(user *entities.User) (*entities.User, map[string]string) {
	return u.us.SaveUser(user)
}

func (u *userApp) GetUser(userId uint64) (*entities.User, error) {
	return u.us.GetUser(userId)
}

func (u *userApp) GetUsers(count uint64) ([]entities.User, error) {
	return u.us.GetUsers(count)
}

func (u *userApp) GetUserByEmailAndPassword(user *entities.User) (*entities.User, map[string]string) {
	return u.us.GetUserByEmailAndPassword(user)
}
func (u *userApp) GetPhotosByUser(userId uint64) ([]entities.File, error) {
	return u.us.GetPhotosByUser(userId)
}

func (u *userApp) AddUserPhoto(photo *entities.UserPhoto) (*entities.UserPhoto, map[string]string) {
	return u.us.AddUserPhoto(photo)
}

func (u *userApp) DeleteUserPhoto(fileId uint64) error {
	return u.us.DeleteUserPhoto(fileId)
}

func (u *userApp) DeleteAllUserPhotos(userId uint64) error {
	return u.us.DeleteAllUserPhotos(userId)
}
