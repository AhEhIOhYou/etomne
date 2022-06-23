package server

import "github.com/go-sql-driver/mysql"

func GetConfig() mysql.Config {
	return mysql.Config{
		User:   "root",
		Passwd: "we92kjkszp",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "models_db",
	}
}
