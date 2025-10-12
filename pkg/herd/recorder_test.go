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

	"github.com/mattdowdell/sandbox/mocks/domain/mockrepositories"
	"github.com/mattdowdell/sandbox/pkg/herd"
)

func Test_newRecorder(t *testing.T) {
	tests := map[string]struct {
		table string
	}{
		"system": {
			table: herd.TableNameSystem,
		},
		"user": {
			table: herd.TableNameUser,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			clock := mockrepositories.NewClock(t)

			// act
			recorder := herd.NewRecorder(clock.Now, tt.table, testCodeVersion, testCodeRevision)

			// assert
			assert.NotNil(t, recorder)
		})
	}
}

func Test_recorder_GetCurrentVersion_Success(t *testing.T) {
	tests := map[string]struct {
		table string
		query string
	}{
		"system": {
			table: herd.TableNameSystem,
			query: herd.SystemVersionQuery,
		},
		"user": {
			table: herd.TableNameUser,
			query: herd.UserVersionQuery,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			clock := mockrepositories.NewClock(t)
			recorder := herd.NewRecorder(clock.Now, tt.table, testCodeVersion, testCodeRevision)

			db, mock := newMockDB(t)
			mock.ExpectBegin()
			mock.ExpectQuery(tt.query).WillReturnRows(sqlmock.NewRows(testQueryRows).AddRow(1))

			tx, err := db.BeginTx(t.Context(), &sql.TxOptions{})
			require.NoError(t, err)

			// act
			version, err := recorder.GetCurrentVersion(t.Context(), tx)

			// assert
			assert.Equal(t, int64(1), version)
			assert.NoError(t, err)
		})
	}
}

func Test_recorder_GetCurrentVersion_Error(t *testing.T) {
	tests := map[string]struct {
		table string
		db    func(*testing.T) *sql.DB
		want  string
	}{
		"invalid table": {
			table: "invalid",
			db: func(t *testing.T) *sql.DB {
				t.Helper()

				db, mock := newMockDB(t)
				mock.ExpectBegin()

				return db
			},
			want: "internal error: unexpected table name: invalid",
		},
		"query error": {
			table: herd.TableNameSystem,
			db: func(t *testing.T) *sql.DB {
				t.Helper()

				db, mock := newMockDB(t)
				mock.ExpectBegin()
				mock.ExpectQuery(herd.SystemVersionQuery).WillReturnError(errors.New("query error"))

				return db
			},
			want: "failed to scan current migration version: query error",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(_ *testing.T) {
			// arrange
			clock := mockrepositories.NewClock(t)
			recorder := herd.NewRecorder(clock.Now, tt.table, testCodeVersion, testCodeRevision)

			db := tt.db(t)

			tx, err := db.BeginTx(t.Context(), &sql.TxOptions{})
			require.NoError(t, err)

			// act
			version, err := recorder.GetCurrentVersion(t.Context(), tx)

			// assert
			assert.Zero(t, version)
			assert.EqualError(t, err, tt.want)
		})
	}
}

func Test_recorder_RecordMigration_Success(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	tests := map[string]struct {
		table string
		query string
		args  []driver.Value
	}{
		"system": {
			table: herd.TableNameSystem,
			query: herd.SystemRecordQuery,
			args: []driver.Value{
				testMigrationVersion,
				now,
				testCodeVersion,
				testCodeRevision,
			},
		},
		"user": {
			table: herd.TableNameUser,
			query: herd.UserRecordQuery,
			args: []driver.Value{
				testMigrationVersion,
				now,
				testCodeVersion,
				testCodeRevision,
				testHerdVersion,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			clock := mockrepositories.NewClock(t)
			clock.EXPECT().Now().Return(now).Once()

			recorder := herd.NewRecorder(clock.Now, tt.table, testCodeVersion, testCodeRevision)

			db, mock := newMockDB(t)
			mock.ExpectBegin()
			mock.ExpectExec(tt.query).WithArgs(tt.args...).WillReturnResult(driver.ResultNoRows)

			tx, err := db.BeginTx(t.Context(), &sql.TxOptions{})
			require.NoError(t, err)

			// act
			err = recorder.RecordMigration(t.Context(), tx, testMigrationVersion, testHerdVersion)

			// assert
			assert.NoError(t, err)
		})
	}
}

func Test_recorder_RecordMigration_Error(t *testing.T) {
	tests := map[string]struct {
		table string
		db    func(*testing.T) *sql.DB
		want  string
	}{
		"invalid table": {
			table: "invalid",
			db: func(t *testing.T) *sql.DB {
				t.Helper()

				db, mock := newMockDB(t)
				mock.ExpectBegin()

				return db
			},
			want: "internal error: unexpected table name: invalid",
		},
		"exec error": {
			table: herd.TableNameSystem,
			db: func(t *testing.T) *sql.DB {
				t.Helper()

				db, mock := newMockDB(t)
				mock.ExpectBegin()
				mock.ExpectExec(herd.SystemRecordQuery).WillReturnError(errors.New("exec error"))

				return db
			},
			want: "exec error",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(_ *testing.T) {
			// arrange
			now := time.Now().UTC().Truncate(time.Second)

			clock := mockrepositories.NewClock(t)
			clock.EXPECT().Now().Return(now).Maybe()

			recorder := herd.NewRecorder(clock.Now, tt.table, testCodeVersion, testCodeRevision)

			db := tt.db(t)

			tx, err := db.BeginTx(t.Context(), &sql.TxOptions{})
			require.NoError(t, err)

			// act
			err = recorder.RecordMigration(t.Context(), tx, testMigrationVersion, testHerdVersion)

			// assert
			assert.EqualError(t, err, tt.want)
		})
	}
}
