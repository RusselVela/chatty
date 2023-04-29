package service

import (
	chattyredis "github.com/RusselVela/chatty/internal/app/datasourcce/repository/redis"
	"github.com/RusselVela/chatty/internal/app/domain"
)

var (
	broadcaster = make(chan domain.Message)
	clients     = make(map[string]*UserClient)
)

func HandleMessages(redisClient *chattyredis.Client) {
	go handleMessages(redisClient)
}

func handleMessages(redisClient *chattyredis.Client) {
	for {
		msg := <-broadcaster

		target := msg.TargetId
		switch msg.Type {
		case domain.TypeUser:
			wsClient, found := clients[target]
			if !found {
				redisClient.StoreMessage(msg)
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
