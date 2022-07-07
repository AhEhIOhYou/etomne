package server

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

func Connect() *sql.DB {

	var db *sql.DB

	yfile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]string)
	yaml.Unmarshal(yfile, &data)

	cfg := mysql.Config{
		User:   data["User"],
		Passwd: data["Passwd"],
		Net:    data["Net"],
		Addr:   data["Addr"],
		DBName: data["DBName"],
	}

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
