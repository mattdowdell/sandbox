package sqlx

import (
	"context"
	"database/sql"
	"log/slog"
)

// ...
type TxOption interface {
	apply(*sql.TxOptions)
}

// ...
type CommitFn func() error

// ...
type RollbackFn func(*slog.Logger)

// ...
func (d *DB) BeginTx(context.Context, ...TxOption) (Connection, CommitFn, RollbackFn, error) {
	return nil, nil, nil, nil
}
