package herd

import (
	"cmp"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/pkg/herd/internal/herdconv"
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
	tracer     trace.Tracer
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
		tracer:     otelx.Tracer(),
	}, nil
}

// migrate applies pending migrations to the database.
func (m *migrator) Migrate(
	ctx context.Context,
	db *sql.DB,
	herdVersion int64,
) (result *Result, err error) {
	ctx, span := m.tracer.Start(ctx, migratorSpanName(m.recorder.TableName()), trace.WithAttributes(
		herdconv.HerdTableName(m.recorder.TableName()),
	))
	defer span.End()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		span.SetStatus(codes.Error, "failed to begin transaction")
		span.RecordError(err)

		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// If we call Commit before Rollback, an atomic bool is flipped, Rollback will always return
	// sql.ErrTxDone and the driver's rollback implementation is never called. This occurs
	// regardless of whether Commit failed.
	defer func() {
		if err2 := tx.Rollback(); err2 != nil && !errors.Is(err2, sql.ErrTxDone) {
			span.SetStatus(codes.Error, "failed to rollback transaction")
			span.RecordError(err2)

			err = errors.Join(err, fmt.Errorf("failed to rollback transaction: %w", err2))
		}
	}()

	before, err := m.recorder.CurrentVersion(ctx, tx)
	if err != nil {
		span.SetStatus(codes.Error, "failed to get current version")
		span.RecordError(err)

		return nil, err
	}

	span.SetAttributes(herdconv.HerdVersionBefore(int(before)))

	pending := slices.DeleteFunc(slices.Clone(m.migrations), func(m Migration) bool {
		return m.Version() <= before
	})

	for _, migration := range pending {
		if err := m.applyMigration(ctx, tx, migration, herdVersion); err != nil {
			span.SetStatus(codes.Error, "failed to apply migration")
			span.RecordError(err)

			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		span.SetStatus(codes.Error, "failed to commit transaction")
		span.RecordError(err)

		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	after := before
	if len(pending) > 0 {
		after = pending[len(pending)-1].Version()
	}

	span.SetAttributes(herdconv.HerdVersionAfter(int(after)))

	return &Result{
		Before: before,
		After:  after,
	}, nil
}

func (m *migrator) applyMigration(
	ctx context.Context,
	tx *sql.Tx,
	migration Migration,
	herdVersion int64,
) error {
	ctx, span := m.tracer.Start(
		ctx,
		fmt.Sprintf("Apply Migration %d", migration.Version()),
		trace.WithAttributes(
			herdconv.HerdTableName(m.recorder.TableName()),
			herdconv.HerdVersionAfter(int(migration.Version())),
		),
	)
	defer span.End()
	if err := migration.Migrate(ctx, tx); err != nil {
		span.SetStatus(codes.Error, "failed to execute migration")
		span.RecordError(err)

		return fmt.Errorf("failed to execute migration %d: %w", migration.Version(), err)
	}

	if err := m.recorder.RecordMigration(
		ctx,
		tx,
		migration.Version(),
		herdVersion,
	); err != nil {
		span.SetStatus(codes.Error, "failed to record migration")
		span.RecordError(err)

		return fmt.Errorf("failed to record migration %d: %w", migration.Version(), err)
	}

	return nil
}

func migratorSpanName(table string) string {
	switch table {
	case TableNameSystem:
		return "Apply System Migrations"

	case TableNameUser:
		return "Apply User Migrations"

	default:
		panic("unknown table name: " + table)
	}
}
