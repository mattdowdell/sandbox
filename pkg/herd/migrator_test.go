package herd_test

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattdowdell/sandbox/pkg/herd"
)

const (
	testCodeVersion      = "v0.0.1"
	testCodeRevision     = "abcdef0123456789"
	testHerdVersion      = 1
	testMigrationVersion = 2
	testMigration        = "-- example"
	testTargetVersion    = 0
)

var (
	testQueryRows  = []string{"migration_version"}
	testExistsRows = []string{"exists"}
)

func Test_newMigrator_Success(t *testing.T) {
	tests := map[string]struct {
		have []herd.Migration
	}{
		"empty": {},
		"single": {
			have: []herd.Migration{
				herd.NewFileMigration(1, ""),
			},
		},
		"multiple sequential": {
			have: []herd.Migration{
				herd.NewFileMigration(1, ""),
				herd.NewFileMigration(2, ""),
				herd.NewFileMigration(3, ""),
			},
		},
		"multiple non-sequential": {
			have: []herd.Migration{
				herd.NewFileMigration(5, ""),
				herd.NewFileMigration(10, ""),
				herd.NewFileMigration(15, ""),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			recorder := herd.NewRecorder(
				time.Now,
				herd.TableNameUser,
				testCodeVersion,
				testCodeRevision,
			)

			// act
			migrator, err := herd.NewMigrator(tt.have, recorder, testTargetVersion)

			// assert
			assert.NotNil(t, migrator)
			assert.NoError(t, err)
		})
	}
}

func Test_newMigrator_Error(t *testing.T) {
	tests := map[string]struct {
		have []herd.Migration
		want string
	}{
		"non-positive version": {
			have: []herd.Migration{
				herd.NewFileMigration(0, ""),
			},
			want: "migration version must be > 0, found: 0",
		},
		"duplicate versions": {
			have: []herd.Migration{
				herd.NewFileMigration(1, ""),
				herd.NewFileMigration(2, ""),
				herd.NewFileMigration(3, ""),
				herd.NewFileMigration(1, ""),
			},
			want: "duplicate migration version found: 1",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			recorder := herd.NewRecorder(
				time.Now,
				herd.TableNameUser,
				testCodeVersion,
				testCodeRevision,
			)

			// act
			migrator, err := herd.NewMigrator(tt.have, recorder, testTargetVersion)

			// assert
			assert.Nil(t, migrator)
			assert.EqualError(t, err, tt.want)
		})
	}
}

func Test_migrator_Migrate_Success(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	tests := map[string]struct {
		migrations []herd.Migration
		db         func(*testing.T) *sql.DB
		want       *herd.Result
	}{
		"empty migrations, unmigrated db": {
			migrations: []herd.Migration{},
			db: func(t *testing.T) *sql.DB {
				t.Helper()

				db, mock := newMockDB(t)
				mock.ExpectBegin()
				mock.ExpectQuery(herd.UserVersionQuery).
					WillReturnRows(sqlmock.NewRows(testQueryRows))
				mock.ExpectCommit()

				return db
			},
			want: &herd.Result{},
		},
		"no pending migrations, migrated db": {
			migrations: []herd.Migration{
				herd.NewFileMigration(testMigrationVersion, testMigration),
			},
			db: func(t *testing.T) *sql.DB {
				t.Helper()

				db, mock := newMockDB(t)
				mock.ExpectBegin()
				mock.ExpectQuery(herd.UserVersionQuery).
					WillReturnRows(sqlmock.NewRows(testQueryRows).AddRow(testMigrationVersion))
				mock.ExpectCommit()

				return db
			},
			want: &herd.Result{
				Before: testMigrationVersion,
				After:  testMigrationVersion,
			},
		},
		"pending migrations, unmigrated db": {
			migrations: []herd.Migration{
				herd.NewFileMigration(testMigrationVersion, testMigration),
			},
			db: func(t *testing.T) *sql.DB {
				t.Helper()

				db, mock := newMockDB(t)
				mock.ExpectBegin()
				mock.ExpectQuery(herd.UserVersionQuery).
					WillReturnRows(sqlmock.NewRows(testQueryRows))
				mock.ExpectExec(testMigration).WillReturnResult(driver.ResultNoRows)
				mock.ExpectExec(herd.UserRecordQuery).
					WithArgs(testMigrationVersion, now, testCodeVersion, testCodeRevision, testHerdVersion).
					WillReturnResult(driver.ResultNoRows)
				mock.ExpectCommit()

				return db
			},
			want: &herd.Result{
				Before: 0,
				After:  testMigrationVersion,
			},
		},
		"pending migrations, partially migrated db": {
			migrations: []herd.Migration{
				herd.NewFileMigration(testMigrationVersion, testMigration),
			},
			db: func(t *testing.T) *sql.DB {
				t.Helper()

				db, mock := newMockDB(t)
				mock.ExpectBegin()
				mock.ExpectQuery(herd.UserVersionQuery).
					WillReturnRows(sqlmock.NewRows(testQueryRows).AddRow(1))
				mock.ExpectExec(testMigration).WillReturnResult(driver.ResultNoRows)
				mock.ExpectExec(herd.UserRecordQuery).
					WithArgs(testMigrationVersion, now, testCodeVersion, testCodeRevision, testHerdVersion).
					WillReturnResult(driver.ResultNoRows)
				mock.ExpectCommit()

				return db
			},
			want: &herd.Result{
				Before: 1,
				After:  testMigrationVersion,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			recorder := herd.NewRecorder(
				testNowFunc(now),
				herd.TableNameUser,
				testCodeVersion,
				testCodeRevision,
			)

			migrator, err := herd.NewMigrator(tt.migrations, recorder, testTargetVersion)
			require.NoError(t, err)

			db := tt.db(t)

			// act
			result, err := migrator.Migrate(t.Context(), db, testHerdVersion)

			// assert
			assert.Equal(t, tt.want, result)
			assert.NoError(t, err)
		})
	}
}

func Test_migrator_Migrate_Error(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	tests := map[string]struct {
		migrations []herd.Migration
		db         func(*testing.T) *sql.DB
		want       string
	}{
		"begin transaction": {
			migrations: []herd.Migration{},
			db: func(t *testing.T) *sql.DB {
				t.Helper()

				db, mock := newMockDB(t)
				mock.ExpectBegin().WillReturnError(errors.New("begin error"))

				return db
			},
			want: "failed to begin transaction: begin error",
		},
		"get current version": {
			migrations: []herd.Migration{},
			db: func(t *testing.T) *sql.DB {
				t.Helper()

				db, mock := newMockDB(t)
				mock.ExpectBegin()
				mock.ExpectQuery(herd.UserVersionQuery).WillReturnError(errors.New("query error"))
				mock.ExpectRollback()

				return db
			},
			want: "failed to scan current migration version: query error",
		},
		"apply migration": {
			migrations: []herd.Migration{
				herd.NewFileMigration(testMigrationVersion, testMigration),
			},
			db: func(t *testing.T) *sql.DB {
				t.Helper()

				db, mock := newMockDB(t)
				mock.ExpectBegin()
				mock.ExpectQuery(herd.UserVersionQuery).
					WillReturnRows(sqlmock.NewRows(testQueryRows))
				mock.ExpectExec(testMigration).WillReturnError(errors.New("exec error"))
				mock.ExpectRollback()

				return db
			},
			want: "failed to execute migration 2: exec error",
		},
		"record migration": {
			migrations: []herd.Migration{
				herd.NewFileMigration(testMigrationVersion, testMigration),
			},
			db: func(t *testing.T) *sql.DB {
				t.Helper()

				db, mock := newMockDB(t)
				mock.ExpectBegin()
				mock.ExpectQuery(herd.UserVersionQuery).
					WillReturnRows(sqlmock.NewRows(testQueryRows))
				mock.ExpectExec(testMigration).WillReturnResult(driver.ResultNoRows)
				mock.ExpectExec(herd.UserRecordQuery).WillReturnError(errors.New("exec error"))
				mock.ExpectRollback()

				return db
			},
			want: "failed to record migration 2: exec error",
		},
		"commit": {
			migrations: []herd.Migration{
				herd.NewFileMigration(testMigrationVersion, testMigration),
			},
			db: func(t *testing.T) *sql.DB {
				t.Helper()

				db, mock := newMockDB(t)
				mock.ExpectBegin()
				mock.ExpectQuery(herd.UserVersionQuery).
					WillReturnRows(sqlmock.NewRows(testQueryRows))
				mock.ExpectExec(testMigration).WillReturnResult(driver.ResultNoRows)
				mock.ExpectExec(herd.UserRecordQuery).WillReturnResult(driver.ResultNoRows)
				mock.ExpectCommit().WillReturnError(errors.New("commit error"))

				return db
			},
			want: "failed to commit transaction: commit error",
		},
		"rollback": {
			migrations: []herd.Migration{
				herd.NewFileMigration(testMigrationVersion, testMigration),
			},
			db: func(t *testing.T) *sql.DB {
				t.Helper()

				db, mock := newMockDB(t)
				mock.ExpectBegin()
				mock.ExpectQuery(herd.UserVersionQuery).
					WillReturnRows(sqlmock.NewRows(testQueryRows))
				mock.ExpectExec(testMigration).WillReturnResult(driver.ResultNoRows)
				mock.ExpectExec(herd.UserRecordQuery).WillReturnError(errors.New("exec error"))
				mock.ExpectRollback().WillReturnError(errors.New("rollback error"))

				return db
			},
			want: "failed to record migration 2: exec error\n" +
				"failed to rollback transaction: rollback error",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			recorder := herd.NewRecorder(
				testNowFunc(now),
				herd.TableNameUser,
				testCodeVersion,
				testCodeRevision,
			)

			migrator, err := herd.NewMigrator(tt.migrations, recorder, testTargetVersion)
			require.NoError(t, err)

			db := tt.db(t)

			// act
			result, err := migrator.Migrate(t.Context(), db, testHerdVersion)

			// assert
			assert.Nil(t, result)
			assert.EqualError(t, err, tt.want)
		})
	}
}
