package repository

import (
	"database/sql"
	"errors"

	"github.com/samber/oops"
)

var (
	ErrNotFound = oops.New("not found")
)

func isNoRowsErr(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
