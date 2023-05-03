package service

import (
	"github.com/RusselVela/chatty/internal/app/datasourcce/repository/inmemory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

type JwtSuite struct {
	suite.Suite
}

func (j *JwtSuite) SetupTest() {
}

func TestJwt(t *testing.T) {
	suite.Run(t, new(JwtSuite))
}

func (j *JwtSuite) Test() {
	JWTSecret = []byte("fooBar123")

	user := inmemory.UserBean{
		Id:            uuid.New(),
		Username:      "foo",
		Password:      "bar",
		Subscriptions: nil,
		Online:        false,
	}

	token, err := generateJWT(user, 10)
	j.Nil(err)
	j.NotNil(token)

	_, err = parseJWT(token)
	j.Nil(err)
}
