package testservice

import (
	"fmt"
	json "encoding/json"
	jModels "github.com/andrskom/jrpc2hh/models"
	models "models"
	models_2561278490898348873 "some/models"
	
)

func (s *Test2) Call(reqBody *jModels.RequestBody) (interface{}, *jModels.Error) {
	switch reqBody.GetMethod() {
	case "NilArgs":
		if reqBody.HasParams() {
			return nil, jModels.NewError(jModels.ErrorCodeInvalidParams, "That method of service can't has param", nil)
		}
		var args jModels.NilArgs
		var res Test2NilArgsResult
		err := s.NilArgs(args, &res)
		if err != nil {
			return nil, jModels.NewError(jModels.ErrorCodeInternalError, "Internal error", err.Error())
		}
		return res, nil
	case "NilResult":
		var args models.Test2NilResultArgs
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
	case "AnotherPackageResult":
		var args models_2561278490898348873.NilArgs
		if reqBody.HasParams() {
			err := json.Unmarshal(*reqBody.Params, &args)
			if err != nil {
				return nil, jModels.NewError(jModels.ErrorCodeInvalidParams, "Can't unmarshal params to args structure'", err.Error())
			}
		}
		var res models.SomeModel
		err := s.AnotherPackageResult(args, &res)
		if err != nil {
			return nil, jModels.NewError(jModels.ErrorCodeInternalError, "Internal error", err.Error())
		}
		return res, nil
	case "DoubleStarAnotherResult":
		if reqBody.HasParams() {
			return nil, jModels.NewError(jModels.ErrorCodeInvalidParams, "That method of service can't has param", nil)
		}
		var args jModels.NilArgs
		var res *models.SomeModel
		err := s.DoubleStarAnotherResult(args, &res)
		if err != nil {
			return nil, jModels.NewError(jModels.ErrorCodeInternalError, "Internal error", err.Error())
		}
		return res, nil
	case "DoubleStarResult":
		if reqBody.HasParams() {
			return nil, jModels.NewError(jModels.ErrorCodeInvalidParams, "That method of service can't has param", nil)
		}
		var args jModels.NilArgs
		var res *Test2NilArgsResult
		err := s.DoubleStarResult(args, &res)
		if err != nil {
			return nil, jModels.NewError(jModels.ErrorCodeInternalError, "Internal error", err.Error())
		}
		return res, nil
	default:
		return nil, jModels.NewError(jModels.ErrorCodeMethodNotFound, fmt.Sprintf("Unknown method '%s' for service '%s'", reqBody.GetMethod(), "Test2"), nil)
	}
}