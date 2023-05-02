package service

import (
	"fmt"
	"github.com/RusselVela/chatty/internal/app/datasourcce/repository/inmemory"
	"github.com/RusselVela/chatty/internal/app/domain"
)

var channelClients = make(map[string]*ChannelClient)

type ChannelClient struct {
	Channel     *inmemory.ChannelBean
	Subscribe   chan *inmemory.UserBean
	Unsubscribe chan *inmemory.UserBean
	Broadcaster chan domain.Message
}

// NewChannelClient returns a new client for Channel communication
func NewChannelClient(channel *inmemory.ChannelBean) *ChannelClient {
	channelClient := &ChannelClient{
		Channel:     channel,
		Subscribe:   make(chan *inmemory.UserBean),
		Unsubscribe: make(chan *inmemory.UserBean),
		Broadcaster: make(chan domain.Message),
	}
	channelClients[channel.Id.String()] = channelClient

	return channelClient
}

// Start initializes the loop for receiving messages from clients
func (cn *ChannelClient) Start() {
	for {
		select {
		case user := <-cn.Subscribe:
			cn.Channel.Members[user.Id.String()] = user.Id.String()
			msg := cn.createChannelMessage(fmt.Sprintf("User %s joined the channel", user.Username))
			cn.broadcastMessage(msg)
			break
		case user := <-cn.Unsubscribe:
			delete(cn.Channel.Members, user.Id.String())
			msg := cn.createChannelMessage(fmt.Sprintf("User %s left the channel", user.Username))
			cn.broadcastMessage(msg)
			break
		case message := <-cn.Broadcaster:
			cn.broadcastMessage(message)
		}
	}
}

func (cn *ChannelClient) createChannelMessage(message string) domain.Message {
	msg := domain.Message{
		SourceId: cn.Channel.Id.String(),
		TargetId: cn.Channel.Id.String(),
		Type:     domain.TypeChannel,
		Text:     message,
	}
	return msg
}

func (cn *ChannelClient) broadcastMessage(message domain.Message) {
	//TODO Save to Redis

	for _, userId := range cn.Channel.Members {
		wsClient := clients[userId]
		if wsClient != nil {
			wsClient.writeMessage(message)
		}
	}
}
