package usecase

import (
	"errors"
	"fmt"
)

type ErrorType int

const (
	ErrorTypeBadRequest ErrorType = iota
	ErrorTypeNotFound
	ErrorTypeAlreadyExists
)

type Error struct {
	error
	Type ErrorType
}

func newError(typ ErrorType, message string, err error) error {
	if err == nil {
		return &Error{error: errors.New(message), Type: typ}
	}
	return &Error{error: fmt.Errorf("%s: %w", message, err), Type: typ}
}
