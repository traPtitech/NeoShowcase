package repository

import (
	"database/sql"
	"errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

// notFound reports that the entity named by what does not exist, already
// classified so that the boundary answers NotFound without the caller doing
// anything.
func notFound(what string) error {
	return domain.NewError(domain.ErrorTypeNotFound, what+" not found", nil)
}

func isNoRowsErr(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
