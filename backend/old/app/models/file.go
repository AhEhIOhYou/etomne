package models

import (
	"database/sql"
	"etomne/app/entities"
	"log"
)

func CreateFile(file *entities.File, db *sql.DB) (int64, error) {
	result, err := db.Exec("insert into models_db.files(path) values(?)", file.Path)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	lastId, _ := result.LastInsertId()
	return lastId, nil
}
func DeleteFile(id int64, db *sql.DB) (int64, error) {
	result, err := db.Exec("delete from files where id = ?", id)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	lastId, _ := result.LastInsertId()
	return lastId, nil
}
