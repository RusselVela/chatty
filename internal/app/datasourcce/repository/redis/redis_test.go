package redis

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/RusselVela/chatty/internal/app/domain"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/suite"
)

type RedisSuite struct {
	suite.Suite
	client Client
	mock   redismock.ClientMock
}

func (rs *RedisSuite) SetupTest() {
	db, mock := redismock.NewClientMock()
	rs.client = Client{client: db}
	rs.mock = mock
}

func TestRedis(t *testing.T) {
	suite.Run(t, new(RedisSuite))
}

func (rs *RedisSuite) TestRedis() {
	id := "zxc890"
	key := fmt.Sprintf(redisUserKey, id)

	rs.Run("StoreMessage", func() {

		msg := domain.Message{
			Id:                123,
			PreviousMessageId: 122,
			SourceId:          "abc123",
			TargetId:          id,
			Type:              0,
			Text:              "hello world!",
		}
		jsonMsg, _ := json.Marshal(msg)
		rs.mock.ExpectRPush(key, jsonMsg)
		rs.client.StoreMessage(msg)

		if err := rs.mock.ExpectationsWereMet(); err != nil {
			rs.T().Error(err)
		}
	})

	rs.Run("GetMessages", func() {
		rs.mock.ExpectLRange(key, 0, -1)

		rs.client.GetMessages(id)

		if err := rs.mock.ExpectationsWereMet(); err != nil {
			rs.T().Error(err)
		}
	})

	rs.Run("ClearMessageQueue", func() {
		rs.mock.ExpectDel(key)

		rs.client.ClearMessageQueue(id)

		if err := rs.mock.ExpectationsWereMet(); err != nil {
			rs.T().Error(err)
		}

	})
}
