package herd

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
)

//go:embed internal/migrations/*.sql
var migrationFS embed.FS

// Herd is used to migrate a database using a set of migrations.
//
// Migrations are applied in order of version, starting from the lowest to highest. Versions do not
// need to be consecutive. This allows version numbers to be based on dates.
type Herd struct {
	system *migrator
	user   *migrator
}

// New creates a new Herd instance with the given migrations.
func New(migrations []Migration) (*Herd, error) {
	systemMigrations, err := CollectFileMigrationsFromFS(migrationFS)
	if err != nil {
		return nil, fmt.Errorf("failed to collect system migrations: %w", err)
	}

	systemHelper := newMigrationHelper()

	system, err := newMigrator(systemMigrations)
	if err != nil {
		panic(err)
	}

	user, err := newMigrator(migrations)
	if err != nil {
		return nil, err
	}

	return &Herd{
		system: system,
		user:   user,
	}, nil
}

// Migrate migrates the database using the configured migrations.
func (h *Herd) Migrate(ctx context.Context, db *sql.DB) (*Result, error) {
	if _, err := h.system.migrate(ctx, db, false /*dryRun*/); err != nil {
		return nil, fmt.Errorf("failed to execute system migrations: %w", err)
	}

	result, err := h.user.migrate(ctx, db, false /*dryRun*/)
	if err != nil {
		return nil, fmt.Errorf("failed to execute user migrations: %w", err)
	}

	return result, nil
}
