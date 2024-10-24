package app

import (
	"encoding/json"
	"log"
	"main/cmd/types"
	"os"
)

// Level represents the loaded information about a level. This includes the static objects inside the world
type Level struct {
	Objects []types.LevelObject
}

func NewLevel() *Level {
	// TODO Load this from S3 or similiar
	data, err := os.ReadFile(os.Getenv("LEVEL_FILE_PATH"))
	if err != nil {
		log.Print("could not load level.json")
		return nil
	}

	var level Level
	err = json.Unmarshal(data, &level)
	if err != nil {
		log.Print("Could not unmarshal level data")
		return nil
	}

	// Load level data
	return &level
}

func (l *Level) IsObjectColliding(a, b types.LevelObject) bool {
	overlap := a.Position.X < b.Position.X + b.Size.X &&
		a.Position.X + a.Size.X > b.Position.X &&
		a.Position.Y < b.Position.Y + b.Size.Y &&
		a.Position.Y + a.Size.Y > b.Position.Y
	return overlap
}
