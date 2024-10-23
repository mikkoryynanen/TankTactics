package routes

import (
	"io"
	"log"
	"main/cmd/database"
	"main/cmd/types"
	"main/cmd/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CreateNewUserPayload struct {
	Username string
}

func UserRouter() *mux.Router {
	// Inspiration from https://stackoverflow.com/a/52473957
	router := mux.NewRouter()
	// TODO No need for instance, can just use funcs?
	uh := NewUserHandler(nil)

	router.HandleFunc("/user", uh.CreateUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/user/{id}", uh.GetUserHandler).Methods(http.MethodGet)

	return router
}

func NewUser(userId uuid.UUID, username string) *types.User {
	return &types.User{
		Id:       userId,
		Username: username,
	}
}

type UserHandler struct {
	database *database.Database
}

func NewUserHandler(db *database.Database) *UserHandler {
	return &UserHandler{
		database: db,
	}
}

func (uh *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
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

	// Check that the user does not already exist
	userId, err := uh.database.CreateUser(createUserBody.Username)
	if err != nil {
		log.Fatal(err)
	}

	user := NewUser(userId, createUserBody.Username)

	bytes, _ := utils.GetBytes(user)
	w.Write(bytes)
}

func (uh *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting user with id")

	// body, err := io.ReadAll(r.Body)
	// if len(body) <= 0 || err != nil {
	// 	log.Print(err)
	// 	http.Error(w, "Could not read request body", http.StatusInternalServerError)
	// 	return
	// }

	// getUserBody, err := utils.GetType[GetUserPayload](body)
	// if err != nil {
	// 	log.Print(err)
	// 	http.Error(w, "Could not get type from reqest body", http.StatusInternalServerError)
	// 	return
	// }

	// requestUserId := r.URL.Query().Get("userId")
	// user, err := uh.database.GetUser(requestUserId)
	// if err != nil {
	// 	w.Write([]byte("Could not find user"))
	// }

	// w.Write(user.Id[:])
}
