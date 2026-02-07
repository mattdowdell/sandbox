package models_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	examplev1 "github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc/models"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

const (
	testSummary = "example"
)

func Test_AuditEventsFromDomain(t *testing.T) {
	// arrange
	eventID := uuid.Must(uuid.NewV7())
	resourceID := uuid.Must(uuid.NewV7())
	now := time.Now()

	input := []*entities.AuditEvent{
		{
			ID:           eventID,
			Operation:    entities.OperationCreated,
			CreatedAt:    now,
			Summary:      testSummary,
			ResourceID:   resourceID,
			ResourceType: entities.ResourceTypeResource,
		},
	}

	// act
	output := models.AuditEventsFromDomain(input)

	// assert
	expected := []*examplev1.AuditEvent{
		{
			Id:           eventID.String(),
			Operation:    examplev1.Operation_OPERATION_CREATED,
			CreatedAt:    timestamppb.New(now),
			Summary:      testSummary,
			ResourceId:   resourceID.String(),
			ResourceType: "example.v1.Resource",
		},
	}

	assert.Equal(t, expected, output)
}

func Test_AuditEventFromDomain(t *testing.T) {
	eventID := uuid.Must(uuid.NewV7())
	resourceID := uuid.Must(uuid.NewV7())
	now := time.Now()

	tests := map[string]struct {
		input    *entities.AuditEvent
		expected *examplev1.AuditEvent
	}{
		"created resource": {
			input: &entities.AuditEvent{
				ID:           eventID,
				Operation:    entities.OperationCreated,
				CreatedAt:    now,
				Summary:      testSummary,
				ResourceID:   resourceID,
				ResourceType: entities.ResourceTypeResource,
			},
			expected: &examplev1.AuditEvent{
				Id:           eventID.String(),
				Operation:    examplev1.Operation_OPERATION_CREATED,
				CreatedAt:    timestamppb.New(now),
				Summary:      testSummary,
				ResourceId:   resourceID.String(),
				ResourceType: "example.v1.Resource",
			},
		},
		"modified resource": {
			input: &entities.AuditEvent{
				ID:           eventID,
				Operation:    entities.OperationModified,
				CreatedAt:    now,
				Summary:      testSummary,
				ResourceID:   resourceID,
				ResourceType: entities.ResourceTypeResource,
			},
			expected: &examplev1.AuditEvent{
				Id:           eventID.String(),
				Operation:    examplev1.Operation_OPERATION_MODIFIED,
				CreatedAt:    timestamppb.New(now),
				Summary:      testSummary,
				ResourceId:   resourceID.String(),
				ResourceType: "example.v1.Resource",
			},
		},
		"deleted resource": {
			input: &entities.AuditEvent{
				ID:           eventID,
				Operation:    entities.OperationDeleted,
				CreatedAt:    now,
				Summary:      testSummary,
				ResourceID:   resourceID,
				ResourceType: entities.ResourceTypeResource,
			},
			expected: &examplev1.AuditEvent{
				Id:           eventID.String(),
				Operation:    examplev1.Operation_OPERATION_DELETED,
				CreatedAt:    timestamppb.New(now),
				Summary:      testSummary,
				ResourceId:   resourceID.String(),
				ResourceType: "example.v1.Resource",
			},
		},
		"unknown resource": {
			input: &entities.AuditEvent{
				ID:           eventID,
				Operation:    0,
				CreatedAt:    now,
				Summary:      testSummary,
				ResourceID:   resourceID,
				ResourceType: entities.ResourceTypeResource,
			},
			expected: &examplev1.AuditEvent{
				Id:           eventID.String(),
				Operation:    examplev1.Operation_OPERATION_UNSPECIFIED,
				CreatedAt:    timestamppb.New(now),
				Summary:      testSummary,
				ResourceId:   resourceID.String(),
				ResourceType: "example.v1.Resource",
			},
		},
		"modified unknown": {
			input: &entities.AuditEvent{
				ID:           eventID,
				Operation:    entities.OperationModified,
				CreatedAt:    now,
				Summary:      testSummary,
				ResourceID:   resourceID,
				ResourceType: 0,
			},
			expected: &examplev1.AuditEvent{
				Id:           eventID.String(),
				Operation:    examplev1.Operation_OPERATION_MODIFIED,
				CreatedAt:    timestamppb.New(now),
				Summary:      testSummary,
				ResourceId:   resourceID.String(),
				ResourceType: "unknown",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			output := models.AuditEventFromDomain(tt.input)
			assert.Equal(t, tt.expected, output)
		})
	}
}
