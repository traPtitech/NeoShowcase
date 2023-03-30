package repository

import (
	"database/sql"

	"github.com/friendsofgo/errors"

	"github.com/go-sql-driver/mysql"
)

const (
	mysqlErrorDuplicateEntry = 1062
)

var (
	ErrNotFound  = errors.New("not found")
	ErrDuplicate = errors.New("duplicate record exists")
)

func isNoRowsErr(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func isDuplicateKeyErr(err error) bool {
	var myErr *mysql.MySQLError
	return errors.As(err, &myErr) && myErr.Number == mysqlErrorDuplicateEntry
}
