package main

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoint
// type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)

func makeUppercaseEndPoint(svc IStringService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(upperCaseRequest)
		str, err := svc.UpperCase(request.Str)
		if err != nil {
			return upperCaseResponse{Str: str, Err: err.Error()}, err
		}
		return upperCaseResponse{Str: str, Err: ""}, nil
	}
}

func makeCountEndPoint(svc IStringService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(countRequest)
		count := svc.Count(request.Str)
		return countResponse{Count: count, Err: ""}, nil
	}
}
