package world

import (
	"fmt"
	"main/handlers"
	client "main/types/Client"
	"time"
)

/*
World is where where all of the actual client logic is contained, such as Moving or ChatMessages
*/
type World struct {
	// TODO Why is there clients here?
	Clients map[string]*client.Client

	handlers []handlers.Handler

	// Time
	lastFrameTime time.Time
}

func NewWorld() *World {
	// Define the different handlers
	var inputHandler handlers.InputHandler

	return &World{
		Clients: make(map[string]*client.Client),
		handlers: []handlers.Handler{0: &inputHandler},
	}
}

// Run the world simulation loop once. Should be called from room loop
func (w *World) RunOnce() {
	w.lastFrameTime = time.Now()

	for _, client := range w.Clients {
		/* TODO
		Collision checks
		Valid move
		*/

		currentTime := time.Now()
		deltaTime := currentTime.Sub(w.lastFrameTime).Seconds()

		posX := client.Position.PosX
		posX += float32(client.Input.InputX) * 200000 * float32(deltaTime)

		posY := client.Position.PosY
		posY += float32(client.Input.InputY) * 200000 * float32(deltaTime)

		client.Position.PosX = posX
		client.Position.PosY = posY

		fmt.Printf("client status %v\n", client)
	}
}

func (w *World) AddMessage(messageType int32, clientId string, block []byte) {
	if client, exists := w.Clients[clientId]; exists {
		w.handlers[messageType].Handle(block, client)
	}
}
