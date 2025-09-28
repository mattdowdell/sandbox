package herd_test

import (
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattdowdell/sandbox/pkg/herd"
)

const (
	testName     = "name"
	testVersion  = 1
	testContents = "-- example"
)

func Test_NewFileMigrationFromFilename_Success(t *testing.T) {
	tests := map[string]struct {
		have    string
		name    string
		version int
	}{
		"unpadded version": {
			have:    "1_initial.sql",
			name:    "initial",
			version: 1,
		},
		"padded version": {
			have:    "0001_initial.sql",
			name:    "initial",
			version: 1,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange

			// act
			migration, err := herd.NewFileMigrationFromFilename(tt.have, testContents)

			// assert
			if assert.NotNil(t, migration) {
				assert.Equal(t, tt.name, migration.Name())
				assert.Equal(t, tt.version, migration.Version())
			}

			assert.NoError(t, err)
		})
	}
}

func Test_NewFileMigrationFromFilename_Error(t *testing.T) {
	tests := map[string]struct {
		have string
		want string
	}{
		"invalid extension": {
			have: "1_initial.txt",
			want: `unexpected extension for filename "1_initial.txt": want: .sql, found: .txt`,
		},
		"invalid format": {
			have: "1initial.sql",
			want: `unable to extract version and name from filename: "1initial.sql"`,
		},
		"invalid version": {
			have: "invalid_initial.sql",
			want: `unable to parse version to int for filename "invalid_initial.sql": ` +
				`strconv.Atoi: parsing "invalid": invalid syntax`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange

			// act
			migration, err := herd.NewFileMigrationFromFilename(tt.have, testContents)

			// assert
			assert.Nil(t, migration)
			assert.EqualError(t, err, tt.want)
		})
	}
}

func Test_NewFileMigration_Success(t *testing.T) {
	// arrange

	// act
	migration, err := herd.NewFileMigration(testName, testVersion, testContents)

	// assert
	assert.NotNil(t, migration)
	assert.NoError(t, err)
}

func Test_NewFileMigration_Error(t *testing.T) {
	// arrange

	// act
	migration, err := herd.NewFileMigration(testName, -1, testContents)

	// assert
	assert.Nil(t, migration)
	assert.EqualError(t, err, "version must be > 0, found: -1")
}

func Test_FileMigration_Name(t *testing.T) {
	// arrange
	migration, err := herd.NewFileMigration(testName, testVersion, testContents)
	require.NoError(t, err)

	// act
	name := migration.Name()

	// assert
	assert.Equal(t, testName, name)
}

func Test_FileMigration_Version(t *testing.T) {
	// arrange
	migration, err := herd.NewFileMigration(testName, testVersion, testContents)
	require.NoError(t, err)

	// act
	version := migration.Version()

	// assert
	assert.Equal(t, testVersion, version)
}

func Test_FileMigration_Migrate_Success(t *testing.T) {
	// arrange
	migration, err := herd.NewFileMigration(testName, testVersion, testContents)
	require.NoError(t, err)

	db, mock := newMockDB(t)
	mock.ExpectExec(testContents).WillReturnResult(driver.ResultNoRows)

	// act
	err = migration.Migrate(t.Context(), db)

	// assert
	assert.NoError(t, err)
}

func Test_FileMigration_Migrate_Error(t *testing.T) {
	migration, err := herd.NewFileMigration(testName, testVersion, testContents)
	require.NoError(t, err)

	db, mock := newMockDB(t)
	mock.ExpectExec(testContents).WillReturnError(errors.New("example"))

	// act
	err = migration.Migrate(t.Context(), db)

	// assert
	assert.EqualError(t, err, "example")
}
