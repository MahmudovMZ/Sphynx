package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Подключение к БД
var DB *sql.DB

func DataBase() {
	dsn := "postgres://postgres:newpassword@localhost:5432/wordsdb?sslmode=disable"
	var err error

	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
}
