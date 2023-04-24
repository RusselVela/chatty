package service

import "time"

type ChattyService struct {
}

func NewChattyService() *ChattyService {
	return &ChattyService{}
}

func (cs *ChattyService) Login(username string, password string) (string, error) {
	return "", nil
}

func (cs *ChattyService) PostMessage(token string, recipient string, message string) (int64, error) {
	return time.Now().Unix(), nil
}
