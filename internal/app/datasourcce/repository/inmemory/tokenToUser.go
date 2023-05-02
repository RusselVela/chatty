package inmemory

type tokenToUsersTable map[string]string

var tokenToUser tokenToUsersTable

// AddTokenToUser adds a new pair of [user, token] to this table.
// When the token is used to authenticate the websocket connection, the entry is deleted.
func AddTokenToUser(userId string, token string) {
	tokenToUser[userId] = token
}

// GetToken retrieves the token associated with the userId, if any.
func GetToken(userId string) string {
	return tokenToUser[userId]
}

// DeleteTokenToUser removes the [userId, token] pair from the table
func DeleteTokenToUser(userId string) {
	delete(tokenToUser, userId)
}
