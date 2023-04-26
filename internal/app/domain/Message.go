package domain

type MessageType int

const (
	TypeUser MessageType = iota
	TypeChannel
)

type Message struct {
	Sender    string      `json:"sender,omitempty"`
	Recipient string      `json:"recipient"`
	Type      MessageType `json:"type"`
	Text      string      `json:"text"`
}
