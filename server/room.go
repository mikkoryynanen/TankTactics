package main

import (
	client "main/types/Client"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Room struct {
	Id        uuid.UUID
	IsRunning bool
	clients   []*client.Client
	mu        sync.Mutex
}

func NewRoom() *Room {
	return &Room{}
}

func (r *Room) AddConnection(conn *websocket.Conn) {
	r.addClient(conn)
}

func (r *Room) AddConnectionAndRun(conn *websocket.Conn) {
	newClient := r.addClient(conn)

	go newClient.ReadMessages()
}

// Run the world loop with TickRate
// This method loop through the read messages from the clients
func (r *Room) Run(tickRate int32) {
}

// func (r *Room) TryAddClientValue(clientId uuid.UUID, position messageTypes.Position) {
// 	for _, client := range r.clients {
// 		// TODO Make sure we can actually make this move
// 		// TODO We could use the rollback technique here
//
// 		// FOR NEXT TIME WHEN LOOKING AT THIS
// 		// Since this is TurnBased game, we have to request out wanted postion that we receive from the server
// 		// Once that position has been confirmed to bee correct, we send it to the client
//
// 		client.Position = messageTypes.Position{position.PosX, position.PosY}
// 	}
//
// }
//
func (r *Room) addClient(conn *websocket.Conn) *client.Client {
	r.mu.Lock()
	newClient := client.NewClient()
	r.clients = append(r.clients, newClient)
	r.mu.Unlock()

	return newClient
}
