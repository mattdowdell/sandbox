package txn_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/adapters/txn"
	"github.com/mattdowdell/sandbox/internal/adapters/txn/mocktxn"
	"github.com/mattdowdell/sandbox/pkg/slogt"
)

func Test_Func(t *testing.T) {
	// arrange
	logger := slogt.New(t)
	datastore := mocktxn.NewMockDatastore(t)

	ender := mocktxn.NewMockEnder(t)
	ender.EXPECT().Commit().Return(nil).Once()
	ender.EXPECT().Rollback().Return(nil).Once()

	provider := mocktxn.NewMockProvider(t)
	provider.
		EXPECT().
		BeginTx(t.Context()).
		Return(datastore, ender, nil).
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
	datastore := mocktxn.NewMockDatastore(t)

	ender := mocktxn.NewMockEnder(t)
	ender.EXPECT().Commit().Return(nil).Once()
	ender.EXPECT().Rollback().Return(nil).Once()

	provider := mocktxn.NewMockProvider(t)
	provider.
		EXPECT().
		BeginTx(t.Context()).
		Return(datastore, ender, nil).
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
	tests := map[string]struct {
		rollbackErr error
	}{
		"no error": {
			rollbackErr: nil,
		},
		"rollback error": {
			rollbackErr: errors.New("example"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			logger := slogt.New(t)
			datastore := mocktxn.NewMockDatastore(t)

			ender := mocktxn.NewMockEnder(t)
			ender.EXPECT().Commit().Return(nil).Once()
			ender.EXPECT().Rollback().Return(tt.rollbackErr).Once()

			provider := mocktxn.NewMockProvider(t)
			provider.
				EXPECT().
				BeginTx(t.Context()).
				Return(datastore, ender, nil).
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
	tests := map[string]struct {
		provider func(*testing.T) txn.Provider
		fn       func(txn.Datastore) (bool, bool, error)
		want     string
	}{
		"begin error": {
			provider: func(t *testing.T) txn.Provider {
				t.Helper()

				p := mocktxn.NewMockProvider(t)
				p.
					EXPECT().
					BeginTx(t.Context()).
					Return(nil, nil, errors.New("example")).
					Once()

				return p
			},
			fn:   nil,
			want: "failed to begin transaction: example",
		},
		"fn error": {
			provider: func(t *testing.T) txn.Provider {
				t.Helper()

				datastore := mocktxn.NewMockDatastore(t)
				ender := mocktxn.NewMockEnder(t)
				ender.EXPECT().Rollback().Return(nil).Once()

				p := mocktxn.NewMockProvider(t)
				p.
					EXPECT().
					BeginTx(t.Context()).
					Return(datastore, ender, nil).
					Once()

				return p
			},
			fn: func(_ txn.Datastore) (bool, bool, error) {
				return false, false, errors.New("example")
			},
			want: "example",
		},
		"commit error": {
			provider: func(t *testing.T) txn.Provider {
				t.Helper()

				datastore := mocktxn.NewMockDatastore(t)

				ender := mocktxn.NewMockEnder(t)
				ender.EXPECT().Commit().Return(errors.New("example")).Once()
				ender.EXPECT().Rollback().Return(nil).Once()

				p := mocktxn.NewMockProvider(t)
				p.
					EXPECT().
					BeginTx(t.Context()).
					Return(datastore, ender, nil).
					Once()

				return p
			},
			fn: func(_ txn.Datastore) (bool, bool, error) {
				return true, true, nil
			},
			want: "failed to commit transaction: example",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			logger := slogt.New(t)
			provider := tt.provider(t)

			// act
			val1, val2, err := txn.Values(t.Context(), logger, provider, tt.fn)

			// assert
			assert.False(t, val1)
			assert.False(t, val2)
			assert.EqualError(t, err, tt.want)
		})
	}
}
