package herd

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
)

const (
	systemTable = "herd_system_migrations"
	userTable   = "herd_user_migrations"
)

// ...
type Clock interface {
	UTCNow() time.Time
}

// ...
type UUIDGenerator interface {
	NewV7() (uuid.UUID, error)
}

// ...
type migrationHelper struct {
	clock Clock
	uuidgen UUIDGenerator
	table string
	codeVersion string
}

// ...
func newMigrationHelper(
	clock Clock,
	uuidgen UUIDGenerator,
	table string,
	codeVersion string,
) *migrationHelper {
	return &migrationHelper
}

// ...
func (h *migrationHelper) getCurrentVersion(ctx context.Context, db DB) (int64, error) {
	query := getCurrentVersionQuery(h.table)
	row := db.QueryRowContext(ctx, query)

	var version int64
	if err := row.Scan(&version); err != nil {
		return 0, fmt.Errorf("failed to scan current version: %w", err)
	}

	return version, nil
}

// ...
func (h *migrationHelper) recordMigration(
	ctx context.Context,
	db DB,
	migrationVersion int64,
) error {
	query := recordMigrationQuery(h.table)

	migratedAt := h.clock.UTCNow() // TODO: truncate

	id, err := h.uuidgen.NewV7()
	if err != nil {
		return fmt.Errorf("failed to generated id to record migration: %w", err)
	}

	if _, err := db.ExecContext(
		ctx,
		query,
		id,
		migratedAt,
		migrationVersion,
		h.codeVersion,
	); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return nil
}

// ...
func getCurrentVersionQuery(table string) string {
	return fmt.Sprintf(
		"SELECT migration_version FROM %s ORDER BY migration_version LIMIT 1;",
		table,
	)
}

// ...
func recordMigrationQuery(table string) string {
	return fmt.Sprintf(
		"INSERT INTO %s (id, migrated_at, migration_version, code_version) VALUES($1, $2, $3, $4);",
		table,
	)
}
