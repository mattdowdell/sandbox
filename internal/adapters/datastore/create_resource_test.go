package datastore_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

const (
	testResourceName = "name"

	createResourceSQL = "INSERT INTO public.resources (id, name, created_at, updated_at) " +
		"VALUES ($1, $2, $3, $4);"
)

func Test_Datastore_CreateResource_Success(t *testing.T) {
	// arrange
	id := uuid.Must(uuid.NewV7())
	now := time.Now()

	db, mock := newMockDB(t)

	mock.
		ExpectExec(createResourceSQL).
		WithArgs(id, testResourceName, now, now.Add(time.Hour)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	store := datastore.NewDatastore(db)

	resource := &entities.Resource{
		ID:        id,
		Name:      testResourceName,
		CreatedAt: now,
		UpdatedAt: now.Add(time.Hour),
	}

	// act
	err := store.CreateResource(t.Context(), resource)

	// assert
	assert.NoError(t, err)
}

func Test_Datastore_CreateResource_Error(t *testing.T) {
	testCases := []struct {
		name string
		have error
		want string
	}{
		{
			name: "unique violation",
			have: &pgconn.PgError{
				Severity: "ERROR",
				Code:     pgerrcode.UniqueViolation,
				Message:  `duplicate key value violates unique constraint "resources_name_key"`,
			},
			want: `already exists: ERROR: duplicate key value violates unique constraint ` +
				`"resources_name_key" (SQLSTATE 23505)`,
		},
		{
			name: "other error",
			have: errors.New("example"),
			want: "internal error: example",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			id := uuid.Must(uuid.NewV7())
			now := time.Now()

			db, mock := newMockDB(t)

			mock.
				ExpectExec(createResourceSQL).
				WithArgs(id, testResourceName, now, now.Add(time.Hour)).
				WillReturnError(tc.have)

			store := datastore.NewDatastore(db)

			resource := &entities.Resource{
				ID:        id,
				Name:      testResourceName,
				CreatedAt: now,
				UpdatedAt: now.Add(time.Hour),
			}

			// act
			err := store.CreateResource(t.Context(), resource)

			// assert
			assert.EqualError(t, err, tc.want)
		})
	}
}
