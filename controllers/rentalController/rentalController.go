package rentalController

import (
	"database/sql"
	"fmt"
	"log"

	"kroflic/gin-library/db"
	"kroflic/gin-library/models"
	"kroflic/gin-library/response"

	"github.com/gin-gonic/gin"
)

//
// GET RENTALS
//

func GetRentals(c *gin.Context) {

	// Connect

	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
		response.SendInternalServerError(c)
		return
	}

	// Get Rentals

	rentals := []models.Rental{}

	rows, err := db.Query("SELECT id, bookId, userId, returned FROM rentals;")
	if err != nil {
		log.Fatal(err)
		response.SendInternalServerError(c)
		return
	}

	defer rows.Close()

	var id int
	var bookId int
	var userId int
	var returned bool

	for rows.Next() {
		err := rows.Scan(&id, &bookId, &userId, &returned)
		if err != nil {
			log.Fatal(err)
		}
		rentals = append(rentals, models.Rental{ID: id, BookId: bookId, UserId: userId, Returned: returned})
	}

	fmt.Println(rentals)

	// Respond

	response.SendOK(c, rentals)
}

//
// CREATE RENTAL
//

func CreateRental(c *gin.Context) {

	// Connect

	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
		response.SendInternalServerError(c)
		return
	}

	// Create new Rental

	var newRental models.NewRental
	if err := c.ShouldBindJSON(&newRental); err != nil {
		response.SendBadRequestError(c, "Invalid request body")
		return
	}

	// Check if user exists

	userExists := checkUserExists(db, newRental.UserId)

	if !userExists {
		response.SendNotFoundError(c, "user not found")
		return
	}

	// Check if there's available books

	bookAvailable := checkBookAvailability(db, newRental.BookId)

	if !bookAvailable {
		response.SendBadRequestError(c, "all copies rented out")
		return
	}

	// Insert new Rental

	newPrimaryKey := insertRental(db, newRental)

	fmt.Printf("ID = %d", newPrimaryKey)

	// Respond

	response.SendCreated(c, newPrimaryKey)
}

//
// CHECK USER EXISTS (helper)
//

func checkUserExists(db *sql.DB, userId int) bool {

	var id int
	var firstName string
	var lastName string
	err := db.QueryRow("SELECT * FROM users WHERE id = $1", userId).Scan(&id, &firstName, &lastName)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User not found")
		} else {
			log.Fatal(err)
		}
		return false
	}

	return true
}

//
// CHECK BOOK AVAILABILITY (helper)
//

func checkBookAvailability(db *sql.DB, bookId int) bool {

	// Get book quantity

	var quantity int
	err := db.QueryRow("SELECT quantity FROM books WHERE id = $1", bookId).Scan(&quantity)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Book not found")
		} else {
			log.Fatal(err)
		}
		return false
	}

	// Get not returned rentals of this book

	var rentedOutCount int
	err = db.QueryRow("SELECT COUNT(*) FROM rentals WHERE bookId = $1 AND returned = false", bookId).Scan(&rentedOutCount)

	if err != nil {
		log.Fatal(err)
		return false
	}

	return rentedOutCount < quantity
}

//
// INSERT RENTAL (helper)
//

func insertRental(db *sql.DB, newRental models.NewRental) int {
	query := `INSERT INTO rentals (bookId, userId, returned) VALUES ($1, $2, $3) RETURNING id`

	var newPrimaryKey int

	err := db.QueryRow(query, newRental.BookId, newRental.UserId, false).Scan(&newPrimaryKey)

	if err != nil {
		log.Fatal(err)
	}

	return newPrimaryKey
}

//
// RETURN RENTAL
//

func ReturnRental(c *gin.Context) {

	// Connect

	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
		response.SendInternalServerError(c)
		return
	}

	// Get rentalId

	var returningRental models.ReturningRental
	if err := c.BindJSON(&returningRental); err != nil {
		response.SendBadRequestError(c, "Invalid request body")
		return
	}

	// Update rental

	updateSuccess, err := markRentalReturned(db, returningRental)

	if !updateSuccess {
		if err == sql.ErrNoRows {
			response.SendNotFoundError(c, "no rental of this book by this user found")
			return
		}
		response.SendInternalServerError(c)
		return
	}

	// Respond

	response.SendOK(c, returningRental)
}

//
// MARK RENTAL AS RETURNED (helper)
//

func markRentalReturned(db *sql.DB, returningRental models.ReturningRental) (bool, error) {
	query := `
		WITH cte AS (
			SELECT id
			FROM rentals
			WHERE bookId = $1 AND userId = $2 AND returned = FALSE
			LIMIT 1
    )
    UPDATE rentals
    SET returned = TRUE
    WHERE id = (SELECT id FROM cte)
    RETURNING id
	`

	var updatedRentalId int

	err := db.QueryRow(query, returningRental.BookId, returningRental.UserId).Scan(&updatedRentalId)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No rental of book %d by user %d found", returningRental.BookId, returningRental.UserId)
		} else {
			log.Fatal(err)
		}
		return false, err
	}

	return true, nil
}
