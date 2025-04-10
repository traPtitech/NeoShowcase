package apiserver

import (
	"github.com/friendsofgo/errors"
)

type ErrorType int

const (
	ErrorTypeBadRequest ErrorType = iota
	ErrorTypeNotFound
	ErrorTypeAlreadyExists
	ErrorTypeForbidden
)

type customError struct {
	error
	typ ErrorType
}

func newError(typ ErrorType, message string, err error) error {
	if err == nil {
		return customError{error: errors.New(message), typ: typ}
	}
	return customError{error: errors.Wrap(err, message), typ: typ}
}

func DecomposeError(err error) (underlying error, typ ErrorType, ok bool) { //nolint:staticcheck
	var cErr customError
	if errors.As(err, &cErr) {
		return cErr.error, cErr.typ, true
	}
	return
}
