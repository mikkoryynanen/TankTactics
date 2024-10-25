package app

import (
	"fmt"
	"main/cmd/database"
	"main/cmd/events"
	"main/cmd/network"
	"main/cmd/routes"
)

type App struct {
	hub         Hub
	network     network.Network
	database    *database.Database
	userHandler *routes.UserHandler

	eventEmitter *events.EventEmitter
}

func NewApp() *App {
	db := database.NewDatabase()
	eventEmitter := events.NewEventEmiter()
	return &App{
		hub:         *NewHub(),
		network:     network.NewWebsocket(*eventEmitter),
		database:    db,
		userHandler: routes.NewUserHandler(db),

		eventEmitter: eventEmitter,
	}
}

func (a *App) Run() {
	a.eventEmitter.On(events.EventOnClientConnected, a.onClientConnected)
	a.eventEmitter.On(events.EventOnClientJoinRoom, a.onClientJoinRoom)

	a.network.Connect()
}

func (a *App) onClientConnected(event events.Event) {
	// TODO Crashes here
	connectionEvent := event.(events.ConnectionEvent)
	a.hub.AddRoom(connectionEvent.Conn)
}

func (a *App) onClientJoinRoom(event events.Event) {
	roomConnectionEvent := event.(events.RoomConnectionEvent)
	isConnected := a.hub.ConnectToRoom(roomConnectionEvent.RoomId, roomConnectionEvent.Conn)
	if !isConnected {
		roomConnectionEvent.Conn.Close()
	}

	fmt.Printf("isConnected? %v to room %v\n", isConnected, roomConnectionEvent.RoomId)
}
