package modelhelpers

import (
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/postgres/public/model"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
func AuditEventsToDomain(inputs []model.AuditEvents) []*entities.AuditEvent {
	outputs := make([]*entities.AuditEvent, 0, len(inputs))

	for _, input := range inputs {
		outputs = append(outputs, AuditEventToDomain(input))
	}

	return outputs
}

// ...
func AuditEventToDomain(input model.AuditEvents) *entities.AuditEvent {
	return &entities.AuditEvent{
		ID:           input.ID,
		Operation:    entities.ParseOperation(input.Operation),
		CreatedAt:    input.CreatedAt,
		Summary:      input.Summary,
		ResourceID:   input.ResourceID,
		ResourceType: entities.ParseResourceType(input.ResourceType),
	}
}

// ...
func AuditEventFromDomain(input *entities.AuditEvent) model.AuditEvents {
	return model.AuditEvents{
		ID:           input.ID,
		Operation:    input.Operation.String(),
		CreatedAt:    input.CreatedAt,
		Summary:      input.Summary,
		ResourceID:   input.ResourceID,
		ResourceType: input.ResourceType.String(),
	}
}
