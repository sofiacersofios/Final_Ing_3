package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Db es una variable global que representa la conexión a la base de datos
var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("mysql", "root:root@tcp(172.23.0.2:3306)/mydatabase")
	if err != nil {
		log.Fatal(err)
	}
	if err := Db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to the database")
}
