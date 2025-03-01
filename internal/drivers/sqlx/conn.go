package sqlx

import (
	"context"
	"database/sql"
)

// ...
type Connection interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)

	// ...
	Exec(string, ...any) (sql.Result, error)

	// ...
	Query(string, ...any) (*sql.Rows, error)
}

// ...
func (d *DB) Conn() Connection {
	return d.db
}
