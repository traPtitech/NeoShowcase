package grpc

import (
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
)

func TestLogError_UnclassifiedPlainError(t *testing.T) {
	// A plain error has no oops in its chain; logging the *connect.Error as is
	// shows only "internal server error" and hides the cause.
	err := handleUseCaseError(errors.New("not found"))

	assert.Equal(t, connect.CodeInternal, connect.CodeOf(err))
	assert.Equal(t, "not found", logError(err))
}
