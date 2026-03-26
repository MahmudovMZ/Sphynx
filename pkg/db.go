package db

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Подключение к БД
var db *sql.DB

func ConnectDB(username, password, dbname, address string) (err error) {
	dns := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, address, dbname)

	db, err = sql.Open("pgx", dns)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}

func CloseDB() {
	db.Close()
}

func GetDB() *sql.DB {
	return db
}
