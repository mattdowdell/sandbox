package txn_test

import (
	"errors"
	"testing"

	"github.com/neilotoole/slogt"
	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/adapters/txn"
	"github.com/mattdowdell/sandbox/mocks/adapters/mocktxn"
)

func Test_Func(t *testing.T) {
	// arrange
	logger := slogt.New(t)
	datastore := mocktxn.NewDatastore(t)

	commit := mocktxn.NewCommitFn(t)
	commit.EXPECT().Execute().Return(nil).Once()

	rollback := mocktxn.NewRollbackFn(t)
	rollback.EXPECT().Execute().Return(nil).Once()

	provider := mocktxn.NewProvider(t)
	provider.
		EXPECT().
		BeginTx(t.Context()).
		Return(datastore, commit.Execute, rollback.Execute, nil).
		Once()

	// act
	err := txn.Func(t.Context(), logger, provider, func(_ txn.Datastore) error {
		return nil
	})

	// assert
	assert.NoError(t, err)
}

func Test_Value(t *testing.T) {
	// arrange
	logger := slogt.New(t)
	datastore := mocktxn.NewDatastore(t)

	commit := mocktxn.NewCommitFn(t)
	commit.EXPECT().Execute().Return(nil).Once()

	rollback := mocktxn.NewRollbackFn(t)
	rollback.EXPECT().Execute().Return(nil).Once()

	provider := mocktxn.NewProvider(t)
	provider.
		EXPECT().
		BeginTx(t.Context()).
		Return(datastore, commit.Execute, rollback.Execute, nil).
		Once()

	// act
	val, err := txn.Value(t.Context(), logger, provider, func(_ txn.Datastore) (bool, error) {
		return true, nil
	})

	// assert
	assert.True(t, val)
	assert.NoError(t, err)
}

func Test_Values_Success(t *testing.T) {
	testCases := []struct {
		name        string
		rollbackErr error
	}{
		{
			name:        "no error",
			rollbackErr: nil,
		},
		{
			name:        "rollback error",
			rollbackErr: errors.New("example"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			logger := slogt.New(t)
			datastore := mocktxn.NewDatastore(t)

			commit := mocktxn.NewCommitFn(t)
			commit.EXPECT().Execute().Return(nil).Once()

			rollback := mocktxn.NewRollbackFn(t)
			rollback.EXPECT().Execute().Return(tc.rollbackErr).Once()

			provider := mocktxn.NewProvider(t)
			provider.
				EXPECT().
				BeginTx(t.Context()).
				Return(datastore, commit.Execute, rollback.Execute, nil).
				Once()

			// act
			val1, val2, err := txn.Values(
				t.Context(),
				logger,
				provider,
				func(_ txn.Datastore) (bool, bool, error) {
					return true, true, nil
				},
			)

			// assert
			assert.True(t, val1)
			assert.True(t, val2)
			assert.NoError(t, err)
		})
	}
}

func Test_Values_Error(t *testing.T) {
	testCases := []struct {
		name     string
		provider func(*testing.T) txn.Provider
		fn       func(txn.Datastore) (bool, bool, error)
		want     string
	}{
		{
			name: "begin error",
			provider: func(t *testing.T) txn.Provider {
				t.Helper()

				p := mocktxn.NewProvider(t)
				p.
					EXPECT().
					BeginTx(t.Context()).
					Return(nil, nil, nil, errors.New("example")).
					Once()

				return p
			},
			fn:   nil,
			want: "failed to begin transaction: example",
		},
		{
			name: "fn error",
			provider: func(t *testing.T) txn.Provider {
				t.Helper()

				datastore := mocktxn.NewDatastore(t)
				commit := mocktxn.NewCommitFn(t)

				rollback := mocktxn.NewRollbackFn(t)
				rollback.EXPECT().Execute().Return(nil).Once()

				p := mocktxn.NewProvider(t)
				p.
					EXPECT().
					BeginTx(t.Context()).
					Return(datastore, commit.Execute, rollback.Execute, nil).
					Once()

				return p
			},
			fn: func(_ txn.Datastore) (bool, bool, error) {
				return false, false, errors.New("example")
			},
			want: "example",
		},
		{
			name: "commit error",
			provider: func(t *testing.T) txn.Provider {
				t.Helper()

				datastore := mocktxn.NewDatastore(t)

				commit := mocktxn.NewCommitFn(t)
				commit.EXPECT().Execute().Return(errors.New("example")).Once()

				rollback := mocktxn.NewRollbackFn(t)
				rollback.EXPECT().Execute().Return(nil).Once()

				p := mocktxn.NewProvider(t)
				p.
					EXPECT().
					BeginTx(t.Context()).
					Return(datastore, commit.Execute, rollback.Execute, nil).
					Once()

				return p
			},
			fn: func(_ txn.Datastore) (bool, bool, error) {
				return true, true, nil
			},
			want: "failed to commit transaction: example",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			logger := slogt.New(t)
			provider := tc.provider(t)

			// act
			val1, val2, err := txn.Values(t.Context(), logger, provider, tc.fn)

			// assert
			assert.False(t, val1)
			assert.False(t, val2)
			assert.EqualError(t, err, tc.want)
		})
	}
}
