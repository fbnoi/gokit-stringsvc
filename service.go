package main

import (
	"errors"
	"strings"
)

// ErrEmpty empty string error
var ErrEmpty = errors.New("Empty string")

// Service

// IStringService StringService
type IStringService interface {
	UpperCase(string) (string, error)
	Count(string) int
}

// StringService IStringService implement
type StringService struct {
}

// UpperCase  return upper case of string
func (svc StringService) UpperCase(str string) (string, error) {
	if str == "" {
		return str, ErrEmpty
	}
	return strings.ToUpper(str), nil
}

// Count return length of given string
func (svc StringService) Count(str string) int {
	return len(str)
}
