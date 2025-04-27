package usecases_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/usecases"
	"github.com/mattdowdell/sandbox/mocks/domain/mockrepositories"
)

func Test_NewListAuditEvents(t *testing.T) {
	// no arrange necessary

	// act
	usecase := usecases.NewListAuditEvents()

	// assert
	assert.NotNil(t, usecase)
}

func Test_ListAuditEvents_Success(t *testing.T) {
	// arrange
	usecase := usecases.NewListAuditEvents()

	expected := []*entities.AuditEvent{
		{
			ID:           uuid.New(),
			Operation:    entities.OperationCreated,
			CreatedAt:    time.Now(),
			Summary:      "example",
			ResourceID:   uuid.New(),
			ResourceType: entities.ResourceTypeResource,
		},
	}

	store := mockrepositories.NewAuditEvent(t)
	store.EXPECT().ListAuditEvents(t.Context()).Return(expected, nil).Once()

	// act
	events, err := usecase.Execute(t.Context(), store)

	// assert
	assert.Equal(t, expected, events)
	assert.NoError(t, err)
}

func Test_ListAuditEvents_Error(t *testing.T) {
	// arrange
	usecase := usecases.NewListAuditEvents()

	store := mockrepositories.NewAuditEvent(t)
	store.EXPECT().ListAuditEvents(t.Context()).Return(nil, domain.ErrInternal).Once()

	// act
	events, err := usecase.Execute(t.Context(), store)

	// assert
	assert.Empty(t, events)
	assert.ErrorIs(t, err, domain.ErrInternal)
}
