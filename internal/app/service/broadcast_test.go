package service

import (
	"encoding/json"
	"fmt"
	"github.com/RusselVela/chatty/internal/app/datasourcce/repository/redis"
	"github.com/RusselVela/chatty/internal/app/domain"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/suite"
	"testing"
)

type BroadcastSuite struct {
	suite.Suite
}

func (bs *BroadcastSuite) SetupTest() {
}

func TestBroadcast(t *testing.T) {
	suite.Run(t, new(BroadcastSuite))
}

func (bs *BroadcastSuite) TestHandleMessages() {
	bs.Run("TestUserBroadcast", func() {
		db, mock := redismock.NewClientMock()
		client := redis.NewRedisClient(db)

		id := "zxc890"
		key := fmt.Sprintf("user-queue-%s", id)

		msg := domain.Message{
			Id:                123,
			PreviousMessageId: 122,
			SourceId:          "abc123",
			TargetId:          id,
			Type:              0,
			Text:              "hello world!",
		}
		jsonMsg, _ := json.Marshal(msg)

		mock.ExpectRPush(key, jsonMsg)
		HandleMessages(client)

		broadcaster <- msg

		if err := mock.ExpectationsWereMet(); err != nil {
			bs.T().Error(err)
		}
	})

	bs.Run("TestChannelBroadcast", func() {
		db, mock := redismock.NewClientMock()
		client := redis.NewRedisClient(db)

		id := "klj964"
		key := fmt.Sprintf("user-queue-%s", id)

		msg := domain.Message{
			Id:                123,
			PreviousMessageId: 122,
			SourceId:          "abc123",
			TargetId:          id,
			Type:              1,
			Text:              "hello channel!",
		}
		jsonMsg, _ := json.Marshal(msg)

		mock.ExpectRPush(key, jsonMsg)
		HandleMessages(client)

		broadcaster <- msg

		if err := mock.ExpectationsWereMet(); err != nil {
			bs.T().Error(err)
		}
	})
}
