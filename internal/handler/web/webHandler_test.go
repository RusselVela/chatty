package web

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type WebHandlerSuite struct {
	suite.Suite
}

func (wh *WebHandlerSuite) SetupTest() {
}

func TestWebHandler(t *testing.T) {
	suite.Run(t, new(WebHandlerSuite))
}

func (wh *WebHandlerSuite) TestWebHandler() {
	//lifecycle := fxtest.NewLifecycle(wh.T())
	//
	//e := echo.New()
	//wh := NewWebHandler(nil)
}
