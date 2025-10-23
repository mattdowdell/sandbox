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
	updateResourceSQL = `UPDATE public.resources SET (name, updated_at) = ($1, $2) ` +
		`WHERE resources.id = $3::uuid RETURNING resources.id AS "resources.id", ` +
		`resources.name AS "resources.name", resources.created_at AS "resources.created_at", ` +
		`resources.updated_at AS "resources.updated_at";`
)

var updateResourceColumns = []string{
	"resources.id",
	"resources.name",
	"resources.created_at",
	"resources.updated_at",
}

func Test_Datastore_UpdateResource_Success(t *testing.T) {
	// arrange
	id := uuid.Must(uuid.NewV7())
	now := time.Now()

	db, mock := newMockDB(t)

	mock.
		ExpectQuery(updateResourceSQL).
		WithArgs(testResourceName, now.Add(time.Hour), id).
		WillReturnRows(
			sqlmock.NewRows(updateResourceColumns).
				AddRow(
					id,
					testResourceName,
					now,
					now.Add(time.Hour),
				),
		)

	store := datastore.NewDatastore(db)

	have := &entities.Resource{
		ID:        id,
		Name:      testResourceName,
		UpdatedAt: now.Add(time.Hour),
	}

	// act
	got, err := store.UpdateResource(t.Context(), have)

	// assert
	want := &entities.Resource{
		ID:        id,
		Name:      testResourceName,
		CreatedAt: now,
		UpdatedAt: now.Add(time.Hour),
	}

	assert.Equal(t, want, got)
	assert.NoError(t, err)
}

func Test_Datastore_UpdateResource_NotFound(t *testing.T) {
	// arrange
	id := uuid.Must(uuid.NewV7())
	now := time.Now()

	db, mock := newMockDB(t)

	mock.
		ExpectQuery(updateResourceSQL).
		WithArgs(testResourceName, now.Add(time.Hour), id).
		WillReturnRows(sqlmock.NewRows(updateResourceColumns))

	store := datastore.NewDatastore(db)

	have := &entities.Resource{
		ID:        id,
		Name:      testResourceName,
		UpdatedAt: now.Add(time.Hour),
	}

	// act
	got, err := store.UpdateResource(t.Context(), have)

	// assert
	assert.Nil(t, got)
	assert.EqualError(t, err, "not found: qrm: no rows in result set")
}

func Test_Datastore_UpdateResource_QueryError(t *testing.T) {
	tests := map[string]struct {
		have error
		want string
	}{
		"unique violation": {
			have: &pgconn.PgError{
				Severity: "ERROR",
				Code:     pgerrcode.UniqueViolation,
				Message:  `duplicate key value violates unique constraint "resources_name_key"`,
			},
			want: `already exists: jet: ERROR: duplicate key value violates unique constraint ` +
				`"resources_name_key" (SQLSTATE 23505)`,
		},
		"other error": {
			have: errors.New("example"),
			want: "internal error: jet: example",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			id := uuid.Must(uuid.NewV7())
			now := time.Now()

			db, mock := newMockDB(t)

			mock.
				ExpectQuery(updateResourceSQL).
				WithArgs(testResourceName, now.Add(time.Hour), id).
				WillReturnError(tt.have)

			store := datastore.NewDatastore(db)

			have := &entities.Resource{
				ID:        id,
				Name:      testResourceName,
				UpdatedAt: now.Add(time.Hour),
			}

			// act
			got, err := store.UpdateResource(t.Context(), have)

			// assert
			assert.Nil(t, got)
			assert.EqualError(t, err, tt.want)
		})
	}
}
