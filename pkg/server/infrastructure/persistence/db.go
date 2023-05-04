package persistence

import (
	"fmt"

	"github.com/AhEhIOhYou/etomne/pkg/server/domain/entities"
	"github.com/AhEhIOhYou/etomne/pkg/server/domain/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repos struct {
	User  repository.UserRepository
	Model repository.ModelRepository
	File  repository.FileRepository
	db    *gorm.DB
}

func NewRepo(DbUser, DbPassword, DbPort, DbHost, DbName string) (*Repos, error) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		DbHost, DbUser, DbPassword, DbName, DbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	db.Logger.LogMode(0)

	return &Repos{
		User:  NewUserRepo(db),
		Model: NewModelRepo(db),
		File:  NewFileRepo(db),
		db:    db,
	}, nil
}

func (s *Repos) Migrate() string {
	return s.db.AutoMigrate(&entities.User{}, &entities.Model{}).Error()
}
