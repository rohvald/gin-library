package userController

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
// GET USERS
//

func GetUsers(c *gin.Context) {

	// Connect

	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
		response.SendInternalServerError(c)
		return
	}

	// Get Users

	users := []models.User{}

	rows, err := db.Query("SELECT id, firstName, lastName FROM users")
	if err != nil {
		log.Fatal(err)
		response.SendInternalServerError(c)
		return
	}

	defer rows.Close()

	var id int
	var firstName string
	var lastName string

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, models.User{ID: id, FirstName: firstName, LastName: lastName})
	}

	fmt.Println(users)

	response.SendOK(c, users)
}

//
// CREATE USER
//

func CreateUser(c *gin.Context) {

	// Connect

	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
		response.SendInternalServerError(c)
		return
	}

	// Create User

	var newUser models.NewUser
	if err := c.BindJSON(&newUser); err != nil {
		response.SendBadRequestError(c, "Invalid request body")
		return
	}

	// Insert new User

	newPrimaryKey := insertUser(db, newUser)

	fmt.Printf("ID = %d", newPrimaryKey)

	response.SendCreated(c, gin.H{"success": true, "id": newPrimaryKey})
}

//
// INSERT USER (helper)
//

func insertUser(db *sql.DB, newUser models.NewUser) int {
	query := `INSERT INTO users (firstName, lastName) VALUES ($1, $2) RETURNING id`

	var newPrimaryKey int

	err := db.QueryRow(query, newUser.FirstName, newUser.LastName).Scan(&newPrimaryKey)

	if err != nil {
		log.Fatal(err)
	}

	return newPrimaryKey
}
