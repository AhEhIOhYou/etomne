package server

import "github.com/go-sql-driver/mysql"

func GetConfig() mysql.Config {
	return mysql.Config{
		User:   "admin",
		Passwd: "Password10$",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "models_db",
	}
}
