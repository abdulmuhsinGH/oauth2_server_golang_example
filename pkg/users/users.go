package users

import (
	uuid "github.com/satori/go.uuid"
)

/*
Crud describes all methods used in the service and repository
*/
type Crud interface {
	AddUser(*User) error
	GetAllUsers() ([]User, error)
}

/*
User defines the properties of a user
*/
type User struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Firstname  string    `json:"firstname"`
	Middlename string    `json:"middlename"`
	Lastname   string    `json:"lastname"`
	Email      string    `json:"email"`
	Gender     string    `json:"gender"`
}
