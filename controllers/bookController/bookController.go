package bookController

import (
	"fmt"
	"log"

	"kroflic/gin-library/db"
	"kroflic/gin-library/models"
	"kroflic/gin-library/response"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {

	// Connect

	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
		response.SendInternalServerError(c)
		return
	}

	// Get Books

	books := []models.Book{}

	rows, err := db.Query("SELECT id, title, author, quantity FROM books;")
	if err != nil {
		log.Fatal(err)
		response.SendInternalServerError(c)
		return
	}

	defer rows.Close()

	var id int
	var title string
	var author string
	var quantity int

	for rows.Next() {
		err := rows.Scan(&id, &title, &author, &quantity)
		if err != nil {
			log.Fatal(err)
			response.SendInternalServerError(c)
			return
		}
		books = append(books, models.Book{ID: id, Title: title, Author: author, Quantity: quantity})
	}

	fmt.Println(books)

	response.SendOK(c, books)
}

//
// GET AVAILABLE BOOKS
//

func GetAvailableBooks(c *gin.Context) {

	// Connect

	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
		response.SendInternalServerError(c)
		return
	}

	// Get Books

	books := []models.Book{}

	rows, err := db.Query(`
		SELECT b.id, b.title, b.author, b.quantity
		FROM books b
		LEFT JOIN rentals r ON r.bookId = b.id AND r.returned = false
		GROUP BY b.id
		HAVING COUNT(r.id) < b.quantity;
	`)
	if err != nil {
		log.Fatal(err)
		response.SendInternalServerError(c)
		return
	}

	defer rows.Close()

	var id int
	var title string
	var author string
	var quantity int

	for rows.Next() {
		err := rows.Scan(&id, &title, &author, &quantity)
		if err != nil {
			log.Fatal(err)
			response.SendInternalServerError(c)
			return
		}
		books = append(books, models.Book{ID: id, Title: title, Author: author, Quantity: quantity})
	}

	fmt.Println(books)

	response.SendOK(c, books)
}
