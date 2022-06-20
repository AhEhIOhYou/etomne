package models

import (
	"database/sql"
	"etomne/app/entities"
	"fmt"
	"log"
)

func GetAllModels(db *sql.DB) ([]entities.Model, error) {

	var models []entities.Model

	rows, _ := db.Query("select *from models")
	defer rows.Close()

	for rows.Next() {
		var model entities.Model
		if err := rows.Scan(&model.Id, &model.Name, &model.CreateDate, &model.Description, &model.FileId); err != nil {
			return models, fmt.Errorf("%v", err)
		}
		models = append(models, model)
	}

	if err := rows.Err(); err != nil {
		return models, fmt.Errorf("%v", err)
	}

	return models, nil
}
func GetModelById(id int, db *sql.DB) (entities.Model, error) {

	var model entities.Model

	rows, _ := db.Query("select * from models where models.id = ?", id)
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&model.Id, &model.Name, &model.CreateDate, &model.Description); err != nil {
			return model, fmt.Errorf("%v", err)
		}
	}

	if err := rows.Err(); err != nil {
		return model, fmt.Errorf("%v", err)
	}

	return model, nil
}
func CreateModel(model entities.Model, db *sql.DB) (int64, error) {
	result, err := db.Exec("insert into models(name, cre_date, descr, file_id) values(?, ?, ?, ?)",
		model.Name, model.CreateDate, model.Description, model.FileId)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	lastId, _ := result.LastInsertId()
	return lastId, nil
}
func EditModel(name string, createDate string, description string, db *sql.DB) (int64, error) {

	return 0, nil
}
func DeleteModel(id int, db *sql.DB) (int64, error) {

	return 0, nil
}
