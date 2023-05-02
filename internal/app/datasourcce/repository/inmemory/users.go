package inmemory

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
)

// UserBean is the struct that saves user information
type UserBean struct {
	Id            uuid.UUID
	Username      string
	Password      string
	Subscriptions []string
	Online        bool
}

type usersTable map[string]*UserBean

var users usersTable
var usersByUsername map[string]string
var userMu sync.Mutex

// NewUser creates a new user in the table
func NewUser(username string, password string) (*UserBean, error) {
	userMu.Lock()
	defer func() {
		userMu.Unlock()
	}()

	if userId := usersByUsername[username]; userId != "" {
		return nil, fmt.Errorf("user %s already exists", username)
	}

	id := uuid.New()
	user := &UserBean{
		Id:       id,
		Username: username,
		Password: password,
		Online:   false,
	}
	users[user.Id.String()] = user
	usersByUsername[user.Username] = user.Id.String()

	return user, nil
}

// GetUser retrieves the user that matches the given id
func GetUser(id string) *UserBean {
	return users[id]
}

// GetUserByName retrieves the user that matches the given username
func GetUserByName(username string) *UserBean {
	id, found := usersByUsername[username]
	if !found {
		return nil
	}
	return users[id]
}

// GetUsers retrieves all users in the table
func GetUsers() []*UserBean {
	userList := make([]*UserBean, 0, len(users))
	for _, v := range users {
		userList = append(userList, v)
	}

	return userList
}
