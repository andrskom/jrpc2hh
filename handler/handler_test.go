package handlers

import (
	"testing"
	"context"
	"github.com/andrskom/jrpc2hh/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
)

type MockService struct {
	mock.Mock
}

func (mc *MockService) Call(method string) (interface{}, *models.Error) {
	args := mc.Called(method)
	return args.Get(0).(interface{}), args.Get(0).(*models.Error)
}

func TestHandler_Register(t *testing.T) {
	a := assert.New(t)
	h := NewHandler(new(context.Context))
	h.Register(new(MockService))
	a.Len(h.sMap, 1)
	_, ok := h.sMap["MockService"]
	a.True(ok)
}
