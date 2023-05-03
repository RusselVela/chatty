package service

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ErrorSuite struct {
	suite.Suite
}

func (e *ErrorSuite) SetupTest() {
}

func TestError(t *testing.T) {
	suite.Run(t, new(ErrorSuite))
}

func (e *ErrorSuite) Test() {
	err := ErrorCode{
		Status:  400,
		Code:    9999,
		Message: "some message",
		Args:    nil,
	}

	e.Equal(err.Error(), err.String())
	cloneErr := err.Clone()
	e.Equal(err.Status, cloneErr.Status)
	e.Equal(err.Code, cloneErr.Code)

}
