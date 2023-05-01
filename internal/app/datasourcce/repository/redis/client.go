package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RusselVela/chatty/internal/app/domain"
	"github.com/knadh/koanf"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

const (
	redisConfigKey = "redis"
	redisUserKey   = "user-queue-%s"
)

type Config struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type Client struct {
	client *redis.Client
}

func ConfigureRedisDB(k *koanf.Koanf) (*Client, error) {
	redisConfig, err := readRedisConfig(k)
	if err != nil {
		return nil, err
	}
	return newRedisClient(redisConfig), nil
}

func readRedisConfig(k *koanf.Koanf) (*Config, error) {
	redisConfig := &Config{}
	if err := k.UnmarshalWithConf(redisConfigKey, redisConfig, koanf.UnmarshalConf{Tag: "json"}); err != nil {
		return nil, err
	}
	return redisConfig, nil
}

func newRedisClient(config *Config) *Client {
	options := &redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.DB,
	}
	client := redis.NewClient(options)

	return &Client{client: client}
}

func (c *Client) StoreMessage(message domain.Message) {
	id := time.Now().UnixMilli()
	message.Id = id
	jsonMsg, err := json.Marshal(message)
	if err != nil {
		zap.S().Errorf("failed to marshall message for storing: %v", message)
	}
	switch message.Type {
	case domain.TypeUser:
		key := fmt.Sprintf(redisUserKey, message.TargetId)
		c.client.RPush(context.Background(), key, jsonMsg)
		break
	case domain.TypeChannel:
		break
	}

}

func (c *Client) GetMessages(userID string) []*domain.Message {
	messages := make([]*domain.Message, 0)
	key := fmt.Sprintf(redisUserKey, userID)

	result, err := c.client.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return messages
	}
	for _, m := range result {
		msg := &domain.Message{}
		_ = json.Unmarshal([]byte(m), msg)
		messages = append(messages, msg)
	}

	return messages
}

func (c *Client) ClearMessageQueue(userID string) (int64, error) {
	key := fmt.Sprintf(redisUserKey, userID)

	result, err := c.client.Del(context.Background(), key).Result()
	if err != nil {
		zap.S().Errorf("error deleting message queue from db: %s", err.Error())
		return -1, err
	}

	return result, nil
}
