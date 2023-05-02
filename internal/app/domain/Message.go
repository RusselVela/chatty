package domain

type MessageType int

const (
	TypeUser MessageType = iota
	TypeChannel
)

// Message is the domain struct passed along between the service and its clients
type Message struct {
	Id                int64       `json:"id"`
	PreviousMessageId int64       `json:"previousMessageId"`
	SourceId          string      `json:"sourceId"`
	TargetId          string      `json:"targetId"`
	Type              MessageType `json:"type"`
	Text              string      `json:"text"`
}
