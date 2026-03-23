package datastore

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// isPgErr tests if the given error is a PostgreSQL error with the given code.
//
// Constants for all PostgreSQL error codes can be found in [pgerrcode].
//
// [pgerrcode]: github.com/jackc/pgerrcode
func isPgErr(err error, code string) bool {
	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	return ok && pgErr.Code == code
}
