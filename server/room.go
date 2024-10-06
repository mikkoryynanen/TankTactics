package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Room struct {
	Id        uuid.UUID
	IsRunning bool
	clients   []*websocket.Conn
	mu        sync.Mutex
}

func (r *Room) AddClientAndRun(client *websocket.Conn) {
	r.mu.Lock()
	r.clients = append(r.clients, client)
	r.mu.Unlock()

	go r.readClientMessages(client)
}

func (r *Room) readClientMessages(client *websocket.Conn) {
	defer client.Close()

	for {
		_, msg, err := client.ReadMessage()
		if err != nil {
			fmt.Printf("Failed to read message. err: %v\n", err)
			return
		}

		message := Message{}
		err = json.Unmarshal(msg, &message)
		if err != nil {
			fmt.Println("Failed to unmarshal json from message")
		}

		fmt.Printf("Received message. (roomId/addr:message) %v/%v: %v\n", client.NetConn().LocalAddr().String(), r.Id, message)

		bytes, _ := json.Marshal(message)
		err = client.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			fmt.Println("Write message failed")
		}
	}
}
