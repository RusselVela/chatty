package inmemory

import (
	"fmt"
	"github.com/google/uuid"
)

type ChannelBean struct {
	Id         string
	Name       string
	Owner      string
	Visibility string
	Members    map[string]string
}

type channelsTable map[string]*ChannelBean

var Channels channelsTable

func (ct channelsTable) NewChannel(name string, owner string, visibility string) (*ChannelBean, error) {
	if channel := Channels[name]; channel != nil {
		return nil, fmt.Errorf("channel %s already exists", name)
	}
	id := uuid.New().String()
	channel := &ChannelBean{
		Id:         id,
		Name:       name,
		Owner:      owner,
		Visibility: visibility,
		Members:    make(map[string]string),
	}
	// Owner is the first member of the Channel
	channel.Members[owner] = owner

	ct[channel.Name] = channel

	return channel, nil
}

func (ct channelsTable) Get(name string) *ChannelBean {
	return ct[name]
}
