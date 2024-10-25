package events

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	EventOnClientConnected = "OnClientConnected"
	EventOnClientJoinRoom = "OnClientJoinRoom"
)

type ConnectionEvent struct {
	Conn *websocket.Conn
}

type RoomConnectionEvent struct {
	RoomId uuid.UUID
	Conn *websocket.Conn
}