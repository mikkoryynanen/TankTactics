package app

import (
	"encoding/json"
	"log"
	"main/cmd/types"
	"main/cmd/utils"
	"os"
	"path/filepath"
)

// Level represents the loaded information about a level. This includes the static objects inside the world
type Level struct {
	Objects []types.LevelObject
}

func NewLevel() *Level {
	levelPath, _ := utils.FindProjectRoot()
	levelPath = filepath.Join(levelPath, "data", "level.json")
	data, err := os.ReadFile(levelPath)
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
