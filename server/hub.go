package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Hub struct {
	rooms map[uuid.UUID]*Room
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[uuid.UUID]*Room),
	}
}

func (h *Hub) AddRoom(client *websocket.Conn) {
	newRoomId := uuid.New()

	var newRoom Room
	newRoom.Id = newRoomId

	h.rooms[newRoomId] = &newRoom

	fmt.Printf("Added room %v\n", newRoom.Id)

	go newRoom.AddClientAndRun(client)
}

func (h *Hub) ConnectToRoom(roomId uuid.UUID, client *websocket.Conn) bool {
	if room, exist := h.rooms[roomId]; exist {
		room.mu.Lock()
		room.clients = append(room.clients, client)
		h.rooms[roomId] = room
		room.mu.Unlock()

		go room.AddClientAndRun(client)

		return true
	}
	return false
}
