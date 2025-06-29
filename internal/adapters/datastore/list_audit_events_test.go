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
	listAuditEventsSQL = `SELECT audit_events.id AS "audit_events.id", ` +
		`audit_events.operation AS "audit_events.operation", ` +
		`audit_events.created_at AS "audit_events.created_at", ` +
		`audit_events.summary AS "audit_events.summary", ` +
		`audit_events.resource_id AS "audit_events.resource_id", ` +
		`audit_events.resource_type AS "audit_events.resource_type" FROM public.audit_events ` +
		`ORDER BY resources.id ASC;`
)

var listAuditEventsColumns = []string{
	"audit_events.id",
	"audit_events.operation",
	"audit_events.created_at",
	"audit_events.summary",
	"audit_events.resource_id",
	"audit_events.resource_type",
}

func Test_Datastore_ListAuditEvents_Success(t *testing.T) {
	// arrange
	id := uuid.Must(uuid.NewV7())
	now := time.Now()
	resourceID := uuid.Must(uuid.NewV7())

	db, mock := newMockDB(t)

	mock.
		ExpectQuery(listAuditEventsSQL).
		WillReturnRows(
			sqlmock.NewRows(listAuditEventsColumns).
				AddRow(
					id,
					"created",
					now,
					"summary",
					resourceID,
					"resource",
				),
		)

	store := datastore.NewDatastore(db)

	// act
	got, err := store.ListAuditEvents(t.Context())

	// assert
	want := []*entities.AuditEvent{
		{
			ID:           id,
			Operation:    entities.OperationCreated,
			CreatedAt:    now,
			Summary:      "summary",
			ResourceID:   resourceID,
			ResourceType: entities.ResourceTypeResource,
		},
	}

	assert.Equal(t, want, got)
	assert.NoError(t, err)
}

func Test_Datastore_ListAuditEvents_Error(t *testing.T) {
	// arrange
	db, mock := newMockDB(t)

	mock.
		ExpectQuery(listAuditEventsSQL).
		WillReturnError(errors.New("example"))

	store := datastore.NewDatastore(db)

	// act
	got, err := store.ListAuditEvents(t.Context())

	// assert
	assert.Empty(t, got)
	assert.EqualError(t, err, "internal error: jet: example")
}
