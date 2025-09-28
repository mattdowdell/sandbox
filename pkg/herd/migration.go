package herd

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
)

// Executable contains the sql ExecContext method.
type Executable interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

// Migration must be implemented for each migration.
//
// If migrations are exclusively SQL files, FileMigration can and should be used. However, custom
// implementations may be required for complex data manipulation. It is strongly recommended that
// all migrations are deterministic as replaying a completed migration is not supported.
type Migration interface {
	// Name returns the migration name.
	Name() string

	// Version returns the migration version.
	Version() int

	// Migrate executes the migration.
	Migrate(context.Context, Executable) error
}

// FileMigration is an implementation of Migration and can represent a single SQL file containing a
// migration.
type FileMigration struct {
	name     string
	version  int
	contents string
}

// CollectFileMigrationsFromFS walks the given filesystem looking for files containing migrations.
//
// For each filename with a ".sql" extension, a migration is created using NewFileMigrationFromFS.
func CollectFileMigrationsFromFS(filesystem fs.FS) ([]Migration, error) {
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

// NewFileMigrationFromFS wraps NewFileMigrationFromFilename, providing the file contents read from
// the filesystem at the given path.
func NewFileMigrationFromFS(filesystem fs.FS, path string) (*FileMigration, error) {
	contents, err := fs.ReadFile(filesystem, path)
	if err != nil {
		return nil, err
	}

	return NewFileMigrationFromFilename(filepath.Base(path), string(contents))
}

// NewFileMigrationFromFilename wraps NewFileMigration, deriving the name and version from the given
// filename.
//
// The filename is expected to be in the format "{version}_{name}.sql" where the version is an
// integer. The version may have leading zeros to help with filename sorting. For example, "0001"
// will be parsed as 1.
func NewFileMigrationFromFilename(filename, contents string) (*FileMigration, error) {
	trimmed, ok := strings.CutSuffix(filename, ".sql")
	if !ok {
		return nil, fmt.Errorf(
			"unexpected extension for filename %q: want: .sql, found: %s",
			filename,
			filepath.Ext(filename),
		)
	}

	version, name, ok := strings.Cut(trimmed, "_")
	if !ok {
		return nil, fmt.Errorf("unable to extract version and name from filename: %q", filename)
	}

	v, err := strconv.Atoi(version)
	if err != nil {
		return nil, fmt.Errorf("unable to parse version to int for filename %q: %w", filename, err)
	}

	return NewFileMigration(name, v, contents)
}

// NewFileMigration creates a new FileMigration.
func NewFileMigration(name string, version int, contents string) (*FileMigration, error) {
	if version < 0 {
		return nil, fmt.Errorf("version must be > 0, found: %d", version)
	}

	return &FileMigration{
		name:     name,
		version:  version,
		contents: contents,
	}, nil
}

// Name returns the migration name.
func (m *FileMigration) Name() string {
	return m.name
}

// Version returns the migration version.
func (m *FileMigration) Version() int {
	return m.version
}

// Migrate executes the migration.
func (m *FileMigration) Migrate(ctx context.Context, tx Executable) error {
	_, err := tx.ExecContext(ctx, m.contents)
	return err
}
