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

type customError struct {
	error
	typ ErrorType
}

func newError(typ ErrorType, message string, err error) error {
	if err == nil {
		return customError{error: errors.New(message), typ: typ}
	}
	return customError{error: fmt.Errorf("%s: %w", message, err), typ: typ}
}

func GetErrorType(err error) (typ ErrorType, ok bool) {
	var cErr customError
	if errors.As(err, &cErr) {
		return cErr.typ, true
	}
	return
}
