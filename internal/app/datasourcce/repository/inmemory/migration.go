package inmemory

import "github.com/google/uuid"

func InitDatabase() {
	admin := UserBean{
		Id:       uuid.New(),
		Username: "admin",
		Password: "admin",
	}
	rvela := UserBean{
		Id:       uuid.New(),
		Username: "rvela",
		Password: "newPassword!",
	}

	//Initializing Users table
	users = make(usersTable, 0)
	users[admin.Id.String()] = &admin
	users[rvela.Id.String()] = &rvela

	// Initializing Channels table
	channels = make(channelsTable, 0)
}
