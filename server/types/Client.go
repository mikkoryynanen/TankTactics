package messageTypes

// TODO Maybe just types?

import (
	"encoding/json"
	"fmt"
	messageTypes "main/types"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Position messageTypes.Position
}

// To be called once as goroutine
func (c *Client) ReadMessages() {
	defer c.Conn.Close()

	for {
		_, msg, err := c.Conn.ReadMessage()
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

		// fmt.Printf("Received message. (roomId/addr:message) %v/%v: %v\n", c.Conn.NetConn().LocalAddr().String(), r.Id, receivedPosition)

		// TODO
		// - Handle the message logic, what do we do when we get a message package
		// r.world.TryAddClientValue(clien.id, receivedPosition)
		// - Send back the computed response to that message

		// TODO Response is written back once we've calculated the actual position inside world
		// 	err = client.WriteMessage(websocket.TextMessage, msg)
		// 	i
		// f err != nil {
		// 		fmt.Println("Write message failed")
		// 	}
	}
}