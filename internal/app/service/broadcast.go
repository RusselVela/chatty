package service

import "github.com/RusselVela/chatty/internal/app/domain"

var (
	broadcaster = make(chan domain.Message)
	clients     = make(map[string]*WsClient)
)

func HandleMessages() {
	go handleMessages()
}

func handleMessages() {
	for {
		msg := <-broadcaster
		// Store in redis

		switch msg.Type {
		case domain.TypeUser:
			target := msg.Recipient
			wsClient, found := clients[target]
			if !found {
				continue
			}
			wsClient.writeMessage(msg)
		case domain.TypeChannel:
		}
	}
}

func removeClient(username string) {
	delete(clients, username)
}
