package domain

import (
	"errors"
	"testing"

	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
)

func TestDecomposeError(t *testing.T) {
	notFound := NewError(ErrorTypeNotFound, "application not found", nil)

	tests := []struct {
		name       string
		err        error
		wantOK     bool
		wantType   ErrorType
		wantPublic string
	}{
		{
			name:   "plain error is unclassified",
			err:    errors.New("boom"),
			wantOK: false,
		},
		{
			name:   "oops error without a kind is unclassified",
			err:    oops.Wrapf(errors.New("boom"), "doing something"),
			wantOK: false,
		},
		{
			// A repository not found passes through usecases that only add context.
			name:       "kind survives wraps without a kind",
			err:        oops.Wrapf(oops.Wrapf(notFound, "getting application"), "checking owner"),
			wantOK:     true,
			wantType:   ErrorTypeNotFound,
			wantPublic: "application not found",
		},
		{
			name:       "outermost kind wins",
			err:        NewError(ErrorTypeFailedPrecondition, "inconsistent", oops.Wrapf(notFound, "getting application")),
			wantOK:     true,
			wantType:   ErrorTypeFailedPrecondition,
			wantPublic: "inconsistent",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			public, typ, ok := DecomposeError(tt.err)
			assert.Equal(t, tt.wantOK, ok)
			assert.Equal(t, tt.wantType, typ)
			assert.Equal(t, tt.wantPublic, public)
		})
	}
}
