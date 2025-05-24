package datastore_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore"
)

const (
	deleteResourceSQL = "DELETE FROM public.resources WHERE resources.id = $1;"
)

func Test_Datastore_DeleteResource_Success(t *testing.T) {
	// arrange
	id := uuid.New()

	db, mock := newMockDB(t)

	mock.
		ExpectExec(deleteResourceSQL).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	store := datastore.NewDatastore(db)

	// act
	err := store.DeleteResource(t.Context(), id)

	// assert
	assert.NoError(t, err)
}

func Test_Datastore_DeleteResource_ExecError(t *testing.T) {
	// arrange
	id := uuid.New()

	db, mock := newMockDB(t)

	mock.
		ExpectExec(deleteResourceSQL).
		WithArgs(id).
		WillReturnError(errors.New("example"))

	store := datastore.NewDatastore(db)

	// act
	err := store.DeleteResource(t.Context(), id)

	// assert
	assert.EqualError(t, err, "internal error: example")
}

func Test_Datastore_DeleteResource_RowsAffectedError(t *testing.T) {
	// arrange
	id := uuid.New()

	db, mock := newMockDB(t)

	mock.
		ExpectExec(deleteResourceSQL).
		WithArgs(id).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("example")))

	store := datastore.NewDatastore(db)

	// act
	err := store.DeleteResource(t.Context(), id)

	// assert
	assert.EqualError(t, err, "internal error: example")
}

func Test_Datastore_DeleteResource_RowsAffectedInvalid(t *testing.T) {
	id := uuid.New()

	testCases := []struct {
		name string
		have int64
		want string
	}{
		{
			name: "no rows affected",
			have: 0,
			want: fmt.Sprintf("not found: %s", id),
		},
		{
			name: "multiple rows affected",
			have: 2,
			want: "internal error: too many deletes: 2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			db, mock := newMockDB(t)

			mock.
				ExpectExec(deleteResourceSQL).
				WithArgs(id).
				WillReturnResult(sqlmock.NewResult(0, tc.have))

			store := datastore.NewDatastore(db)

			// act
			err := store.DeleteResource(t.Context(), id)

			// assert
			assert.EqualError(t, err, tc.want)
		})
	}
}
