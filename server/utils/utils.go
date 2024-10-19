package utils

import (
	"encoding/json"
	"fmt"
	client "main/types/Client"
)

func RemoveDisconnectedClients(clients []*client.Client) map[string]*client.Client {
	var filteredClients []*client.Client
	for _, client := range clients {
		if client.IsConnected {
			filteredClients = append(filteredClients, client)
		}
	}

	clientMap := make(map[string]*client.Client)
	for _, client := range filteredClients {
		clientMap[client.Id] = client
	}

	return clientMap
}

func GetMapValues[K comparable, V any](m map[K]V) []V {
    var values []V
    for _, value := range m {
        values = append(values, value)
    }
    return values
}

func GetBytes[T any](data T) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Failed to unmarshal ")
		return nil, err
	}
	return bytes, nil
}

func GetType[T any](bytes []byte) (T, error) {
	var t T
	err := json.Unmarshal(bytes, &t)
	if err != nil {
		fmt.Println("Failed to unmarshal ")
		return t, err
	}
	return t, nil
}
