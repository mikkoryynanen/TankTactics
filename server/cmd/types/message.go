package types

type BaseMessage struct {
	Type int32 `json:"type"`
	ClientId string `json:"ClientId"`
}

type MessagePayload struct {
	Body string `json:"body"`
}

type PositionPayload struct {
	PosX float32 `json:"posx"`
	PosY float32 `json:"posy"`
}

type InputPayload struct {
	InputX int8 `json:"InputX"`
	InputY int8 `json:"InputY"`
}

type ServerState struct {
	BaseMessage

	PosX float32 `json:"posx"`
	PosY float32 `json:"posy"`
}

type PlayerMetadata struct {
	BaseMessage
}