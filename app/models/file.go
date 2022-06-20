package models

import (
	"database/sql"
	"etomne/app/entities"
	"log"
)

func CreateFile(file *entities.File, db *sql.DB) (int64, error) {
	result, err := db.Exec("INSERT INTO models_db.files(path) values(?)", file.Path)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	lastId, _ := result.LastInsertId()
	return lastId, nil
}
