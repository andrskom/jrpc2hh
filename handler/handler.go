package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"context"
	"github.com/andrskom/jrpc2hh/models"
	"errors"
	"sync"
	"reflect"
)

type Caller interface {
	Call(reqBody *models.RequestBody) (interface{}, *models.Error)
}

type Handler struct {
	ctx *context.Context
	mu sync.Mutex
	sMap map[string]Caller
}

func NewHandler(ctx *context.Context) *Handler {
	return &Handler{ctx: ctx, sMap: make(map[string]Caller)}
}

func (h *Handler) Register(c Caller) error {
	t := reflect.TypeOf(c)
	var n string
	if t.Kind() == reflect.Ptr {
		n = t.Elem().Name()
	} else {
		n = t.Name()
	}
	return h.RegisterName(n, c)
}

func (h *Handler) RegisterName(name string, c Caller) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.sMap[name]; ok {
		return errors.New(fmt.Sprintf("Service with name '%s' already registered", name))
	}
	h.sMap[name] = c
	return nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	jErr := models.ValidateHeaders(req)
	if jErr != nil {
		models.JsonResponse(w, models.NewResponseError(jErr, nil), http.StatusBadRequest)
		return
	}

	reqBodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jErr := models.NewError(models.ErrorCodeInternalError, "Can't read request body", err.Error())
		models.JsonResponse(w, models.NewResponseError(jErr, nil), http.StatusInternalServerError)
		return
	}
	req.Body.Close()

	// If hasn't error, that it is simple request
	var jReq *models.RequestBody
	errReq := json.Unmarshal(reqBodyBytes, &jReq)
	if errReq == nil {
		err := jReq.Validate()
		if err != nil {
			jErr := models.NewError(models.ErrorCodeParseError, err.Error(), nil)
			models.JsonResponse(w, models.NewResponseError(jErr, jReq.Id), http.StatusBadRequest)
			return
		}
		s, ok := h.sMap[jReq.GetService()]
		if !ok {
			jErr := models.NewError(models.ErrorCodeMethodNotFound, "Unknown service", nil)
			models.JsonResponse(w, models.NewResponseError(jErr, jReq.Id), http.StatusNotFound)
			return
		}
		res, jErr := s.Call(jReq)
		if jErr != nil {
			models.JsonResponse(w, models.NewResponseError(jErr, jReq.Id), http.StatusInternalServerError)
			return
		}

		if res == nil {
			models.JsonResponse(w, models.NewResponseBody(nil, jReq.Id), http.StatusOK)
		} else {
			resByte, err := json.Marshal(res)
			if err != nil {
				jErr := models.NewError(models.ErrorCodeInternalError, "Can't marshal response", err.Error())
				models.JsonResponse(w, models.NewResponseError(jErr, jReq.Id), http.StatusInternalServerError)
				return
			}
			models.JsonResponse(w, models.NewResponseBody(&json.RawMessage(resByte), jReq.Id), http.StatusOK)
		}

		return
	}

	// If hasn't error, that it is batch request
	//var jReqBatchSlice []json.RawMessage
	//errReqBatch := json.Unmarshal(reqBodyBytes, &jReqBatchSlice)
	//if errReqBatch == nil {
	//	jRespBatchSlice := struct {
	//		sync.Mutex
	//		Slice []*models.JsonRpcResponse
	//	}{Slice: make([]*models.JsonRpcResponse, 0)}
	//	w := sync.WaitGroup{}
	//	for _, req := range jReqBatchSlice {
	//		w.Add(1)
	//		go func(req json.RawMessage) {
	//			defer w.Done()
	//			var jReqBatch *models.JsonRpcRequest
	//			err := json.Unmarshal(req, &jReqBatch)
	//			if err != nil {
	//				jRespBatchSlice.Lock()
	//				defer jRespBatchSlice.Unlock()
	//				jRespBatchSlice.Slice = append(
	//					jRespBatchSlice.Slice,
	//					models.NewJsonRpcResponseError(
	//						models.JsonRpcErrorCodeInvalidRequest,
	//						err.Error(),
	//						nil,
	//						nil))
	//			} else {
	//				respBody, _ := h.call(jReqBatch)
	//				jRespBatchSlice.Lock()
	//				defer jRespBatchSlice.Unlock()
	//				jRespBatchSlice.Slice = append(
	//					jRespBatchSlice.Slice,
	//					respBody)
	//			}
	//		}(req)
	//	}
	//	w.Wait()
	//	web.JSONResponse(resp, jRespBatchSlice.Slice, http.StatusOK)
	//	return
	//}
	//
	//// Else that error
	//web.JSONResponse(
	//	resp,
	//	models.NewJsonRpcResponseError(
	//		models.JsonRpcErrorCodeParseError,
	//		"Can't parse request json to json rpc 2.0 struct",
	//		nil,
	//		nil),
	//	http.StatusBadRequest)
	//return
}

//func (h *JsonRpcBaseHandler) call(reqBody *models.JsonRpcRequest) (*models.JsonRpcResponse, int) {
//	plan := NewPlan(h.db)
//	runPlan := NewRunPlan(h.db)
//	job := NewJob(h.db)
//	runJob := NewRunJob(h.db)
//	test := NewTest(h.db)
//	runTest := NewRunTest(h.db)
//	runTestResult := NewRunTestResult(h.db)
//	incident := NewIncident(h.db)
//	switch reqBody.Method {
//	case "Plan.List":
//		if reqBody.HasParams() {
//			return invalidParamsResponse("This method is not support params", &reqBody.Id)
//		}
//		res := make([]*models.Plan, 0)
//		err := plan.List(nil, &res)
//		if err != nil {
//			return internalErrorResponse(err.Error(), &reqBody.Id)
//		}
//		return models.NewJsonRpcResponseResult(res, &reqBody.Id), http.StatusOK
//	case "Plan.EditName":
//		if !reqBody.HasParams() {
//			return invalidParamsResponse("This method need params", &reqBody.Id)
//		}
//		var args models.PlanEditNameArgs
//		err := json.Unmarshal(*reqBody.Params, &args)
//		if err != nil {
//			return invalidParamsResponse(err.Error(), &reqBody.Id)
//		}
//		var res models.JsonRpcResult
//		err = plan.EditName(args, &res)
//		if err != nil {
//			return internalErrorResponse(err.Error(), &reqBody.Id)
//		}
//		return models.NewJsonRpcResponseResult(res, &reqBody.Id), http.StatusOK
//	case "Plan.Get":
//		if !reqBody.HasParams() {
//			return invalidParamsResponse("This method need params", &reqBody.Id)
//		}
//		var args models.PlanGetArgs
//		err := json.Unmarshal(*reqBody.Params, &args)
//		if err != nil {
//			return invalidParamsResponse(err.Error(), &reqBody.Id)
//		}
//		var res *models.Plan
//		err = plan.Get(args, &res)
//		if err != nil {
//			return internalErrorResponse(err.Error(), &reqBody.Id)
//		}
//		return models.NewJsonRpcResponseResult(res, &reqBody.Id), http.StatusOK
//	case "RunPlan.List":
//		var args models.RunPlanListArgs
//		if reqBody.HasParams() {
//			err := json.Unmarshal(*reqBody.Params, &args)
//			if err != nil {
//				return invalidParamsResponse(err.Error(), &reqBody.Id)
//			}
//		}
//		var res models.RunPlanList
//		err := runPlan.List(args, &res)
//		if err != nil {
//			return internalErrorResponse(err.Error(), &reqBody.Id)
//		}
//		return models.NewJsonRpcResponseResult(res, &reqBody.Id), http.StatusOK
//	case "Job.List":
//		var args models.JobListArgs
//		if reqBody.HasParams() {
//			err := json.Unmarshal(*reqBody.Params, &args)
//			if err != nil {
//				return invalidParamsResponse(err.Error(), &reqBody.Id)
//			}
//		}
//		var res []*models.Job
//		err := job.List(args, &res)
//		if err != nil {
//			return internalErrorResponse(err.Error(), &reqBody.Id)
//		}
//		return models.NewJsonRpcResponseResult(res, &reqBody.Id), http.StatusOK
//	case "Job.Get":
//		if !reqBody.HasParams() {
//			return invalidParamsResponse("This method need params", &reqBody.Id)
//		}
//		var args models.JobGetArgs
//		err := json.Unmarshal(*reqBody.Params, &args)
//		if err != nil {
//			return invalidParamsResponse(err.Error(), &reqBody.Id)
//		}
//		var res *models.Job
//		err = job.Get(args, &res)
//		if err != nil {
//			return internalErrorResponse(err.Error(), &reqBody.Id)
//		}
//		return models.NewJsonRpcResponseResult(res, &reqBody.Id), http.StatusOK
//	case "Job.EditName":
//		if !reqBody.HasParams() {
//			return invalidParamsResponse("This method need params", &reqBody.Id)
//		}
//		var args models.JobEditNameArgs
//		err := json.Unmarshal(*reqBody.Params, &args)
//		if err != nil {
//			return invalidParamsResponse(err.Error(), &reqBody.Id)
//		}
//		var res models.JsonRpcResult
//		err = job.EditName(args, &res)
//		if err != nil {
//			return internalErrorResponse(err.Error(), &reqBody.Id)
//		}
//		return models.NewJsonRpcResponseResult(res, &reqBody.Id), http.StatusOK
//	case "RunJob.List":
//		var args models.RunJobListArgs
//		if reqBody.HasParams() {
//			err := json.Unmarshal(*reqBody.Params, &args)
//			if err != nil {
//				return invalidParamsResponse(err.Error(), &reqBody.Id)
//			}
//		}
//		var res models.RunJobList
//		err := runJob.List(args, &res)
//		if err != nil {
//			return internalErrorResponse(err.Error(), &reqBody.Id)
//		}
//		return models.NewJsonRpcResponseResult(res, &reqBody.Id), http.StatusOK
//	case "Test.Get":
//		if !reqBody.HasParams() {
//			return invalidParamsResponse("This method need params", &reqBody.Id)
//		}
//		var args models.TestGetArgs
//		err := json.Unmarshal(*reqBody.Params, &args)
//		if err != nil {
//			return invalidParamsResponse(err.Error(), &reqBody.Id)
//		}
//		var res *models.Test
//		err = test.Get(args, &res)
//		if err != nil {
//			return internalErrorResponse(err.Error(), &reqBody.Id)
//		}
//		return models.NewJsonRpcResponseResult(res, &reqBody.Id), http.StatusOK
//	case "RunTest.List":
//		var args models.RunTestListArgs
//		if reqBody.HasParams() {
//			err := json.Unmarshal(*reqBody.Params, &args)
//			if err != nil {
//				return invalidParamsResponse(err.Error(), &reqBody.Id)
//			}
//		}
//		var res models.RunTestList
//		err := runTest.List(args, &res)
//		if err != nil {
//			return internalErrorResponse(err.Error(), &reqBody.Id)
//		}
//		return models.NewJsonRpcResponseResult(res, &reqBody.Id), http.StatusOK
//	case "RunTestResult.List":
//		var args models.RunTestResultListArgs
//		if reqBody.HasParams() {
//			err := json.Unmarshal(*reqBody.Params, &args)
//			if err != nil {
//				return invalidParamsResponse(err.Error(), &reqBody.Id)
//			}
//		}
//		var res []*models.RunTestResult
//		err := runTestResult.List(args, &res)
//		if err != nil {
//			return internalErrorResponse(err.Error(), &reqBody.Id)
//		}
//		return models.NewJsonRpcResponseResult(res, &reqBody.Id), http.StatusOK
//	case "Incident.List":
//		var args models.IncidentListArgs
//		if reqBody.HasParams() {
//			err := json.Unmarshal(*reqBody.Params, &args)
//			if err != nil {
//				return invalidParamsResponse(err.Error(), &reqBody.Id)
//			}
//		}
//		var res []*models.IncidentResp
//		err := incident.List(args, &res)
//		if err != nil {
//			return internalErrorResponse(err.Error(), &reqBody.Id)
//		}
//		return models.NewJsonRpcResponseResult(res, &reqBody.Id), http.StatusOK
//	case "RunTestResult.StatusList":
//		var args models.RunTestResultStatusListArgs
//		if reqBody.HasParams() {
//			err := json.Unmarshal(*reqBody.Params, &args)
//			if err != nil {
//				return invalidParamsResponse(err.Error(), &reqBody.Id)
//			}
//		}
//		var res []*models.RunTestResultStatus
//		err := runTestResult.StatusList(args, &res)
//		if err != nil {
//			return internalErrorResponse(err.Error(), &reqBody.Id)
//		}
//		return models.NewJsonRpcResponseResult(res, &reqBody.Id), http.StatusOK
//	}
//
//	return models.NewJsonRpcResponseError(
//		models.JsonRpcErrorCodeMethodNotFound,
//		fmt.Sprintf("Unknown method '%s'", reqBody.Method),
//		nil,
//		nil),
//		http.StatusNotFound
//}
//
//func invalidParamsResponse(m string, id *uint) (*models.JsonRpcResponse, int) {
//	return models.NewJsonRpcResponseError(models.JsonRpcErrorCodeInvalidParams, m, nil, id), http.StatusBadRequest
//}
//
//func internalErrorResponse(m string, id *uint) (*models.JsonRpcResponse, int) {
//	return models.NewJsonRpcResponseError(models.JsonRpcErrorCodeInternalError, m, nil, id), http.StatusInternalServerError
//}
