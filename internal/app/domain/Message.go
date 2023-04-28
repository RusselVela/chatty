package domain

type MessageType int

const (
	TypeUser MessageType = iota
	TypeChannel
)

type Message struct {
	SourceId string      `json:"sourceId"`
	TargetId string      `json:"targetId"`
	Type     MessageType `json:"type"`
	Text     string      `json:"text"`
}
