package app

import (
	"fmt"
	"main/cmd/handlers"
	"main/cmd/types"
	"time"
)

/*
World is where where all of the actual client logic is contained, such as Moving or ChatMessages
*/
type World struct {
	Clients  map[string]*types.Client
	handlers []handlers.Handler
	level    *Level

	// Time
	lastFrameTime time.Time
}

func NewWorld() *World {
	// Define the different handlers
	var inputHandler handlers.InputHandler

	level := NewLevel()

	return &World{
		Clients: make(map[string]*types.Client),
		handlers: []handlers.Handler{
			0: &inputHandler,
		},
		level: level,
	}
}

// Run the world simulation loop once. Should be called from room loop
func (w *World) SimulateOnce() {
	w.lastFrameTime = time.Now()

	for _, client := range w.Clients {
		currentTime := time.Now()
		deltaTime := currentTime.Sub(w.lastFrameTime).Seconds()

		posX := client.Object.Position.X
		posX += float32(client.Input.InputX) * 200000 * float32(deltaTime)

		posY := client.Object.Position.Y
		posY += float32(client.Input.InputY) * 200000 * float32(deltaTime)

		// Collision check
		for _, object := range w.level.Objects {
			// TODO Do not check the whole level, check collisions at a radious
			if w.level.IsObjectColliding(client.Object, object) {
				fmt.Printf("%v collided wih %v\n", client.Object.Name, object.Name)
				// If we collide do not move to that position, move it to the edge
				posX = object.Position.X 
				posY = object.Position.Y
			}
		}

		client.Object.Position.X = posX
		client.Object.Position.Y = posY

		fmt.Printf("client status %v\n", client)
	}
}

func (w *World) AddMessage(messageType int32, clientId string, block []byte) {
	if client, exists := w.Clients[clientId]; exists {
		w.handlers[messageType].Handle(block, client)
	} else {
		fmt.Printf("Did not find client with %v\n", clientId)
	}
}
