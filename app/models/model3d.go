package models

import (
	"database/sql"
	"etomne/app/entities"
	"log"
)

type Model3dModel struct {
	Db *sql.DB
}

func (model3dModel Model3dModel) GetAllModels() ([]entities.Model3d, error) {

	var models []entities.Model3d

	rows, _ := model3dModel.Db.Query("select m.*, f.path from models m join files f on f.id = m.file_id")
	defer rows.Close()

	for rows.Next() {
		var model entities.Model3d
		if err := rows.Scan(&model.Id, &model.Name, &model.CreateDate, &model.Description, &model.FileId, &model.FilePath); err != nil {
			return models, err
		}
		models = append(models, model)
	}

	if err := rows.Err(); err != nil {
		return models, err
	}

	return models, nil
}
func (model3dModel Model3dModel) GetModelById(id int) (entities.Model3d, error) {

	var model entities.Model3d

	rows, _ := model3dModel.Db.Query("select m.*, f.path from models m join files f on f.id = m.file_id where m.id = ?", id)
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&model.Id, &model.Name, &model.CreateDate, &model.Description, &model.FileId, &model.FilePath); err != nil {
			return model, err
		}
	}

	if err := rows.Err(); err != nil {
		return model, err
	}

	return model, nil
}
func (model3dModel Model3dModel) CreateModel(model entities.Model3d) (int64, error) {
	result, err := model3dModel.Db.Exec("insert into models(name, cre_date, descr, file_id) values(?, ?, ?, ?)",
		model.Name, model.CreateDate, model.Description, model.FileId)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	lastId, _ := result.LastInsertId()
	return lastId, nil
}
func (model3dModel Model3dModel) EditModel(model entities.Model3d) (int64, error) {
	result, err := model3dModel.Db.Exec("update models set name = ?, cre_date = ?, descr = ?, file_id = ? where id = ?",
		model.Name, model.CreateDate, model.Description, model.FileId, model.Id)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	lastId, _ := result.LastInsertId()
	return lastId, nil
}
func (model3dModel Model3dModel) DeleteModel(id int) (int64, error) {
	result, err := model3dModel.Db.Exec("delete from models where id = ?", id)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	lastId, _ := result.LastInsertId()
	return lastId, nil
}
