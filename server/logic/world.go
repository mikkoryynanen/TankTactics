package logic

import (
	messageTypes "main/types"

	"github.com/google/uuid"
)

type World struct {
	clients map[uuid.UUID]*messageTypes.Client
}

func (w World) Run() {
}

func (w World) TryAddClientValue(clientId uuid.UUID, position messageTypes.Position) {
	if client, exists := w.clients[clientId]; exists {
		// TODO Make sure we can actually make this move
		// TODO We could use the rollback technique here

		// FOR NEXT TIME WHEN LOOKING AT THIS
		// Since this is TurnBased game, we have to request out wanted postion that we receive from the server
		// Once that position has been confirmed to bee correct, we send it to the client

		client.Position = messageTypes.Position{position.PosX, position.PosY}
	}

}
