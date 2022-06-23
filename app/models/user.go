package models

import (
	"database/sql"
	"etomne/app/entities"
)

type UserModel struct {
	Db *sql.DB
}

func (userModel *UserModel) Login(login string, pass string) (entities.User, error) {

	var user entities.User

	row, err := userModel.Db.Query("select id, login, name from users where login = ? and hashed_pass = ? ", login, pass)
	if err != nil {
		return user, err
	}

	defer row.Close()

	if row.Next() {
		if err := row.Scan(&user.Id, &user.Login, &user.Name); err != nil {
			return user, err
		}
	}

	if err := row.Err(); err != nil {
		return user, err
	}

	return user, err
}
func (userModel *UserModel) Create(login string, pass string, name string) (int64, error) {
	result, err := userModel.Db.Exec("insert into users(login, hashed_pass, name) values(?, ?, ?)", login, pass, name)
	if err != nil {
		return 0, err
	}
	lastId, _ := result.LastInsertId()
	return lastId, nil
}
