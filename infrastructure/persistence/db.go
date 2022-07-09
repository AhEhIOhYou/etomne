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

func NewRepo(DbDriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*Repos, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		DbHost, DbPort, DbUser, DbName, DbPassword)
	//db, err := gorm.Open(mysql.Open(dsn), DBURL)

	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
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
