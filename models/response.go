package models

import (
	"encoding/json"
	"net/http"
)

type ResponseBody struct {
	JsonRpc string           `json:"jsonrpc"`
	Result  *json.RawMessage `json:"result,omitempty"`
	Error   *Error           `json:"error,omitempty"`
	Id      *uint            `json:"id"`
}

func NewResponseError(error *Error, id *uint) *ResponseBody {
	return &ResponseBody{
		JsonRpc: "2.0",
		Error:   error,
		Id:      id,
	}
}

func NewResponseBody(result *json.RawMessage, id *uint) *ResponseBody {
	return &ResponseBody{
		JsonRpc: "2.0",
		Result:  result,
		Id:      id,
	}
}

type JsonRpcResult string

const (
	JsonRpcResultOk JsonRpcResult = "Ok"
)

func JsonResponse(w http.ResponseWriter, data interface{}, code int) (int, error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	body, err := json.Marshal(data)
	if err != nil {
		body = []byte(
			`{"error": "Unknown and unpredictable error with huge, massive and catastrophic consequences!"}`)
	}

	return w.Write(body)
}
