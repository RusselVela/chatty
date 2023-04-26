package inmemory

import (
	"fmt"
	"github.com/google/uuid"
)

type UserBean struct {
	Id            string
	Username      string
	Password      string
	Subscriptions []string
}

type usersTable map[string]*UserBean

var Users usersTable

func (u usersTable) NewUser(username string, password string) (*UserBean, error) {
	if user := u[username]; user != nil {
		return nil, fmt.Errorf("user %s already exists", username)
	}

	id := uuid.New().String()
	user := &UserBean{
		Id:       id,
		Username: username,
		Password: password,
	}
	u[username] = user

	return user, nil
}

func (u usersTable) Get(username string) *UserBean {
	return u[username]
}
