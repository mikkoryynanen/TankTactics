package app

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

	newRoom := NewRoom()
	newRoom.Id = newRoomId
	go newRoom.Run()

	h.rooms[newRoomId] = newRoom

	fmt.Printf("Added room %v\n", newRoom.Id)

	go newRoom.AddConnectionAndRun(client)
}

func (h *Hub) ConnectToRoom(roomId uuid.UUID, client *websocket.Conn) bool {
	if room, exist := h.rooms[roomId]; exist {
		room.mu.Lock()
		// room.clients = append(room.clients, client)
		h.rooms[roomId] = room
		room.mu.Unlock()

		go room.AddConnectionAndRun(client)

		return true
	}
	return false
}
