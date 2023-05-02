package service

import (
	chattyredis "github.com/RusselVela/chatty/internal/app/datasourcce/repository/redis"
	"github.com/RusselVela/chatty/internal/app/domain"
)

var (
	broadcaster = make(chan domain.Message)
	clients     = make(map[string]*UserClient)
)

// HandleMessages is the exposed function to start the handleMessages routine
func HandleMessages(redisClient *chattyredis.Client) {
	go handleMessages(redisClient)
}

// handleMessages receives all incoming messages from clients and delivers them to the correct recipient
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
			wsClient.message <- msg
		case domain.TypeChannel:
			channel, found := channelClients[target]
			if !found {
				continue
			}
			channel.broadcastMessage(msg)
		}
	}
}

// removeClient removes a websocket struct from the pool of connected websockets
func removeClient(id string) {
	delete(clients, id)
}
