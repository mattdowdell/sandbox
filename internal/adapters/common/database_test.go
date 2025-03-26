package common_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/adapters/common"
	"github.com/mattdowdell/sandbox/mocks/adapters/mockcommon"
)

func Test_TxFunc(t *testing.T) {
	// arrange
	datastore := mockcommon.NewDatastore(t)

	commit := mockcommon.NewCommitFn(t)
	commit.EXPECT().Execute().Return(nil).Once()

	rollback := mockcommon.NewRollbackFn(t)
	rollback.EXPECT().Execute().Return(nil).Once()

	provider := mockcommon.NewProvider(t)
	provider.
		EXPECT().
		BeginTx(t.Context()).
		Return(datastore, commit.Execute, rollback.Execute, nil).
		Once()

	// act
	err := common.TxFunc(t.Context(), provider, func(_ common.Datastore) error {
		return nil
	})

	// assert
	assert.NoError(t, err)
}

func Test_TxValue(t *testing.T) {
	// arrange
	datastore := mockcommon.NewDatastore(t)

	commit := mockcommon.NewCommitFn(t)
	commit.EXPECT().Execute().Return(nil).Once()

	rollback := mockcommon.NewRollbackFn(t)
	rollback.EXPECT().Execute().Return(nil).Once()

	provider := mockcommon.NewProvider(t)
	provider.
		EXPECT().
		BeginTx(t.Context()).
		Return(datastore, commit.Execute, rollback.Execute, nil).
		Once()

	// act
	val, err := common.TxValue(t.Context(), provider, func(_ common.Datastore) (bool, error) {
		return true, nil
	})

	// assert
	assert.True(t, val)
	assert.NoError(t, err)
}

func Test_TxValues_Success(t *testing.T) {
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
			datastore := mockcommon.NewDatastore(t)

			commit := mockcommon.NewCommitFn(t)
			commit.EXPECT().Execute().Return(nil).Once()

			rollback := mockcommon.NewRollbackFn(t)
			rollback.EXPECT().Execute().Return(tc.rollbackErr).Once()

			provider := mockcommon.NewProvider(t)
			provider.
				EXPECT().
				BeginTx(t.Context()).
				Return(datastore, commit.Execute, rollback.Execute, nil).
				Once()

			// act
			val1, val2, err := common.TxValues(
				t.Context(),
				provider,
				func(_ common.Datastore) (bool, bool, error) {
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

func Test_TxValues_Error(t *testing.T) {
	testCases := []struct {
		name     string
		provider func(*testing.T) common.Provider
		fn       func(common.Datastore) (bool, bool, error)
		want     string
	}{
		{
			name: "begin error",
			provider: func(t *testing.T) common.Provider {
				t.Helper()

				p := mockcommon.NewProvider(t)
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
			provider: func(t *testing.T) common.Provider {
				t.Helper()

				datastore := mockcommon.NewDatastore(t)
				commit := mockcommon.NewCommitFn(t)

				rollback := mockcommon.NewRollbackFn(t)
				rollback.EXPECT().Execute().Return(nil).Once()

				p := mockcommon.NewProvider(t)
				p.
					EXPECT().
					BeginTx(t.Context()).
					Return(datastore, commit.Execute, rollback.Execute, nil).
					Once()

				return p
			},
			fn: func(_ common.Datastore) (bool, bool, error) {
				return false, false, errors.New("example")
			},
			want: "example",
		},
		{
			name: "commit error",
			provider: func(t *testing.T) common.Provider {
				t.Helper()

				datastore := mockcommon.NewDatastore(t)

				commit := mockcommon.NewCommitFn(t)
				commit.EXPECT().Execute().Return(errors.New("example")).Once()

				rollback := mockcommon.NewRollbackFn(t)
				rollback.EXPECT().Execute().Return(nil).Once()

				p := mockcommon.NewProvider(t)
				p.
					EXPECT().
					BeginTx(t.Context()).
					Return(datastore, commit.Execute, rollback.Execute, nil).
					Once()

				return p
			},
			fn: func(_ common.Datastore) (bool, bool, error) {
				return true, true, nil
			},
			want: "failed to commit transaction: example",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			provider := tc.provider(t)

			// act
			val1, val2, err := common.TxValues(t.Context(), provider, tc.fn)

			// assert
			assert.False(t, val1)
			assert.False(t, val2)
			assert.EqualError(t, err, tc.want)
		})
	}
}
