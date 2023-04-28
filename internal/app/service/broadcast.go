package service

import "github.com/RusselVela/chatty/internal/app/domain"

var (
	broadcaster = make(chan domain.Message)
	clients     = make(map[string]*UserClient)
)

func HandleMessages() {
	go handleMessages()
}

func handleMessages() {
	for {
		msg := <-broadcaster
		// TODO: Store in redis
		target := msg.TargetId
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

func removeClient(id string) {
	delete(clients, id)
}
