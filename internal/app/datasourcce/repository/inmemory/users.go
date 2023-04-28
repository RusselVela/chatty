package inmemory

import (
	"fmt"
	"github.com/google/uuid"
)

type UserBean struct {
	Id            uuid.UUID
	Username      string
	Password      string
	Subscriptions []string
}

type usersTable map[string]*UserBean

var users usersTable
var usersByUsername map[string]string

func NewUser(username string, password string) (*UserBean, error) {
	if userId := usersByUsername[username]; userId != "" {
		return nil, fmt.Errorf("user %s already exists", username)
	}

	id := uuid.New()
	user := &UserBean{
		Id:       id,
		Username: username,
		Password: password,
	}
	users[user.Id.String()] = user
	usersByUsername[user.Username] = user.Id.String()

	return user, nil
}

func GetUser(id string) *UserBean {
	return users[id]
}

func GetUserByName(username string) *UserBean {
	id, found := usersByUsername[username]
	if !found {
		return nil
	}
	return users[id]
}
