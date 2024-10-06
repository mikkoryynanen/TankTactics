package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Room struct {
	Id        uuid.UUID
	IsRunning bool
	clients   []*websocket.Conn
	mu        sync.Mutex
}

func (r *Room) Run() {
	r.IsRunning = true

	go r.readMessages()

	// Endless loop
	for {
		r.mu.Lock()
		fmt.Printf("Room %v has clients %v \n", r.Id, r.clients)
		r.mu.Unlock()
		// for _, client := range r.clients {
		// 	fmt.Printf("client %v", client)
		// 	// fmt.Printf("Room %v has clients %v \n", r.Id, r.clients)
		// }

		// TODO We could be doing something here with the messages that we've read from another thread

		// r.readMessages()

		// Tick rate of the room
		time.Sleep(1 * time.Second)

	}
}

func (r *Room) readMessages() {
	for _, client := range r.clients {
		// TODO This should be closed from somewhere else?
		// defer client.Close()

		// Waiting for message, blocking
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

		fmt.Println(message)

		// 	err = c.WriteMessage(websocket.text)
		// 	if err != nil {
		// 		log.Fatal("Write message failed")
		// 	}
	}
}
