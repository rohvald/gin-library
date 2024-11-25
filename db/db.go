package db

import (
	"database/sql"
	"log"
)

const connStr = "postgres://postgres:secret@db:5432/ginlibrary?sslmode=disable"

func Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}

func QueryRows(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}
