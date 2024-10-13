package world

import (
	"main/handlers"
	client "main/types/Client"
)

/*
World is where where all of the actual client logic is contained, such as Moving or ChatMessages
*/
type World struct {
	Clients map[string]*client.Client

	handlers []handlers.Handler
}

func NewWorld() *World {
	// Define the different handlers
	var inputHandler handlers.InputHandler

	return &World{
		handlers: []handlers.Handler{0: &inputHandler},
	}
}

// Run the world simulation loop once. Should be called from room loop
func (w *World) RunOnce() {
	for _, client := range w.Clients {
		/* TODO
			Collision checks
			Valid move
		*/
		client.Position.PosX += int32(client.Input.InputX)
		client.Position.PosY += int32(client.Input.InputY)
	}
}

func (w *World) AddMessage(messageType int32, clientId string, block []byte) {
	if client, exists := w.Clients[clientId]; exists {
		w.handlers[messageType].Handle(block, client)
	}
}
