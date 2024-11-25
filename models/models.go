package models

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type NewUser struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Book struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

type NewBook struct {
	Title    string `json:"title" binding:"required"`
	Author   string `json:"author" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
}

type Rental struct {
	ID       int  `json:"id"`
	BookId   int  `json:"bookId"`
	UserId   int  `json:"userId"`
	Returned bool `json:"returned"`
}

type NewRental struct {
	BookId int `json:"bookId" binding:"required"`
	UserId int `json:"userId" binding:"required"`
}

type ReturningRental struct {
	BookId int `json:"bookId"`
	UserId int `json:"userId"`
}
