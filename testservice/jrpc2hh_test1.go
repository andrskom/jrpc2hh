package testservice

import (
	"fmt"
	jModels "github.com/andrskom/jrpc2hh/models"
	"encoding/json"
	"models"
	
)

func (s *Test1) Call(reqBody *jModels.RequestBody) (interface{}, *jModels.Error) {
	switch reqBody.GetMethod() {
	case "AnotherPackageResult":
		if reqBody.HasParams() {
			return nil, jModels.NewError(jModels.ErrorCodeInvalidParams, "That method of service can't has param", nil)
		}
		var args jModels.NilArgs
		var res models.SomeModel
		err := s.AnotherPackageResult(args, &res)
		if err != nil {
			return nil, jModels.NewError(jModels.ErrorCodeInternalError, "Internal error", err.Error())
		}
		return res, nil
	case "NilArgs":
		if reqBody.HasParams() {
			return nil, jModels.NewError(jModels.ErrorCodeInvalidParams, "That method of service can't has param", nil)
		}
		var args jModels.NilArgs
		var res .Test1NilArgsResult
		err := s.NilArgs(args, &res)
		if err != nil {
			return nil, jModels.NewError(jModels.ErrorCodeInternalError, "Internal error", err.Error())
		}
		return res, nil
	case "NilResult":
		var args Test1NilResultArgs
		if reqBody.HasParams() {
			err := json.Unmarshal(*reqBody.Params, &args)
			if err != nil {
				return nil, jModels.NewError(jModels.ErrorCodeInvalidParams, "Can't unmarshal params to args structure'", err.Error())
			}
		}
		var res jModels.NilResult
		err := s.NilResult(args, &res)
		if err != nil {
			return nil, jModels.NewError(jModels.ErrorCodeInternalError, "Internal error", err.Error())
		}
		return res, nil
	default:
		return nil, jModels.NewError(jModels.ErrorCodeMethodNotFound, fmt.Sprintf("Unknown method '%s' for service '%s'", reqBody.GetMethod(), "Test1"), nil)
	}
}