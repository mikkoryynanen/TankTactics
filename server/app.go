package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type App struct {
	hub      Hub
	upgrader websocket.Upgrader

	database *Database

	userHandler *UserHandler
}

func NewApp() *App {
	db := NewDatabase()
	return &App{
		hub: *NewHub(),
		database: db,
		userHandler: NewUserHandler(db),
	}
}

func (a *App) handleConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connection received")

	c, err := a.upgradeConnection(w, r)
	if err != nil {
		fmt.Println(err)
	}

	a.hub.AddRoom(c)
}

func (a *App) handleRoomConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connection received. Connecting to room...")

	queryParams := r.URL.Query()
	roomId, err := uuid.Parse(queryParams.Get("roomId"))
	if err != nil {
		fmt.Println("failed to parse roomId")
		return
	}
	c, err := a.upgradeConnection(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
	isConnected := a.hub.ConnectToRoom(roomId, c)
	if !isConnected {
		c.Close()
	}

	fmt.Printf("isConnected? %v to room %v\n", isConnected, roomId)
}

func (a App) upgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	c, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Could not upgrade HTTP connection to Websocket")
		return nil, err
	}
	return c, nil
}

func (a App) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/c", a.handleConnection)
	r.HandleFunc("/c/room", a.handleRoomConnection)

	// TODO Disabled for now,continuing later
	// r.HandleFunc("/user", a.userHandler.CreateUserHandler).Methods("POST")
	// r.HandleFunc("/user/{id}", a.userHandler.GetUserHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
