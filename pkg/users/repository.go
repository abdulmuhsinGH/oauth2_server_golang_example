package users

import (
	"log"
	"os"

	"github.com/go-pg/pg/v9"
)

/*
Repository is an interface for anyone using this
*/
type Repository interface {
	Crud
	FindUserByUsername(string) (*User, error)
	FindOrAddUser(*User) (*User, error)
}

type repository struct {
	db *pg.DB
}

var (
	logger *log.Logger
)

/*
NewRepository creates a users repository with the necessary dependencies
*/
func NewRepository(db *pg.DB) Repository {
	logger = log.New(os.Stdout, "user_repository \n", log.LstdFlags|log.Lshortfile)
	return repository{db}

}

/*
	AddUser creates a new user
*/
func (r repository) AddUser(user *User) error {
	return r.db.Insert(user)
}

/*
FindOrAddUser finds user or saves user if not found to the user's table
*/
func (r repository) FindOrAddUser(user *User) (*User, error) {
	_, err := r.db.Model(user).
		Column("id").
		Where("email = ?email").
		OnConflict("DO NOTHING"). // OnConflict is optional
		Returning("id").
		SelectOrInsert()
	if err != nil {
		logger.Println("FindORAddUser_Error", err.Error())
		return &User{}, err
	}

	return user, nil

}

/*
GetAllUsers returns all users from the user's table
*/
func (r repository) GetAllUsers() ([]User, error) {
	users := []User{}
	err := r.db.Model(&users).Select()
	if err != nil {
		logger.Println("GetAllusers_Repo_Error", err.Error())
		return nil, err
	}
	return users, nil
}

/*
FindUserByUsername return a user based on the username
*/
func (r repository) FindUserByUsername(username string) (*User, error) {
	user := new(User)

	err := r.db.Model(user).Where("username = ?", username).Select()
	if err != nil {
		logger.Println("FindUserByUsername_Error", err.Error())
		return &User{}, err
	}
	return user, nil
}
