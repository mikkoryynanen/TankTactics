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

type Room struct {
	Id        uuid.UUID
	IsRunning bool
	clients   []*client.Client
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

func (r *Room) AddConnection(conn *websocket.Conn) {
	r.addClient(conn)
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
		// Cleanup disconnected clients
		utils.RemoveDisconnectedClients(utils.GetMapValues(r.world.Clients))

		r.world.RunOnce()

		// Send the world state back to clients
		for _, client := range r.clients {
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

		// fmt.Println("Room tick")
		time.Sleep(time.Duration(150) * time.Millisecond)
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
	r.clients = append(r.clients, newClient)
	// TODO Rewrites it every time
	r.world.Clients = make(map[string]*client.Client)
	r.world.Clients["client"] = newClient
	r.mu.Unlock()

	return newClient
}
