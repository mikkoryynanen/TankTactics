package main

import (
	"encoding/json"
	"fmt"
	"main/logic"
	messageTypes "main/types"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Room struct {
	Id        uuid.UUID
	IsRunning bool
	clients   []*websocket.Conn
	mu        sync.Mutex
	world     logic.World
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

		receivedPosition := messageTypes.Position{}
		err = json.Unmarshal(msg, &receivedPosition)
		if err != nil {
			fmt.Println("Failed to unmarshal json from message")
		}

		// TODO Validate the payload

		fmt.Printf("Received message. (roomId/addr:message) %v/%v: %v\n", client.NetConn().LocalAddr().String(), r.Id, receivedPosition)

		// TODO
		// - Handle the message logic, what do we do when we get a message package
		r.world.TryAddClientValue(receivedPosition)
		// - Send back the computed response to that message

		err = client.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("Write message failed")
		}
	}
}
