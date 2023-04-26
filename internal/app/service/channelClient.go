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

func NewChannelClient(channel *inmemory.ChannelBean) *ChannelClient {
	channelClient := &ChannelClient{
		Channel:     channel,
		Subscribe:   make(chan *inmemory.UserBean),
		Unsubscribe: make(chan *inmemory.UserBean),
		Broadcaster: make(chan domain.Message),
	}
	channelClients[channel.Name] = channelClient

	return channelClient
}

func (cn *ChannelClient) Start() {
	for {
		select {
		case user := <-cn.Subscribe:
			cn.Channel.Members[user.Username] = user.Username
			msg := cn.createChannelMessage(fmt.Sprintf("UserBean %s joined the channel", user.Username))
			cn.broadcastMessage(msg)
			break
		case user := <-cn.Unsubscribe:
			delete(cn.Channel.Members, user.Username)
			msg := cn.createChannelMessage(fmt.Sprintf("UserBean %s left the channel", user.Username))
			cn.broadcastMessage(msg)
			break
		case message := <-cn.Broadcaster:
			cn.broadcastMessage(message)
		}
	}
}

func (cn *ChannelClient) createChannelMessage(message string) domain.Message {
	msg := domain.Message{
		Sender:    cn.Channel.Name,
		Recipient: cn.Channel.Name,
		Type:      domain.TypeChannel,
		Text:      message,
	}
	return msg
}

func (cn *ChannelClient) broadcastMessage(message domain.Message) {
	//TODO Save to Redis

	for _, user := range cn.Channel.Members {
		wsClient := clients[user]
		if wsClient != nil {
			wsClient.writeMessage(message)
		}
	}
}
