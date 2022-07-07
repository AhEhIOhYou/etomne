package application

import (
	"etomne/domain/entities"
	"etomne/domain/repository"
)

type userApp struct {
	us repository.UserRepository
}

func (u *userApp) SaveUser(user *entities.User) (*entities.User, map[string]string) {
	return u.us.SaveUser(user)
}

func (u *userApp) GetUser(userId uint64) (*entities.User, error) {
	return u.us.GetUser(userId)
}

func (u *userApp) GetUsers() ([]entities.User, error) {
	return u.us.GetUsers()
}

func (u *userApp) GetUserByEmailAndPassword(user *entities.User) (*entities.User, map[string]string) {
	return u.us.GetUserByEmailAndPassword(user)
}

var _ UserApiInterface = &userApp{}

type UserApiInterface interface {
	SaveUser(*entities.User) (*entities.User, map[string]string)
	GetUser(uint64) (*entities.User, error)
	GetUsers() ([]entities.User, error)
	GetUserByEmailAndPassword(*entities.User) (*entities.User, map[string]string)
}
