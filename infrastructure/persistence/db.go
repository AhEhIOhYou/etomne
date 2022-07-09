package persistence

import (
	"etomne/domain/entities"
	"etomne/domain/repository"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Repos struct {
	User  repository.UserRepository
	Model repository.ModelRepository
	db    *gorm.DB
}

func NewRepo(DbUser, DbPassword, DbPort, DbHost, DbName string) (*Repos, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DbUser, DbPassword, DbHost, DbPort, DbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

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
