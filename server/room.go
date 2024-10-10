package main

import (
	"main/logic"
	messageTypes "main/types"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Room struct {
	Id        uuid.UUID
	IsRunning bool
	clients   []*messageTypes.Client
	mu        sync.Mutex
	world     logic.World
}

func NewRoom() *Room {
	return &Room{
		world: logic.World{},
	}
}

func (r *Room) AddConnection(conn *websocket.Conn) {
	r.addClient(conn)
}

func (r *Room) AddConnectionAndRun(conn *websocket.Conn) {
	newClient := r.addClient(conn)

	go newClient.ReadMessages()
}

func (r *Room) addClient(conn *websocket.Conn) *messageTypes.Client {
	r.mu.Lock()
	newClient := &messageTypes.Client{*conn}
	r.clients = append(r.clients, newClient)
	r.mu.Unlock()

	return newClient
}
