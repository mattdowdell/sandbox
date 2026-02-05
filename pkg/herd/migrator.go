package herd

import (
	"cmp"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
)

// Result contains the migration version before and after the pending migrations were applied.
type Result struct {
	// Before is the migration version before applying any migrations. A value of 0 indicates that
	// no migrations had been applied previously.
	Before int64

	// After is the migration version after applying any migrations. After will equal Before when no
	// pending migrations were found.
	After int64
}

// migrator applies system or user-defined migrations to the database.
type migrator struct {
	migrations []Migration
	recorder   *recorder
}

// newMigrator creates a new migrator.
func newMigrator(migrations []Migration, rec *recorder) (*migrator, error) {
	slices.SortFunc(migrations, func(a, b Migration) int {
		return cmp.Compare(a.Version(), b.Version())
	})

	for i, m := range migrations {
		if m.Version() < 1 {
			return nil, fmt.Errorf("migration version must be > 0, found: %d", m.Version())
		}

		// rely on sorting to detect duplicates
		if i > 0 && m.Version() == migrations[i-1].Version() {
			return nil, fmt.Errorf("duplicate migration version found: %d", m.Version())
		}
	}

	return &migrator{
		migrations: migrations,
		recorder:   rec,
	}, nil
}

// migrate applies pending migrations to the database.
func (m *migrator) Migrate(
	ctx context.Context,
	db *sql.DB,
	herdVersion int64,
) (result *Result, err error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// If we call Commit before Rollback, an atomic bool is flipped, Rollback will always return
	// sql.ErrTxDone and the driver's rollback implementation is never called. This occurs
	// regardless of whether Commit failed.
	defer func() {
		if err2 := tx.Rollback(); err2 != nil && !errors.Is(err2, sql.ErrTxDone) {
			err = errors.Join(err, fmt.Errorf("failed to rollback transaction: %w", err2))
		}
	}()

	before, err := m.recorder.GetCurrentVersion(ctx, tx)
	if err != nil {
		return nil, err
	}

	pending := slices.DeleteFunc(slices.Clone(m.migrations), func(m Migration) bool {
		return m.Version() <= before
	})

	for _, migration := range pending {
		if err := migration.Migrate(ctx, tx); err != nil {
			return nil, fmt.Errorf("failed to apply migration %d: %w", migration.Version(), err)
		}

		if err := m.recorder.RecordMigration(
			ctx,
			tx,
			migration.Version(),
			herdVersion,
		); err != nil {
			return nil, fmt.Errorf("failed to record migration %d: %w", migration.Version(), err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	after := before
	if len(pending) > 0 {
		after = pending[len(pending)-1].Version()
	}

	return &Result{
		Before: before,
		After:  after,
	}, nil
}
