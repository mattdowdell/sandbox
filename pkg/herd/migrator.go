package herd

import (
	"cmp"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ...
type migrateOpts struct {
	dryRun bool
}

// ...
type migrator struct {
	migrations []Migration
}

// ...
func newMigrator(migrations []Migration) *migrator {
	// TODO: require 1 migration to be present?

	slices.SortFunc(migrations, func(a, b Migration) int {
		return cmp.Compare(a.Version(), b.Version())
	})

	return &migrator{
		migrations: migrations,
	}
}

// ...
func (m *migrator) migrate(ctx context.Context, db *sql.DB, opts *migrateOpts) (*Result, error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	// TODO: return this error
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			slog.ErrorContext(ctx, "failed to rollback transaction", slogx.Err(err))
		}
	}()

	// TODO: filter to just pending migrations

	for _, migration := range m.migrations {
		if err := migration.Migrate(ctx, tx); err != nil {
			return nil, fmt.Errorf("migration %s %d failed: %w", migration.Name(), migration.Version(), err)
		}
	}

	// TODO: support dry run
	if !opts.dryRun {
		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("failed to commit transaction: %w", err)
		}
	}

	return &Result{
		// TODO: return previous + new version
	}, nil
}
