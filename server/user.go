package main

import (
	"io"
	"log"
	"main/utils"
	"net/http"

	"github.com/google/uuid"
)

type CreateNewUserPayload struct {
	Username string
}

type User struct {
	Id 			string
	Username 	string
	CurrentRoom	*Room
}

func NewUser(username string) *User{
	return &User{
		Id: uuid.NewString(),
		Username: username,
	}
}

type UserHandler struct {
	database *Database
}

func NewUserHandler(db *Database) *UserHandler {
	return &UserHandler{
		database: db,
	}
}

func (uh * UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating user...")

	body, err := io.ReadAll(r.Body)
	if len(body) <= 0 || err != nil {
		log.Print(err)
		http.Error(w, "Could not read request body", http.StatusInternalServerError)
		return
	}

	createUserBody, err := utils.GetType[CreateNewUserPayload](body)
	if err != nil {
		log.Print(err)
		http.Error(w, "Could not get type from reqest body", http.StatusInternalServerError)
		return
	}

	// Perform validation
	if createUserBody.Username == "" {
		http.Error(w, "username missing from request body", http.StatusBadRequest)
		return
	}

	user := NewUser(createUserBody.Username)
	// TODO Save user to database
	uh.database.CreateUser()

	bytes, _ := utils.GetBytes(user)
	w.Write(bytes)
}