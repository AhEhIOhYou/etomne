package repository

import "etomne/domain/entities"

type UserRepository interface {
	SaveUser(*entities.User) (*entities.User, map[string]string)
	GetUser(uint64) (*entities.User, error)
	GetUsers() ([]entities.User, error)
	GetUserByEmailAndPassword(*entities.User) (*entities.User, map[string]string)
}
