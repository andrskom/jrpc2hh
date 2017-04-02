package handlers

import (
	"github.com/andrskom/jrpc2hh/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockService struct {
	mock.Mock
}

func (mc *MockService) Call(reqBody *models.RequestBody) (interface{}, *models.Error) {
	args := mc.Called(reqBody)
	return args.Get(0).(interface{}), args.Get(0).(*models.Error)
}

func TestHandler_Register(t *testing.T) {
	a := assert.New(t)
	h := NewHandler()
	h.Register(new(MockService))
	a.Len(h.sMap, 1)
	_, ok := h.sMap["MockService"]
	a.True(ok)
}
