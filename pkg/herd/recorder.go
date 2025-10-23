package herd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	tableNameSystem = "herd_system_migrations"
	tableNameUser   = "herd_user_migrations"
)

const (
	systemExistsQuery = "SELECT 1 FROM pg_tables WHERE schemaname='public' AND tablename='pg_type';"

	systemVersionQuery = "SELECT migration_version FROM herd_system_migrations " +
		"ORDER BY migration_version DESC LIMIT 1;"
	userVersionQuery = "SELECT migration_version FROM herd_user_migrations " +
		"ORDER BY migration_version DESC LIMIT 1;"

	systemRecordQuery = "INSERT INTO herd_system_migrations " +
		"(migration_version, migrated_at, code_version, code_revision) " +
		"VALUES($1, $2, $3, $4);"
	userRecordQuery = "INSERT INTO herd_user_migrations " +
		"(migration_version, migrated_at, code_version, code_revision, herd_version) " +
		"VALUES($1, $2, $3, $4, $5, $6);"
)

// recorder provides a table agnostic set of helpers for recording applied migrations.
type recorder struct {
	nowFunc      func() time.Time
	tableName    string
	codeVersion  string
	codeRevision string
}

// newRecorder creates a new recorder.
func newRecorder(
	nowFunc func() time.Time,
	tableName string,
	codeVersion string,
	codeRevision string,
) *recorder {
	return &recorder{
		nowFunc:      nowFunc,
		tableName:    tableName,
		codeVersion:  codeVersion,
		codeRevision: codeRevision,
	}
}

// GetCurrentVersion gets the current migration version, or 0 if no migrations were previously
// applied.
func (r *recorder) GetCurrentVersion(ctx context.Context, tx *sql.Tx) (int64, error) {
	exists, err := r.checkSystemExists(ctx, tx)
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, nil
	}

	query, err := r.getCurrentVersionQuery()
	if err != nil {
		return 0, err
	}

	row := tx.QueryRowContext(ctx, query)

	var version int64
	if err := row.Scan(&version); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("failed to scan current migration version: %w", err)
	}

	return version, nil
}

// checkSystemExists checks that the table for tracking system migrations exists before attempting
// to query it.
//
// If a query uses it before it exists, the transaction gets aborted and prevents application of
// migrations.
func (r *recorder) checkSystemExists(ctx context.Context, tx *sql.Tx) (bool, error) {
	if r.tableName != tableNameSystem {
		return true, nil
	}

	row := tx.QueryRowContext(ctx, systemExistsQuery)

	var exists int64
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("failed to test if system table exists: %w", err)
	}

	return true, nil
}

func (r *recorder) getCurrentVersionQuery() (string, error) {
	switch r.tableName {
	case tableNameSystem:
		return systemVersionQuery, nil

	case tableNameUser:
		return userVersionQuery, nil

	default:
		return "", fmt.Errorf("internal error: unexpected table name: %s", r.tableName)
	}
}

// RecordMigrations records the application of a migration.
func (r *recorder) RecordMigration(
	ctx context.Context,
	tx *sql.Tx,
	migrationVersion int64,
	herdVersion int64,
) error {
	query, err := r.recordMigrationQuery()
	if err != nil {
		return err
	}

	args, err := r.recordMigrationArgs(migrationVersion, herdVersion)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	return err
}

func (r *recorder) recordMigrationQuery() (string, error) {
	switch r.tableName {
	case tableNameSystem:
		return systemRecordQuery, nil

	case tableNameUser:
		return userRecordQuery, nil

	default:
		return "", fmt.Errorf("internal error: unexpected table name: %s", r.tableName)
	}
}

func (r *recorder) recordMigrationArgs(
	migrationVersion int64,
	herdVersion int64,
) ([]any, error) {
	now := r.nowFunc().UTC().Truncate(time.Second)

	switch r.tableName {
	case tableNameSystem:
		return []any{migrationVersion, now, r.codeVersion, r.codeRevision}, nil

	case tableNameUser:
		return []any{migrationVersion, now, r.codeVersion, r.codeRevision, herdVersion}, nil

	default:
		return nil, fmt.Errorf("internal error: unexpected table name: %s", r.tableName)
	}
}
