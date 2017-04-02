package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andrskom/jrpc2hh/models"
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"
)

type Caller interface {
	Call(reqBody *models.RequestBody) (interface{}, *models.Error)
}

type Handler struct {
	mu   sync.Mutex
	sMap map[string]Caller
}

func NewHandler() *Handler {
	return &Handler{sMap: make(map[string]Caller)}
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

func (h *Handler) getService(sN string) (Caller, *models.Error) {
	c, ok := h.sMap[sN]
	if !ok {
		return nil, models.NewError(
			models.ErrorCodeMethodNotFound,
			"Unknown service",
			map[string]string{"methodName": sN})
	}

	return c, nil
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
		rB, httpSt := h.doProcedure(jReq)
		models.JsonResponse(w, rB, httpSt)
		return
	}

	// If hasn't error, that it is batch request
	var jReqBatchSlice []json.RawMessage
	errReqBatch := json.Unmarshal(reqBodyBytes, &jReqBatchSlice)
	if errReqBatch == nil {
		jRespBatchSlice := struct {
			sync.Mutex
			Slice []*models.ResponseBody
		}{Slice: make([]*models.ResponseBody, 0)}
		wg := sync.WaitGroup{}
		for _, req := range jReqBatchSlice {
			wg.Add(1)
			go func(req json.RawMessage) {
				defer wg.Done()
				var jReqBatch *models.RequestBody
				err := json.Unmarshal(req, &jReqBatch)
				if err != nil {
					jRespBatchSlice.Lock()
					defer jRespBatchSlice.Unlock()
					jRespBatchSlice.Slice = append(
						jRespBatchSlice.Slice,
						models.NewResponseError(
							models.NewError(
								models.ErrorCodeInvalidRequest,
								err.Error(),
								nil),
							nil,
						))
				} else {
					rB, _ := h.doProcedure(jReqBatch)
					jRespBatchSlice.Lock()
					defer jRespBatchSlice.Unlock()
					jRespBatchSlice.Slice = append(
						jRespBatchSlice.Slice,
						rB)
				}
			}(req)
		}
		wg.Wait()
		models.JsonResponse(w, jRespBatchSlice.Slice, http.StatusOK)
		return
	}

	// Else that error
	models.JsonResponse(
		w,
		models.NewError(models.ErrorCodeParseError, "Can't parse request json to json rpc 2.0 struct", nil),
		http.StatusBadRequest)
	return
}

func (h *Handler) doProcedure(jReq *models.RequestBody) (*models.ResponseBody, int) {
	err := jReq.Validate()
	if err != nil {
		jErr := models.NewError(models.ErrorCodeParseError, err.Error(), nil)
		return models.NewResponseError(jErr, jReq.Id), http.StatusBadRequest
	}
	s, mErr := h.getService(jReq.GetService())
	if mErr != nil {
		return models.NewResponseError(mErr, jReq.Id), http.StatusNotFound
	}
	res, jErr := s.Call(jReq)
	if jErr != nil {
		return models.NewResponseError(jErr, jReq.Id), http.StatusInternalServerError
	}

	if res == nil {
		return models.NewResponseBody(nil, jReq.Id), http.StatusOK
	} else {
		resByte, err := json.Marshal(res)
		if err != nil {
			jErr := models.NewError(models.ErrorCodeInternalError, "Can't marshal response", err.Error())
			return models.NewResponseError(jErr, jReq.Id), http.StatusInternalServerError
		}
		jsonRes := json.RawMessage(resByte)
		return models.NewResponseBody(&jsonRes, jReq.Id), http.StatusOK
	}
}
