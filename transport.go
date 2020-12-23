package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Transport

type upperCaseRequest struct {
	Str string `json:"str"`
}

type upperCaseResponse struct {
	Str string `json:"str"`
	Err string `json:"err,omitempty"`
}

type countRequest struct {
	Str string `json:"str"`
}

type countResponse struct {
	Count int    `json:"count"`
	Err   string `json:"err,omitempty"`
}

func decodeUpperCaseRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request upperCaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeUpperCaseResponse(ctx context.Context, w *http.Response) (interface{}, error) {
	var response upperCaseResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	return response, err
}

func decodeCountRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// encode interface{} (local response) to http response, used to response to http request
// response : local => http
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	err := json.NewEncoder(w).Encode(response)
	return err
}

// encode local request (local request) to http request for client to request for http
// request: local => http
func encodeRequest(ctx context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		log.Print(err)
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}
