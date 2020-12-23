package main

import (
	"context"
	"encoding/json"
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

func decodeCountRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
