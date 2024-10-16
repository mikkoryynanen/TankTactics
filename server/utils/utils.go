package utils

import client "main/types/Client"

func RemoveDisconnectedClients(clients []*client.Client) []*client.Client {
	var filteredClients []*client.Client
	for _, client := range clients {
		if client.IsConnected {
			filteredClients = append(filteredClients, client)
		}
	}

	return filteredClients
}

func GetMapValues[K comparable, V any](m map[K]V) []V {
    var values []V
    for _, value := range m {
        values = append(values, value)
    }
    return values
}