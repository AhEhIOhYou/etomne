package application

import (
	"github.com/AhEhIOhYou/etomne/pkg/server/domain/entities"
	"github.com/AhEhIOhYou/etomne/pkg/server/domain/repository"
)

type userApp struct {
	us repository.UserRepository
}

var _ UserAppInterface = &userApp{}

type UserAppInterface interface {
	SaveUser(*entities.User) (*entities.User, error)
	GetUser(uint64) (*entities.User, error)
	GetUsers(uint64) ([]entities.User, error)
	GetUserByEmailAndPassword(*entities.User) (*entities.User, error)
}

func (u *userApp) SaveUser(user *entities.User) (*entities.User, error) {
	return u.us.SaveUser(user)
}

func (u *userApp) GetUser(userId uint64) (*entities.User, error) {
	return u.us.GetUser(userId)
}

func (u *userApp) GetUsers(count uint64) ([]entities.User, error) {
	return u.us.GetUsers(count)
}

func (u *userApp) GetUserByEmailAndPassword(user *entities.User) (*entities.User, error) {
	return u.us.GetUserByEmailAndPassword(user)
}
