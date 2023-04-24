package inmemory

import (
	"fmt"
	"github.com/google/uuid"
)

type User struct {
	Id       string
	Username string
	Password string
}

type usersTable map[string]*User

var Users usersTable

func (u usersTable) NewUser(username string, password string) (error, *User) {
	if user := u[username]; user != nil {
		return fmt.Errorf("user %s already exists", username), nil
	}

	id := uuid.New().String()
	user := User{
		Id:       id,
		Username: username,
		Password: password,
	}
	u[username] = &user

	return nil, &user
}

func (u usersTable) Get(username string) *User {
	return u[username]
}
