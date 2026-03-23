package herd_test

import (
	"database/sql/driver"
	"path/filepath"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/pkg/herd"
)

func collectVersions(migrations []herd.Migration) []int64 {
	versions := make([]int64, 0, len(migrations))
	for _, migration := range migrations {
		versions = append(versions, migration.Version())
	}

	slices.Sort(versions)
	return versions
}

func Test_CollectFileMigrations_Success(t *testing.T) {
	// arrange
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "1_initial.sql"))
	writeFile(t, filepath.Join(dir, "0002_next.sql"))
	writeFile(t, filepath.Join(dir, "3_ignored.txt"))
	writeFile(t, filepath.Join(dir, "4_non-sequential.sql"))

	filesystem := openRootFS(t, dir)

	// act
	migrations, err := herd.CollectFileMigrations(filesystem)

	// assert
	if assert.Len(t, migrations, 3) {
		versions := collectVersions(migrations)
		assert.Equal(t, []int64{1, 2, 4}, versions)
	}

	assert.NoError(t, err)
}

func Test_CollectFileMigrations_Error(t *testing.T) {
	tests := map[string]struct {
		filename string
		want     string
	}{
		"undetectable version": {
			filename: "1.sql",
			want:     "could not extract version from filename: 1.sql",
		},
		"invalid version": {
			filename: "invalid_initial.sql",
			want: `could not parse version from filename invalid_initial.sql: ` +
				`strconv.ParseInt: parsing "invalid": invalid syntax`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			dir := t.TempDir()
			writeFile(t, filepath.Join(dir, tt.filename))

			filesystem := openRootFS(t, dir)

			// act
			migrations, err := herd.CollectFileMigrations(filesystem)

			// assert
			assert.Empty(t, migrations)
			assert.EqualError(t, err, tt.want)
		})
	}
}

func Test_NewFileMigrationFromFS_Error(t *testing.T) {
	// arrange
	dir := t.TempDir()
	filesystem := openRootFS(t, dir)

	// act
	migration, err := herd.NewFileMigrationFromFS(filesystem, "does_not_exist.sql")

	// assert
	assert.Nil(t, migration)
	assert.EqualError(
		t,
		err,
		"failed to read migration from filesystem: openat does_not_exist.sql: "+
			"no such file or directory",
	)
}

func Test_FileMigration_Migrate(t *testing.T) {
	// arrange
	migration := herd.NewFileMigration(1, "-- example")

	db, mock := newMockDB(t)
	mock.ExpectExec("-- example").WillReturnResult(driver.ResultNoRows)

	// act
	err := migration.Migrate(t.Context(), db)

	// assert
	assert.NoError(t, err)
}
