package models

import (
	"encoding/json"
	"errors"
	"strings"
	"unicode/utf8"
	"net/http"
)

type RequestBody struct {
	JsonRpc string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Id      *uint             `json:"id"`
	Params  *json.RawMessage `json:"params,omitempty"`
}

func (r *RequestBody) Validate() error {
	if r.JsonRpc != "2.0" {
		return errors.New("Bad request, field 'jsonrpc' is empty")
	}

	dotI := strings.Index(r.Method, ".")
	if dotI < 0 || strings.Count(r.Method, ".") != 1 ||
		dotI == 0 || dotI == (utf8.RuneCountInString(r.Method)-1){
		return errors.New("Bad request, bad format field 'method'")
	}

	if r.Id == nil || *r.Id == 0 {
		return errors.New("Bad request, bad format field 'id'")
	}

	return nil
}

func (r *RequestBody) GetService() string {
	s := strings.Split(r.Method, ".")
	return s[0]
}

func (r *RequestBody) GetMethod() string {
	s := strings.Split(r.Method, ".")
	return s[1]
}

func (r *RequestBody) HasParams() bool {
	return r.Params != nil
}

func ValidateHeaders(req *http.Request) *Error {
	if hAccept := req.Header.Get("Accept"); hAccept != "application/json" {
		return NewError(ErrorCodeInvalidRequest, "Header 'Accept' is not 'application/json'", nil)
	}
	if hContentType := req.Header.Get("Content-Type"); strings.Count(hContentType, "application/json") != 1 {
		return NewError(ErrorCodeInvalidRequest, "Header 'Content-Type' is not 'application/json'", nil)
	}
	return nil
}