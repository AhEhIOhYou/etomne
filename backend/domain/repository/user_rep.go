package repository

import (
	"github.com/AhEhIOhYou/etomne/backend/domain/entities"
)

type UserRepository interface {
	SaveUser(*entities.User) (*entities.User, map[string]string)
	GetUser(uint64) (*entities.User, error)
	GetUsers(uint64) ([]entities.User, error)
	GetUserByEmailAndPassword(*entities.User) (*entities.User, map[string]string)
}
