package herd_test

import (
	"database/sql"
	"io/fs"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newMockDB is used to create mocks for database queries.
func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	return db, mock
}

// writeFile creates an empty file at the given path.
func writeFile(t *testing.T, path string) {
	t.Helper()

	err := os.WriteFile(path, nil, 0o600)
	require.NoErrorf(t, err, "failed to create file: %s", path)
}

// openRootFS creates a fs.FS for the given directory.
func openRootFS(t *testing.T, dir string) fs.FS {
	t.Helper()

	root, err := os.OpenRoot(dir)
	require.NoError(t, err)

	return root.FS()
}
