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
	// newRoomId := uuid.New()
	newRoomId, _ := uuid.Parse("c494027b-9d5a-4785-8b3b-89187f72c44a")

	var newRoom Room
	newRoom.Id = newRoomId
	newRoom.clients = append(newRoom.clients, client)

	h.rooms[newRoomId] = &newRoom

	fmt.Printf("Added room %v\n", newRoom.Id)

	go newRoom.Run()
}

func (h *Hub) ConnectToRoom(roomId uuid.UUID, client *websocket.Conn) bool {
	if room, exist := h.rooms[roomId]; exist {
		room.mu.Lock()
		room.clients = append(room.clients, client)
		h.rooms[roomId] = room
		room.mu.Unlock()

		return true
	}
	return false
}
