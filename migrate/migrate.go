package migrate

import (
	"database/sql"
	"fmt"
	"log"

	"kroflic/gin-library/db"
	"kroflic/gin-library/models"
)

func Run() {

	// Connect

	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Check, create, insert

	createUserTable(db)
	createBookTable(db)
	createRentalTable(db)

	var books = []models.NewBook{
		{Title: "Ameriški bogovi", Author: "Neil Gaiman", Quantity: 3},
		{Title: "Deseti december", Author: "George Saunders", Quantity: 2},
		{Title: "Klavnica 5", Author: "Kurt Vonnegut", Quantity: 1},
	}

	insertBook(db, books[0])
	insertBook(db, books[1])
	insertBook(db, books[2])

	defer db.Close()
}

func createUserTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS "users" (
		id SERIAL PRIMARY KEY,
		firstName VARCHAR(100) NOT NULL,
		lastName VARCHAR(100) NOT NULL
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func createBookTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS "books" (
		id SERIAL PRIMARY KEY,
		title VARCHAR(100) NOT NULL,
		author VARCHAR(100) NOT NULL,
		quantity INTEGER NOT NULL,
		UNIQUE (title, author)
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func createRentalTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS "rentals" (
		id SERIAL PRIMARY KEY,
		userId INTEGER NOT NULL,
		bookId INTEGER NOT NULL,
		returned BOOLEAN NOT NULL,
		FOREIGN KEY (userId) REFERENCES users(id),
		FOREIGN KEY (bookId) REFERENCES books(id)
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func insertBook(db *sql.DB, newBook models.NewBook) (int, error) {
	query := `
		INSERT INTO books (title, author, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (title, author) DO NOTHING
		RETURNING id
	`

	var newPrimaryKey int
	err := db.QueryRow(query, newBook.Title, newBook.Author, newBook.Quantity).Scan(&newPrimaryKey)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("book already exists")
		}
		return 0, err
	}

	return newPrimaryKey, nil
}
