package main

import (
	"encoding/json"
	"fmt"
	client "main/types/Client"
	world "main/types/World"
	messageTypes "main/types/payloads"
	"main/utils"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	TickRate     = 60                                              // Target tick rate in Hz
	TickDuration = time.Duration(1000/TickRate) * time.Millisecond // ~16.67 ms per tick
)

type Room struct {
	Id        uuid.UUID
	IsRunning bool
	mu        sync.Mutex

	// Replicate of the game word. Everything related to clients is contained in here
	world *world.World
	// Stream of messages
	stream chan []byte
}

func NewRoom() *Room {
	return &Room{
		// Create buffered channel, the channel size is still a question
		stream: make(chan []byte, 10),
		world:  world.NewWorld(),
	}
}

func (r *Room) AddConnectionAndRun(conn *websocket.Conn) {
	newClient := r.addClient(conn)

	go newClient.ReadMessages(r.stream)
}

// Run the world loop with TickRate
// This method loop through the read messages from the clients
func (r *Room) Run() {
	fmt.Println("starting room")

	go r.receiveMessages()

	for {
		startTime := time.Now()

		// Cleanup disconnected clients
		// TODO This could be done with an interval if deemed to take too long
		r.world.Clients = utils.RemoveDisconnectedClients(utils.GetMapValues(r.world.Clients))

		r.world.SimulateOnce()

		// TODO Consider moving this to world
		// Send the world state back to clients
		for _, client := range r.world.Clients {
			serverState := &messageTypes.ServerState{
				PosX: client.Position.PosX,
				PosY: client.Position.PosY,
			}
			data, err := json.Marshal(serverState)
			if err != nil {
				fmt.Println("Failed to marshal data")
			}
			client.Conn.WriteMessage(websocket.TextMessage, data)
		}

		elapsedTime := time.Since(startTime)
		sleepDuration := TickDuration - elapsedTime
		fmt.Printf("sleeping for %v\n", sleepDuration)
		if sleepDuration > 0 {
			time.Sleep(time.Duration(sleepDuration))
		} else {
			fmt.Println("skipping sleep, frame took too long to render")
		}
	}
}

func (r *Room) receiveMessages() {
	for {
		for block := range r.stream {
			var baseMsg messageTypes.BaseMessage
			err := json.Unmarshal(block, &baseMsg)
			if err != nil {
				fmt.Println("Failed to unmarshal json from message")
			}

			fmt.Printf("received message: %v\n", baseMsg)

			r.world.AddMessage(baseMsg.Type, baseMsg.ClientId, block)
		}
	}
}

func (r *Room) addClient(conn *websocket.Conn) *client.Client {
	r.mu.Lock()
	newClient := client.NewClient(conn)
	r.world.Clients[newClient.Id] = newClient
	r.mu.Unlock()

	// TODO These write messages could be moved to somewhere unified so we do not call them willy-nilly
	// Send player their generated metadata
	playerMetadata := &messageTypes.PlayerMetadata{
		ClientId: newClient.Id,
	}
	data, err := json.Marshal(playerMetadata)
	if err != nil {
		fmt.Println("Failed to marshal data")
	}
	conn.WriteMessage(websocket.TextMessage, data)

	return newClient
}
