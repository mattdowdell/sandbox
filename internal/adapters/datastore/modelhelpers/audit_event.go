package modelhelpers

import (
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/public/model"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// AuditEventsToDomain converts multiple database audit event models to the equivalent domain
// representation.
func AuditEventsToDomain(inputs []*model.AuditEvents) []*entities.AuditEvent {
	outputs := make([]*entities.AuditEvent, 0, len(inputs))

	for _, input := range inputs {
		outputs = append(outputs, AuditEventToDomain(input))
	}

	return outputs
}

// AuditEventToDomain converts a database audit event model to the equivalent domain representation.
func AuditEventToDomain(input *model.AuditEvents) *entities.AuditEvent {
	return &entities.AuditEvent{
		ID:           input.ID,
		Operation:    entities.ParseOperation(input.Operation),
		CreatedAt:    input.CreatedAt,
		Summary:      input.Summary,
		ResourceID:   input.ResourceID,
		ResourceType: entities.ParseResourceType(input.ResourceType),
	}
}

// AuditEventFromDomain converts a domain audit event model to the equivalent database
// representation.
func AuditEventFromDomain(input *entities.AuditEvent) *model.AuditEvents {
	return &model.AuditEvents{
		ID:           input.ID,
		Operation:    input.Operation.String(),
		CreatedAt:    input.CreatedAt,
		Summary:      input.Summary,
		ResourceID:   input.ResourceID,
		ResourceType: input.ResourceType.String(),
	}
}
