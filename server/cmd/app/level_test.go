package app

import (
	"main/cmd/utils"
	"path/filepath"
	"testing"
)

func TestLevel(t *testing.T) {
	levelPath, _ := utils.FindProjectRoot()
	levelPath = filepath.Join(levelPath, "data", "level.json")
	t.Setenv("LEVEL_FILE_PATH", levelPath)

	level := NewLevel()	
	if level == nil {
		t.Error("Cound not create level")	
	}
}