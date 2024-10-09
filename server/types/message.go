package messageTypes

type Message struct {
	Body string `json:"body"`
}

type Position struct {
	PosX int32 `json:"posx"`
	PosY int32 `json:"posy"`
}
