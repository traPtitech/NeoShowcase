package repository

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("NotFound")

func isNoRowsErr(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
