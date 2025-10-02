package herd

import (
	"context"
	"database/sql"
)

// DB contains the methods from [sql.DB] that migrations can use.
type DB interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

// Migration implementations apply schema or data changes to a database.
type Migration interface {
	// Version returns the migration version. It must be greater than 0 and unique across all
	// migrations.
	Version() int64

	// Migrate applies the schema or data changes for the migration.
	Migrate(context.Context, DB) error
}