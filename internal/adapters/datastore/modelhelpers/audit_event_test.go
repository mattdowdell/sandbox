package modelhelpers_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore/modelhelpers"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/public/model"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

const (
	testSummary = "summary"
)

func Test_AuditEventsToDomain(t *testing.T) {
	// arrange
	id := uuid.New()
	now := time.Now()
	resourceID := uuid.New()

	have := []*model.AuditEvents{
		{
			ID:           id,
			Operation:    "created",
			CreatedAt:    now,
			Summary:      testSummary,
			ResourceID:   resourceID,
			ResourceType: "resource",
		},
	}

	// act
	got := modelhelpers.AuditEventsToDomain(have)

	// assert
	want := []*entities.AuditEvent{
		{
			ID:           id,
			Operation:    entities.OperationCreated,
			CreatedAt:    now,
			Summary:      testSummary,
			ResourceID:   resourceID,
			ResourceType: entities.ResourceTypeResource,
		},
	}

	assert.Equal(t, want, got)
}

func Test_AuditEventToDomain(t *testing.T) {
	// arrange
	id := uuid.New()
	now := time.Now()
	resourceID := uuid.New()

	have := &model.AuditEvents{
		ID:           id,
		Operation:    "created",
		CreatedAt:    now,
		Summary:      testSummary,
		ResourceID:   resourceID,
		ResourceType: "resource",
	}

	// act
	got := modelhelpers.AuditEventToDomain(have)

	// assert
	want := &entities.AuditEvent{
		ID:           id,
		Operation:    entities.OperationCreated,
		CreatedAt:    now,
		Summary:      testSummary,
		ResourceID:   resourceID,
		ResourceType: entities.ResourceTypeResource,
	}

	assert.Equal(t, want, got)
}

func Test_AuditEventFromDomain(t *testing.T) {
	// arrange
	id := uuid.New()
	now := time.Now()
	resourceID := uuid.New()

	have := &entities.AuditEvent{
		ID:           id,
		Operation:    entities.OperationCreated,
		CreatedAt:    now,
		Summary:      testSummary,
		ResourceID:   resourceID,
		ResourceType: entities.ResourceTypeResource,
	}

	// act
	got := modelhelpers.AuditEventFromDomain(have)

	// assert
	want := &model.AuditEvents{
		ID:           id,
		Operation:    "created",
		CreatedAt:    now,
		Summary:      testSummary,
		ResourceID:   resourceID,
		ResourceType: "resource",
	}

	assert.Equal(t, want, got)
}
