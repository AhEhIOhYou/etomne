package persistence

import (
	"etomne/domain/entities"
	"etomne/domain/repository"
	"fmt"
	"gorm.io/gorm"
)

type Repos struct {
	User  repository.UserRepository
	Model repository.ModelRepository
	db    *gorm.DB
}

func NewRepo(DbDriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*Repos, error) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err := gorm.Open(DbDriver, DBURL)
	if err != nil {
		return nil, err
	}
	db.Logger.LogMode(0)

	return &Repos{
		User:  NewUserRepo(db),
		Model: NewModelRepo(db),
		db:    db,
	}, nil
}

func (s *Repos) Migrate() string {
	return s.db.AutoMigrate(&entities.User{}, &entities.Model{}).Error()
}
