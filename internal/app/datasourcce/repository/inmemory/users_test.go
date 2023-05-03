package inmemory

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UsersSuite struct {
	suite.Suite
}

func (us *UsersSuite) SetupTest() {
	InitDatabase()
}

func TestUsers(t *testing.T) {
	suite.Run(t, new(UsersSuite))
}

func (us *UsersSuite) TestUsers_NewUser() {
	u, err := NewUser("admin", "admin")
	us.NotNil(err)

	u, err = NewUser("foo", "123")

	u = GetUser(u.Id.String())

	us.Nil(err)
	us.NotNil(u)
	us.Equal("foo", u.Username)
	us.Equal("123", u.Password)

	u = GetUserByName("bar")
	us.Nil(u)

	u = GetUserByName("admin")
	us.NotNil(u)
	us.Equal("admin", u.Username)
	us.Equal("admin", u.Password)

	uList := GetUsers()
	us.Equal(3, len(uList))
}
