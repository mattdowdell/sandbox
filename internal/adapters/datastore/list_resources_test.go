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
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
)

const (
	listResourcesSQL = `SELECT resources.id AS "resources.id", resources.name AS "resources.name", ` +
		`resources.created_at AS "resources.created_at", ` +
		`resources.updated_at AS "resources.updated_at" FROM public.resources ` +
		` ORDER BY resources.id ASC LIMIT $1;`

	testLimit = 50
)

var listResourcesColumns = []string{
	"resources.id",
	"resources.name",
	"resources.created_at",
	"resources.updated_at",
}

func Test_Datastore_ListResources_Success(t *testing.T) {
	// arrange
	id := uuid.Must(uuid.NewV7())
	now := time.Now()

	db, mock := newMockDB(t)

	mock.
		ExpectQuery(listResourcesSQL).
		WithArgs(testLimit).
		WillReturnRows(
			sqlmock.NewRows(listResourcesColumns).
				AddRow(
					id,
					testResourceName,
					now,
					now.Add(time.Hour),
				),
		)

	store := datastore.NewDatastore(db)
	pager := repositories.Pager{
		Limit: testLimit,
	}

	// act
	got, err := store.ListResources(t.Context(), pager)

	// assert
	want := &repositories.Paged[*entities.Resource]{
		Items: []*entities.Resource{
			{
				ID:        id,
				Name:      testResourceName,
				CreatedAt: now,
				UpdatedAt: now.Add(time.Hour),
			},
		},
	}

	assert.Equal(t, want, got)
	assert.NoError(t, err)
}

func Test_Datastore_ListResources_Error(t *testing.T) {
	// arrange
	db, mock := newMockDB(t)

	mock.
		ExpectQuery(listResourcesSQL).
		WillReturnError(errors.New("example"))

	store := datastore.NewDatastore(db)
	pager := repositories.Pager{
		Limit: testLimit,
	}

	// act
	got, err := store.ListResources(t.Context(), pager)

	// assert
	assert.Empty(t, got)
	assert.EqualError(t, err, "internal error: jet: example")
}
