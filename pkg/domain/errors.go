package domain

import (
	"github.com/samber/oops"
)

// ErrorType is the kind of failure attached to an error by the layer that can
// tell what went wrong. It travels with the error as the oops error code, and the
// transport boundary maps it to a Connect code. An error without one is Internal
// by definition.
type ErrorType string

const (
	ErrorTypeBadRequest         ErrorType = "bad_request"
	ErrorTypeNotFound           ErrorType = "not_found"
	ErrorTypeAlreadyExists      ErrorType = "already_exists"
	ErrorTypeForbidden          ErrorType = "forbidden"
	ErrorTypeFailedPrecondition ErrorType = "failed_precondition"
)

// NewError tags err with a kind of failure and with the message the client is
// allowed to see. err may be nil.
//
// message is attached as the oops "public" message: it is the only thing the
// boundary returns to the client. The wrap chain underneath stays internal and
// reaches the single boundary log instead.
func NewError(typ ErrorType, message string, err error) error {
	b := oops.Code(string(typ)).Public(message)
	if err == nil {
		return b.New(message)
	}
	return b.Wrapf(err, "%s", message)
}

// DecomposeError reports the kind of failure attached to err, along with the
// client-facing message. ok is false when no layer of err carries a kind.
//
// It reads the classification from the outermost layer that carries one, not from
// oopsErr.Code()/Public() directly: those resolve to the *deepest* value in the
// chain, so a classified error wrapped as the cause of another would otherwise
// hijack the tag and public message. Walking Layers() (outermost to innermost)
// keeps the classification the top-most NewError set.
func DecomposeError(err error) (publicMessage string, typ ErrorType, ok bool) {
	oopsErr, found := oops.AsOops(err)
	if !found {
		return "", "", false
	}
	for _, layer := range oopsErr.Layers() {
		code, isString := layer.Code.(string)
		if isString && code != "" {
			return layer.Public, ErrorType(code), true
		}
	}
	return "", "", false
}
