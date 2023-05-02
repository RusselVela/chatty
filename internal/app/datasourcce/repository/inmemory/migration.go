package inmemory

// InitDatabase initializes two users for convenience
func InitDatabase() {
	//Initializing Users table
	users = make(usersTable, 0)
	usersByUsername = make(map[string]string, 0)
	// Initializing Channels table
	channels = make(channelsTable, 0)
	channelsByName = make(map[string]string)

	_, _ = NewUser("admin", "admin")
	_, _ = NewUser("rvela", "newPassword!")
}
