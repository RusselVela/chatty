package inmemory

import (
	"fmt"
	"github.com/google/uuid"
)

type ChannelBean struct {
	Id         uuid.UUID
	Name       string
	OwnerId    string
	Visibility string
	Members    map[string]string
}

type channelsTable map[string]*ChannelBean

var channels channelsTable
var channelsByName map[string]string

func NewChannel(name string, ownerId string, visibility string) (*ChannelBean, error) {
	if channelId := channelsByName[name]; channelId != "" {
		return nil, fmt.Errorf("channel %s already exists", name)
	}
	id := uuid.New()
	channel := &ChannelBean{
		Id:         id,
		Name:       name,
		OwnerId:    ownerId,
		Visibility: visibility,
		Members:    make(map[string]string),
	}
	// OwnerId is the first member of the Channel
	channel.Members[ownerId] = ownerId

	channels[channel.Id.String()] = channel
	channelsByName[channel.Name] = channel.Id.String()

	return channel, nil
}

func GetChannel(id string) *ChannelBean {
	return channels[id]
}

func GetChannelByName(name string) *ChannelBean {
	id, found := channelsByName[name]
	if !found {
		return nil
	}
	return channels[id]
}
