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
	limit int64
}

// ...
type migrator struct {
	helper     *migrationHelper
	migrations []Migration

}

// ...
func newMigrator(helper *migrationHelper, migrations []Migration) (*migrator, error) {
	if len(migrations) == 0 {
		return nil, fmt.Errorf("at least 1 migration must be present")
	}

	slices.SortFunc(migrations, func(a, b Migration) int {
		return cmp.Compare(a.Version(), b.Version())
	})

	for i := range migrations {
		version := migrations[i].Version()
		if version <= 0 {
			return nil, fmt.Errorf("migration versions must be > 0, found: %d", version)
		}

		if i == 0 {
			continue
		}

		if migrations[i - 1].Version() == version {
			return nil, fmt.Errorf("duplicate migration version found: %d", version)
		}
	}

	return &migrator{
		helper: helper,
		migrations: migrations,
	}, nil
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

	currentVersion, err := m.helper.getCurrentVersion(ctx, tx)
	if err != nil {
		return nil, err
	}

	pending := slices.DeleteFunc(m.migrations, func(m Migration) bool {
		return m.Version() <= currentVersion || opts.limit >
	})

	for _, migration := range pending {
		if err := migration.Migrate(ctx, tx); err != nil {
			return nil, fmt.Errorf("migration %s %d failed: %w", migration.Name(), migration.Version(), err)
		}
	}

	if !opts.dryRun {
		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("failed to commit transaction: %w", err)
		}
	}

	return &Result{
		// TODO: return previous + new version
	}, nil
}
