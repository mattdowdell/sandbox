package herd

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
)

// Non-allocation compile-time check for interface compliance.
var _ Migration = (*FileMigration)(nil)

// FileMigration uses the contents of a SQL file as a migration.
type FileMigration struct {
	version  int64
	contents string
}

// CollectFileMigrations walks the filesystem, searching for files with a ".sql" extension. For each
// found file, NewFileMigrationFromFS is called to create the migration.
func CollectFileMigrations(filesystem fs.FS) ([]Migration, error) {
	var migrations []Migration

	if err := fs.WalkDir(filesystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".sql" {
			return nil
		}

		migration, err := NewFileMigrationFromFS(filesystem, path)
		if err != nil {
			return err
		}

		migrations = append(migrations, migration)
		return nil
	}); err != nil {
		return nil, err
	}

	return migrations, nil
}

// NewFileMigrationFromFS reads the file at the path from the filesystem. The last element of the
// path is taken as the filename, and is passed with the contents to NewFileMigrationFromFilename.
func NewFileMigrationFromFS(filesystem fs.FS, path string) (*FileMigration, error) {
	contents, err := fs.ReadFile(filesystem, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read migration from filesystem: %w", err)
	}

	return NewFileMigrationFromFilename(filepath.Base(path), string(contents))
}

// NewFileMigrationFromFilename extracts the migration version from the given filename and calls
// NewFileMigration.
//
// The version must be at the start of the filename followed by a underscore ("_"). For example,
// 1_initial.sql would have a version of 1. Versions can be padded with 0s to make them easier to
// view in a file browser. 0001_initial.sql and 1_initial.sql would have the same version.
func NewFileMigrationFromFilename(filename, contents string) (*FileMigration, error) {
	version, _, ok := strings.Cut(filename, "_")
	if !ok {
		return nil, fmt.Errorf("could not extract version from filename: %s", filename)
	}

	parsed, err := strconv.ParseInt(version, 10 /*base*/, 64 /*bitSize*/)
	if err != nil {
		return nil, fmt.Errorf("could not parse version from filename %s: %w", filename, err)
	}

	return NewFileMigration(parsed, contents), nil
}

// NewFileMigration creates a new FileMigration.
func NewFileMigration(version int64, contents string) *FileMigration {
	return &FileMigration{
		version:  version,
		contents: contents,
	}
}

// Version returns the migration version.
func (m *FileMigration) Version() int64 {
	return m.version
}

// Migrate executes the SQL file contents against the database.
func (m *FileMigration) Migrate(ctx context.Context, db DB) error {
	_, err := db.ExecContext(ctx, m.contents)
	return err
}
