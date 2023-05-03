package service

import (
	"github.com/RusselVela/chatty/internal/app/datasourcce/repository/inmemory"
	"github.com/RusselVela/chatty/internal/app/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ChannelClientSuite struct {
	suite.Suite
}

func (ccs *ChannelClientSuite) SetupTest() {
}

func TestChannelClient(t *testing.T) {
	suite.Run(t, new(ChannelClientSuite))
}

func (ccs *ChannelClientSuite) TestNew() {
	ch := &inmemory.ChannelBean{
		Id:         uuid.New(),
		Name:       "TestChannel",
		OwnerId:    "4567",
		Visibility: "public",
		Members:    map[string]string{"4567": ""},
	}
	channelClient := NewChannelClient(ch, nil)
	ccs.NotNil(channelClient.Channel)
	ccs.NotNil(channelClient.Subscribe)
	ccs.NotNil(channelClient.Unsubscribe)
	ccs.NotNil(channelClient.Broadcaster)

}

func (css *ChannelClientSuite) TestCreateChannelMessage() {
	ch := &inmemory.ChannelBean{
		Id:         uuid.New(),
		Name:       "TestChannel",
		OwnerId:    "4567",
		Visibility: "public",
		Members:    map[string]string{"4567": ""},
	}
	channelClient := NewChannelClient(ch, nil)
	msg := channelClient.createChannelMessage("fooBar123")
	css.Equal(ch.Id.String(), msg.SourceId)
	css.Equal(ch.Id.String(), msg.TargetId)
	css.Equal(domain.TypeChannel, msg.Type)
	css.Equal("fooBar123", msg.Text)
}
