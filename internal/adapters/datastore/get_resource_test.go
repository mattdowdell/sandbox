package datastore_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

const (
	getResourceSQL = `SELECT resources.id AS "resources.id", resources.name AS "resources.name", ` +
		`resources.created_at AS "resources.created_at", ` +
		`resources.updated_at AS "resources.updated_at" FROM public.resources ` +
		`WHERE resources.id = $1::uuid;`
)

var getResourceColumns = []string{
	"resources.id",
	"resources.name",
	"resources.created_at",
	"resources.updated_at",
}

func Test_Datastore_GetResource_Success(t *testing.T) {
	// arrange
	id := uuid.Must(uuid.NewV7())
	now := time.Now()

	db, mock := newMockDB(t)

	mock.
		ExpectQuery(getResourceSQL).
		WithArgs(id).
		WillReturnRows(
			sqlmock.NewRows(getResourceColumns).
				AddRow(
					id,
					testResourceName,
					now,
					now.Add(time.Hour),
				),
		)

	store := datastore.NewDatastore(db)

	// act
	got, err := store.GetResource(t.Context(), id)

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

func Test_Datastore_GetResource_NotFound(t *testing.T) {
	// arrange
	id := uuid.Must(uuid.NewV7())

	db, mock := newMockDB(t)

	mock.
		ExpectQuery(getResourceSQL).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows(getResourceColumns))

	store := datastore.NewDatastore(db)

	// act
	got, err := store.GetResource(t.Context(), id)

	// assert
	assert.Empty(t, got)
	assert.EqualError(t, err, "not found: qrm: no rows in result set")
}

func Test_Datastore_GetResource_Error(t *testing.T) {
	// arrange
	id := uuid.Must(uuid.NewV7())

	db, mock := newMockDB(t)

	mock.
		ExpectQuery(getResourceSQL).
		WithArgs(id).
		WillReturnError(errors.New("example"))

	store := datastore.NewDatastore(db)

	// act
	got, err := store.GetResource(t.Context(), id)

	// assert
	assert.Nil(t, got)
	assert.EqualError(t, err, "internal error: jet: example")
}
