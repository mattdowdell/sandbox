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
	store, commit, rollback, err := provider.BeginTx(t.Context())

	// assert
	assert.NotNil(t, store)
	assert.NotNil(t, commit)
	assert.NotNil(t, rollback)
	assert.NoError(t, err)
}

func Test_Provider_BeginTx_Error(t *testing.T) {
	// arrange
	db, mock := newMockDB(t)

	mock.ExpectBegin().WillReturnError(errors.New("example"))

	provider := datastore.NewProvider(db)

	// act
	store, commit, rollback, err := provider.BeginTx(t.Context())

	// assert
	assert.Nil(t, store)
	assert.Nil(t, commit)
	assert.Nil(t, rollback)
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

func Test_wrapRollback_Success(t *testing.T) {
	testCases := []struct {
		name string
		have error
	}{
		{
			name: "no error",
			have: nil,
		},
		{
			name: "tx done",
			have: sql.ErrTxDone,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			db, mock := newMockDB(t)

			mock.ExpectBegin()
			mock.ExpectRollback().WillReturnError(tc.have)

			provider := datastore.NewProvider(db)

			_, _, rollback, err := provider.BeginTx(t.Context())
			require.NoError(t, err)

			// act
			err = rollback()

			// assert
			assert.NoError(t, err)
		})
	}
}

func Test_wrapRollback_Error(t *testing.T) {
	// arrange
	db, mock := newMockDB(t)

	mock.ExpectBegin()
	mock.ExpectRollback().WillReturnError(errors.New("example"))

	provider := datastore.NewProvider(db)

	_, _, rollback, err := provider.BeginTx(t.Context())
	require.NoError(t, err)

	// act
	err = rollback()

	// assert
	assert.EqualError(t, err, "example")
}
