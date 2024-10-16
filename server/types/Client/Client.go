package client

import (
	"log"
	messageTypes "main/types/payloads"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Id 			string
	Conn        *websocket.Conn
	Position    messageTypes.PositionPayload
	Input       messageTypes.InputPayload
	IsConnected bool
}

func NewClient(conn *websocket.Conn) *Client {
	newClient := &Client{
		Id: uuid.NewString(),
		Conn:        conn,
		IsConnected: true,
	}
	conn.SetCloseHandler(func(code int, text string) error {
		newClient.IsConnected = false
		return nil
	})

	return newClient
}

// To be called once as goroutine
func (c *Client) ReadMessages(stream chan []byte) {

	for {
		defer c.Conn.Close()

		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("Normal closure:", err)
				c.IsConnected = false
			}
			return
		}

		stream <- msg
	}
}
