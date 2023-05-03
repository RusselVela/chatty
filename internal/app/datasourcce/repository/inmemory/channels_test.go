package inmemory

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ChannelsSuite struct {
	suite.Suite
}

func (cs *ChannelsSuite) SetupTest() {
	InitDatabase()
	NewChannel("foo", "456", "public")
}

func TestChannels(t *testing.T) {
	suite.Run(t, new(ChannelsSuite))
}

func (cs *ChannelsSuite) TestChannels_NewChannel() {
	c, err := NewChannel("channel", "123", "public")
	cs.Nil(err)

	c = GetChannel(c.Id.String())
	cs.Equal("channel", c.Name)
	cs.Equal("123", c.OwnerId)
	cs.Equal("public", c.Visibility)

	c, err = NewChannel("foo", "456", "public")
	cs.NotNil(err)

	c = GetChannelByName("foo")
	cs.NotNil(c)

	c = GetChannelByName("bar")
	cs.Nil(c)

	list := GetChannels()
	cs.NotNil(list)
	cs.Equal(2, len(list))
}
