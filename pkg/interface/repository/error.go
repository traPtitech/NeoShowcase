package repository

import (
	"database/sql"

	"github.com/friendsofgo/errors"
)

var (
	ErrNotFound = errors.New("not found")
)

func isNoRowsErr(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
