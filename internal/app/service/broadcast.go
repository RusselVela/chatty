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
		target := msg.Recipient
		switch msg.Type {
		case domain.TypeUser:
			wsClient, found := clients[target]
			if !found {
				continue
			}
			wsClient.writeMessage(msg)
		case domain.TypeChannel:
			channel, found := channelClients[target]
			if !found {
				continue
			}
			channel.broadcastMessage(msg)
		}
	}
}

func removeClient(username string) {
	delete(clients, username)
}
