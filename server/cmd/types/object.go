package types

type Vector3 struct {
	X, Y, Z float64
}

type Vector2 struct {
	X, Y float64
}

type LevelObjectVector struct {
	X, Y float32
}

type LevelObject struct {
	Id       uint8
	Name     string
	Position LevelObjectVector
	Rotation float32
	Size     LevelObjectVector
}
