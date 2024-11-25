package main

import (
	"kroflic/gin-library/controllers/bookController"
	"kroflic/gin-library/controllers/rentalController"
	"kroflic/gin-library/controllers/userController"
	docs "kroflic/gin-library/docs"
	"kroflic/gin-library/migrate"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	router := gin.Default()

	migrate.Run()

	docs.SwaggerInfo.BasePath = "/api/v1"

	v1 := router.Group("/api/v1")
	{
		v1.GET("/users", userController.GetUsers)
		v1.POST("/users/create", userController.CreateUser)

		v1.GET("/books/all", bookController.GetBooks)
		v1.GET("/books/available", bookController.GetAvailableBooks)

		v1.GET("/rentals", rentalController.GetRentals)
		v1.POST("/rentals/create", rentalController.CreateRental)
		v1.POST("/rentals/return", rentalController.ReturnRental)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run("0.0.0.0:8080")
}
