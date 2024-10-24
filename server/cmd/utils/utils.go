package utils

import (
	"fmt"
    "os"
    "path/filepath"
	"encoding/json"
	"main/cmd/types" 
)

func RemoveDisconnectedClients(clients []*types.Client) map[string]*types.Client {
	var filteredClients []*types.Client
	for _, client := range clients {
		if client.IsConnected {
			filteredClients = append(filteredClients, client)
		}
	}

	clientMap := make(map[string]*types.Client)
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

func FindProjectRoot() (string, error) {
    cwd, err := os.Getwd()
    if err != nil {
        return "", err
    }

    // Traverse upwards from the current working directory
    for cwd != "/" {
        // Check if "go.mod" exists in the current directory (indicating root)
        if _, err := os.Stat(filepath.Join(cwd, "go.mod")); !os.IsNotExist(err) {
            return cwd, nil
        }
        cwd = filepath.Dir(cwd) // Go up one directory
    }

    return "", fmt.Errorf("could not find project root")
}
