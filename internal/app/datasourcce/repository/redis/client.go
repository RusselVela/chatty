package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RusselVela/chatty/internal/app/domain"
	"github.com/knadh/koanf"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"sync"
)

const (
	redisConfigKey = "redis"
	redisUserKey   = "user-queue-%s"
)

var redisMu sync.Mutex

// Config represents a redis instance configuration
type Config struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// Client holds all the convenience functions to store and retrieve from redis instance referenced by client
type Client struct {
	client *redis.Client
}

func NewRedisClient(client *redis.Client) *Client {
	return &Client{client: client}
}

// ConfigureRedisDB takes config values from koanf to establish connection with the redis instance
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

// StoreMessage takes a domain.Message and stores it on the redis instance
func (c *Client) StoreMessage(message domain.Message) {
	redisMu.Lock()
	defer func() {
		redisMu.Unlock()
	}()

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
		key := fmt.Sprintf(redisUserKey, message.TargetId)
		c.client.RPush(context.Background(), key, jsonMsg)
		break
	}
}

// GetMessages retrieves all messages mapped to the given key userID from the redis instance
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

// ClearMessageQueue removes all messages under userID after they have been delivered
func (c *Client) ClearMessageQueue(userID string) (int64, error) {
	redisMu.Lock()
	defer func() {
		redisMu.Unlock()
	}()

	key := fmt.Sprintf(redisUserKey, userID)

	result, err := c.client.Del(context.Background(), key).Result()
	if err != nil {
		zap.S().Errorf("error deleting message queue from db: %s", err.Error())
		return -1, err
	}

	return result, nil
}
