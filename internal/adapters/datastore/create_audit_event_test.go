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
	createAuditEventSQL = `INSERT INTO public.audit_events ` +
		`(id, operation, created_at, summary, resource_id, resource_type) ` +
		`VALUES ($1, $2, $3, $4, $5, $6);`
)

func Test_Datastore_CreateAuditEvent_Success(t *testing.T) {
	// arrange
	id := uuid.Must(uuid.NewV7())
	now := time.Now()
	resourceID := uuid.Must(uuid.NewV7())

	db, mock := newMockDB(t)

	mock.
		ExpectExec(createAuditEventSQL).
		WithArgs(id, "created", now, "summary", resourceID, "resource").
		WillReturnResult(sqlmock.NewResult(0, 1))

	store := datastore.NewDatastore(db)

	event := &entities.AuditEvent{
		ID:           id,
		Operation:    entities.OperationCreated,
		CreatedAt:    now,
		Summary:      "summary",
		ResourceID:   resourceID,
		ResourceType: entities.ResourceTypeResource,
	}

	// act
	err := store.CreateAuditEvent(t.Context(), event)

	// assert
	assert.NoError(t, err)
}

func Test_Datastore_CreateAuditEvent_Error(t *testing.T) {
	// arrange
	id := uuid.Must(uuid.NewV7())
	now := time.Now()
	resourceID := uuid.Must(uuid.NewV7())

	db, mock := newMockDB(t)

	mock.
		ExpectExec(createAuditEventSQL).
		WithArgs(id, "created", now, "summary", resourceID, "resource").
		WillReturnError(errors.New("example"))

	store := datastore.NewDatastore(db)

	event := &entities.AuditEvent{
		ID:           id,
		Operation:    entities.OperationCreated,
		CreatedAt:    now,
		Summary:      "summary",
		ResourceID:   resourceID,
		ResourceType: entities.ResourceTypeResource,
	}

	// act
	err := store.CreateAuditEvent(t.Context(), event)

	// assert
	assert.EqualError(t, err, "internal error: example")
}
