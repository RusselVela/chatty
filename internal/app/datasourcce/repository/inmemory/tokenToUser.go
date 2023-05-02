package inmemory

import "sync"

type tokenToUsersTable map[string]string

var tokenToUser tokenToUsersTable
var tokenMu sync.Mutex

// AddTokenToUser adds a new pair of [user, token] to this table.
// When the token is used to authenticate the websocket connection, the entry is deleted.
func AddTokenToUser(userId string, token string) {
	tokenMu.Lock()
	tokenToUser[userId] = token
	tokenMu.Unlock()
}

// GetToken retrieves the token associated with the userId, if any.
func GetToken(userId string) string {
	return tokenToUser[userId]
}

// DeleteTokenToUser removes the [userId, token] pair from the table
func DeleteTokenToUser(userId string) {
	tokenMu.Lock()
	delete(tokenToUser, userId)
	tokenMu.Unlock()
}
