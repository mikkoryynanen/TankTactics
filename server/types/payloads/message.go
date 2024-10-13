package payloads

type BaseMessage struct {
	Type int32 `json:"type"`
	ClientId string `json:"ClientId"`
}

type MessagePayload struct {
	Body string `json:"body"`
}

type PositionPayload struct {
	PosX int32 `json:"posx"`
	PosY int32 `json:"posy"`
}

type InputPayload struct {
	InputX int8 `json:"InputX"`
	InputY int8 `json:"InputY"`
}
