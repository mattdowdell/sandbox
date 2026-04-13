package datastore_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore"
)

func Test_NewProvider(t *testing.T) {
	// arrange
	db, _ := newMockDB(t)

	// act
	provider := datastore.NewProvider(db)

	// assert
	assert.NotNil(t, provider)
}

func Test_Provider_BeginTx_Success(t *testing.T) {
	// arrange
	db, mock := newMockDB(t)

	mock.ExpectBegin()

	provider := datastore.NewProvider(db)

	// act
	store, ender, err := provider.BeginTx(t.Context())

	// assert
	assert.NotNil(t, store)
	assert.NotNil(t, ender)
	assert.NoError(t, err)
}

func Test_Provider_BeginTx_Error(t *testing.T) {
	// arrange
	db, mock := newMockDB(t)

	mock.ExpectBegin().WillReturnError(errors.New("example"))

	provider := datastore.NewProvider(db)

	// act
	store, ender, err := provider.BeginTx(t.Context())

	// assert
	assert.Nil(t, store)
	assert.Nil(t, ender)
	assert.EqualError(t, err, "example")
}

func Test_Provider_Datastore(t *testing.T) {
	// arrange
	db, _ := newMockDB(t)

	provider := datastore.NewProvider(db)

	// act
	store := provider.Datastore()

	// assert
	assert.NotNil(t, store)
}

func Test_Ender_Rollback_Success(t *testing.T) {
	tests := map[string]struct {
		have error
	}{
		"no error": {
			have: nil,
		},
		"tx done": {
			have: sql.ErrTxDone,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			db, mock := newMockDB(t)

			mock.ExpectBegin()
			mock.ExpectRollback().WillReturnError(tt.have)

			provider := datastore.NewProvider(db)

			_, ender, err := provider.BeginTx(t.Context())
			require.NoError(t, err)

			// act
			err = ender.Rollback()

			// assert
			assert.NoError(t, err)
		})
	}
}

func Test_Ender_Rollback_Error(t *testing.T) {
	// arrange
	db, mock := newMockDB(t)

	mock.ExpectBegin()
	mock.ExpectRollback().WillReturnError(errors.New("example"))

	provider := datastore.NewProvider(db)

	_, ender, err := provider.BeginTx(t.Context())
	require.NoError(t, err)

	// act
	err = ender.Rollback()

	// assert
	assert.EqualError(t, err, "example")
}
