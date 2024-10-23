package handlers

import (
	"encoding/json"
	"fmt"
	"main/cmd/types"
)

type Handler interface {
	Handle(block []byte, client *types.Client)
}

type InputHandler struct {
}

func (ph *InputHandler) Handle(block []byte, client *types.Client) {
	// Since this is real-time, we only receive inputs from the player
	// Once correct position is calculated with the given input, we send that the the client

	var receivedInput types.InputPayload
	err := json.Unmarshal(block, &receivedInput)
	if err != nil {
		fmt.Println("Failed to unmarshal InputPayload from message")
	}

	fmt.Println(receivedInput)

	// TODO
	// Validate the move
	// Check if the move can be made
	// Send the response to client
	// If we position cannot be moved to, just send the current position or last known position?

	// Normalize
	if receivedInput.InputX > 1 {
		receivedInput.InputX = 1
	}
	if receivedInput.InputX < -1 {
		receivedInput.InputX = -1
	}
	if receivedInput.InputY > 1 {
		receivedInput.InputY = 1
	}
	if receivedInput.InputY < -1 {
		receivedInput.InputY = -1
	}

	client.Input.InputX = receivedInput.InputX
	client.Input.InputY = receivedInput.InputY
}
