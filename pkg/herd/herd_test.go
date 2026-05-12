package herd_test

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"runtime/debug"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattdowdell/sandbox/pkg/herd"
)

const (
	testVersion  = "v0.0.1"
	testRevision = "abcdef0"

	systemMigration1 = `-- Initial migration for recording which migrations have been applied via herd.

-- Create a table for recording internal/system migrations. This is mostly a way to track what
-- schema version is in use.
CREATE TABLE herd_system_migrations (
	migration_version BIGINT PRIMARY KEY,
	migrated_at TIMESTAMPTZ (0) NOT NULL,
	code_version TEXT NOT NULL,
	code_revision TEXT NOT NULL
);

-- Create a table for recording user-defined migrations. Also include the herd_version which is the
-- highest herd_system_migrations.migration_version value.
CREATE TABLE herd_user_migrations (
	migration_version BIGINT PRIMARY KEY,
	migrated_at TIMESTAMPTZ (0) NOT NULL,
	code_version TEXT NOT NULL,
	code_revision TEXT NOT NULL,
	herd_version BIGINT NOT NULL
);
`
)

func Test_New_Success(t *testing.T) {
	// arrange
	migrations, err := herd.CollectFileMigrations(migrationFS)
	require.NoError(t, err)

	// act
	migrator, err := herd.New(migrations, herd.WithBuildInfoValues(testVersion, testRevision))

	// assert
	assert.NotNil(t, migrator)
	assert.NoError(t, err)
}

func Test_New_Error(t *testing.T) {
	tests := map[string]struct {
		migrations []herd.Migration
		options    []herd.Option
		want       string
	}{
		"nil buildinfo": {
			options: []herd.Option{
				herd.WithBuildInfo(nil),
			},
			want: "build info is unavailable",
		},
		"missing code version": {
			options: []herd.Option{
				herd.WithBuildInfo(&debug.BuildInfo{}),
			},
			want: "unable to extract code version from build info",
		},
		"missing code revision": {
			options: []herd.Option{
				herd.WithBuildInfo(&debug.BuildInfo{
					Main: debug.Module{
						Version: testVersion,
					},
				}),
			},
			want: "unable to extract code revision from build info",
		},
		"invalid migration version": {
			migrations: []herd.Migration{
				herd.NewFileMigration(0, ""),
			},
			options: []herd.Option{
				herd.WithBuildInfoValues(testVersion, testRevision),
			},
			want: "migration version must be > 0, found: 0",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange

			// act
			migrator, err := herd.New(tt.migrations, tt.options...)

			// assert
			assert.Nil(t, migrator)
			assert.EqualError(t, err, tt.want)
		})
	}
}

func Test_Herd_Migrate_Success(t *testing.T) {
	now := time.Now().UTC().Round(time.Second)

	tests := map[string]struct {
		db   func(*testing.T) *sql.DB
		want *herd.Result
	}{
		"pending system migrations": {
			db: func(t *testing.T) *sql.DB {
				t.Helper()

				db, mock := newMockDB(t)
				mock.ExpectBegin()
				mock.ExpectQuery(herd.SystemExistsQuery).
					WillReturnRows(sqlmock.NewRows(testExistsRows).AddRow(1))
				mock.ExpectQuery(herd.SystemVersionQuery).
					WillReturnRows(sqlmock.NewRows(testQueryRows))
				mock.ExpectExec(systemMigration1).WillReturnResult(driver.ResultNoRows)
				mock.ExpectExec(herd.SystemRecordQuery).
					WithArgs(1, now, testVersion, testRevision).
					WillReturnResult(driver.ResultNoRows)
				mock.ExpectCommit()

				mock.ExpectBegin()
				mock.ExpectQuery(herd.UserVersionQuery).
					WillReturnRows(sqlmock.NewRows(testQueryRows).AddRow(1))
				mock.ExpectCommit()

				return db
			},
			want: &herd.Result{
				Before: 1,
				After:  1,
			},
		},
		"pending user migration": {
			db: func(t *testing.T) *sql.DB {
				t.Helper()

				db, mock := newMockDB(t)
				mock.ExpectBegin()
				mock.ExpectQuery(herd.SystemExistsQuery).
					WillReturnRows(sqlmock.NewRows(testExistsRows).AddRow(1))
				mock.ExpectQuery(herd.SystemVersionQuery).
					WillReturnRows(sqlmock.NewRows(testQueryRows).AddRow(testHerdVersion))
				mock.ExpectCommit()

				mock.ExpectBegin()
				mock.ExpectQuery(herd.UserVersionQuery).
					WillReturnRows(sqlmock.NewRows(testQueryRows))
				mock.ExpectExec("-- placeholder").WillReturnResult(driver.ResultNoRows)
				mock.ExpectExec(herd.UserRecordQuery).
					WithArgs(1, now, testVersion, testRevision, testHerdVersion).
					WillReturnResult(driver.ResultNoRows)
				mock.ExpectCommit()

				return db
			},
			want: &herd.Result{
				Before: 0,
				After:  1,
			},
		},
		"no pending migrations": {
			db: func(t *testing.T) *sql.DB {
				t.Helper()

				db, mock := newMockDB(t)
				mock.ExpectBegin()
				mock.ExpectQuery(herd.SystemExistsQuery).
					WillReturnRows(sqlmock.NewRows(testExistsRows).AddRow(1))
				mock.ExpectQuery(herd.SystemVersionQuery).
					WillReturnRows(sqlmock.NewRows(testQueryRows).AddRow(testHerdVersion))
				mock.ExpectCommit()

				mock.ExpectBegin()
				mock.ExpectQuery(herd.UserVersionQuery).
					WillReturnRows(sqlmock.NewRows(testQueryRows).AddRow(testHerdVersion))
				mock.ExpectCommit()

				return db
			},
			want: &herd.Result{
				Before: 1,
				After:  1,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			migrations, err := herd.CollectFileMigrations(migrationFS)
			require.NoError(t, err)

			migrator, err := herd.New(
				migrations,
				herd.WithNowFunc(testNowFunc(now)),
				herd.WithBuildInfoValues(testVersion, testRevision),
			)
			require.NoError(t, err)

			db := tt.db(t)

			// act
			result, err := migrator.Migrate(t.Context(), db)

			// assert
			assert.Equal(t, tt.want, result)
			assert.NoError(t, err)
		})
	}
}

func Test_Herd_Migrate_Error(t *testing.T) {
	tests := map[string]struct {
		db   func(*testing.T) *sql.DB
		want string
	}{
		"system migrate error": {
			db: func(t *testing.T) *sql.DB {
				t.Helper()

				db, mock := newMockDB(t)
				mock.ExpectBegin().WillReturnError(errors.New("begin error"))

				return db
			},
			want: "failed to execute system migrations: failed to begin transaction: begin error",
		},
		"user migrate error": {
			db: func(t *testing.T) *sql.DB {
				t.Helper()

				db, mock := newMockDB(t)
				mock.ExpectBegin()
				mock.ExpectQuery(herd.SystemExistsQuery).
					WillReturnRows(sqlmock.NewRows(testExistsRows).AddRow(1))
				mock.ExpectQuery(herd.SystemVersionQuery).
					WillReturnRows(sqlmock.NewRows(testQueryRows).AddRow(testHerdVersion))
				mock.ExpectCommit()

				mock.ExpectBegin().WillReturnError(errors.New("begin error"))

				return db
			},
			want: "failed to execute user migrations: failed to begin transaction: begin error",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			migrations, err := herd.CollectFileMigrations(migrationFS)
			require.NoError(t, err)

			migrator, err := herd.New(migrations, herd.WithBuildInfoValues(testVersion, testRevision))
			require.NoError(t, err)

			db := tt.db(t)

			// act
			result, err := migrator.Migrate(t.Context(), db)

			// assert
			assert.Nil(t, result)
			assert.EqualError(t, err, tt.want)
		})
	}
}
