package app

import "testing"

func TestLevel(t *testing.T) {
	// TODO DO NOT COMMIT THIS 
	t.Setenv("LEVEL_FILE_PATH", "/home/mikko/Projects/TankTactics/server/data/level.json")
	level := NewLevel()	
	if level == nil {
		t.Error("Cound not create level")	
	}
}