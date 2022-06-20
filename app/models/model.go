package models

import (
	"database/sql"
	"etomne/app/entities"
	"fmt"
)

func GetAllModels(db *sql.DB) ([]entities.Model, error) {

	var models []entities.Model

	rows, _ := db.Query("SELECT * FROM models")
	defer rows.Close()

	for rows.Next() {
		var model entities.Model
		if err := rows.Scan(&model.Id, &model.Name, &model.CreateDate, &model.Description); err != nil {
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

	rows, _ := db.Query("SELECT * FROM models WHERE models.id = ?", id)
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
func EditModel(name string, createDate string, description string, db *sql.DB) (bool, error) {

	return true, nil
}
func CreateModel(name string, createDate string, description string, db *sql.DB) (bool, error) {

	return true, nil
}
func DeleteModel(id int, db *sql.DB) (bool, error) {

	return true, nil
}
