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
func New(migrations []Migration, options ...Option) (*Herd, error) {
	opts := defaultHerdOpts()
	for _, option := range options {
		option.apply(opts)
	}

	version, revision, err := opts.codeInfo()
	if err != nil {
		return nil, err
	}

	systemMigrations, err := CollectFileMigrations(migrationFS)
	if err != nil {
		return nil, fmt.Errorf("failed to collect system migrations: %w", err)
	}

	systemRecorder := newRecorder(opts.nowFunc, tableNameSystem, version, revision)
	userRecorder := newRecorder(opts.nowFunc, tableNameUser, version, revision)

	system, err := newMigrator(systemMigrations, systemRecorder)
	if err != nil {
		return nil, fmt.Errorf("failed to create system migrator: %w", err)
	}

	user, err := newMigrator(migrations, userRecorder)
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
	systemResult, err := h.system.Migrate(ctx, db, 0 /*herdVersion*/)
	if err != nil {
		return nil, fmt.Errorf("failed to execute system migrations: %w", err)
	}

	userResult, err := h.user.Migrate(ctx, db, systemResult.After)
	if err != nil {
		return nil, fmt.Errorf("failed to execute user migrations: %w", err)
	}

	return userResult, nil
}
