package herd_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattdowdell/sandbox/mocks/domain/mockrepositories"
	"github.com/mattdowdell/sandbox/pkg/herd"
)

func Test_New(t *testing.T) {
	// arrange
	migrations, err := herd.CollectFileMigrations(migrationFS)
	require.NoError(t, err)

	// act
	migrator, err := herd.New(migrations)

	// assert
	assert.NotNil(t, migrator)
	assert.NoError(t, err)
}

func Test_Herd_Migrate_Success(t *testing.T) {
	tests := map[string]struct {
		db   func(testing.TB) *sql.DB
		want int
	}{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			clock := mockrepositories.NewClock(t)

			migrations, err := herd.CollectFileMigrations(migrationFS)
			require.NoError(t, err)

			migrator, err := herd.New(migrations, herd.WithNowFunc(clock.Now))
			require.NoError(t, err)

			db := tt.db(t)

			// act
			version, err := migrator.Migrate(t.Context(), db)

			// assert
			assert.Equal(t, tt.want, version)
			assert.NoError(t, err)
		})
	}
}

func Test_Herd_Migrate_Error(t *testing.T) {
	tests := map[string]struct {
		db   func(testing.TB) *sql.DB
		want int
	}{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			clock := mockrepositories.NewClock(t)

			migrations, err := herd.CollectFileMigrations(migrationFS)
			require.NoError(t, err)

			migrator, err := herd.New(migrations, herd.WithNowFunc(clock.Now))
			require.NoError(t, err)

			db := tt.db(t)

			// act
			version, err := migrator.Migrate(t.Context(), db)

			// assert
			assert.Equal(t, tt.want, version)
			assert.NoError(t, err)
		})
	}
}
