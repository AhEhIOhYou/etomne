package application

import (
	"etomne/backend/domain/entities"
	"etomne/backend/domain/repository"
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

var _ UserAppInterface = &userApp{}

type UserAppInterface interface {
	SaveUser(*entities.User) (*entities.User, map[string]string)
	GetUser(uint64) (*entities.User, error)
	GetUsers() ([]entities.User, error)
	GetUserByEmailAndPassword(*entities.User) (*entities.User, map[string]string)
}
