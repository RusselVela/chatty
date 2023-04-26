package inmemory

import "github.com/google/uuid"

func InitDatabase() {
	admin := UserBean{
		Id:       uuid.New().String(),
		Username: "admin",
		Password: "admin",
	}
	rvela := UserBean{
		Id:       uuid.New().String(),
		Username: "rvela",
		Password: "newPassword!",
	}

	//Initializing Users table
	Users = make(usersTable, 0)
	Users[admin.Username] = &admin
	Users[rvela.Username] = &rvela

	// Initializing Channels table
	Channels = make(channelsTable, 0)
}
