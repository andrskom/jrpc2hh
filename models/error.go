package models

import "fmt"

type ErrorCode int

const (
	ErrorCodeParseError     ErrorCode = -32700
	ErrorCodeInvalidRequest ErrorCode = -32600
	ErrorCodeMethodNotFound ErrorCode = -32601
	ErrorCodeInvalidParams  ErrorCode = -32602
	ErrorCodeInternalError  ErrorCode = -32603
)

type Error struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewError(code ErrorCode, mes string, data interface{}) *Error {
	return &Error{code, mes, data}
}

func (e *Error) Error() string {
	return fmt.Sprintf("Code: %d, Message: '%s', Data, %+v", e.Code, e.Message, e.Data)
}
