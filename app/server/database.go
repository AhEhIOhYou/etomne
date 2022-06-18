package server

import (
	"database/sql"
	"log"
)

func Connect() *sql.DB {

	var db *sql.DB

	cfg := GetConfig()

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	return db
}
