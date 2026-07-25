package repository

import (
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
)

func isNoRowsErr(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
