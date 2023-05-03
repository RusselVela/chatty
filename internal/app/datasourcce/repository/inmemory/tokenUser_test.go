package inmemory

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TokenUserSuite struct {
	suite.Suite
}

func (tu *TokenUserSuite) SetupTest() {
	InitDatabase()
}

func TestTokenUser(t *testing.T) {
	suite.Run(t, new(TokenUserSuite))
}

func (tu *TokenUserSuite) TestTokenUserAddTokenUser() {
	AddTokenToUser("123", "456789")
	token := GetToken("123")
	tu.NotNil(token)
	tu.Equal("456789", token)

	DeleteTokenToUser("123")
	token = GetToken("123")
	tu.Empty(token)
}
