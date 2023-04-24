package inmemory

import "github.com/google/uuid"

func InitDatabase() {
	admin := User{
		Id:       uuid.New().String(),
		Username: "admin",
		Password: "admin",
	}
	Users = make(usersTable, 0)
	Users[admin.Username] = &admin
}
