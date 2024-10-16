package utils

import (
	client "main/types/Client"
	"math"
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

func RoundFloatTo(value float64) float64 {
	// Check if the value is less than 0.01
	// if math.Abs(value) < 0.01 {
	// 	return value // Return the original value if it's less than 0.01
	// }

	// Multiply by 100, round, then divide by 100 to get two decimal places
	return math.Round(value*100) / 100
}
