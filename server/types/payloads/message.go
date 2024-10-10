package payloads

type MessagePayload struct {
	Body string `json:"body"`
}

type PositionPayload struct {
	PosX int32 `json:"posx"`
	PosY int32 `json:"posy"`
}
