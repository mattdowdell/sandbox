package models

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ResourcesFromDomain converts multiple audit events into the equivalent Protobuf messages.
func AuditEventsFromDomain(input []*entities.AuditEvent) []*examplev1.AuditEvent {
	output := make([]*examplev1.AuditEvent, 0, len(input))

	for i := range input {
		output = append(output, AuditEventFromDomain(input[i]))
	}

	return output
}

// ResourcesFromDomain converts an audit event into the equivalent Protobuf message.
func AuditEventFromDomain(input *entities.AuditEvent) *examplev1.AuditEvent {
	return &examplev1.AuditEvent{
		Id:           input.ID.String(),
		Operation:    operationFromDomain(input.Operation),
		CreatedAt:    timestamppb.New(input.CreatedAt),
		Summary:      input.Summary,
		ResourceId:   input.ResourceID.String(),
		ResourceType: resourceTypeFromDomain(input.ResourceType),
	}
}

// operationFromDomain converts an operation into the equivalent Protobuf enum variant.
//
// Unknown values are not rejected under the assumption that invalid results will be rejected by
// Protobuf validation.
func operationFromDomain(input entities.Operation) examplev1.Operation {
	switch input {
	case entities.OperationCreated:
		return examplev1.Operation_OPERATION_CREATED

	case entities.OperationModified:
		return examplev1.Operation_OPERATION_MODIFIED

	case entities.OperationDeleted:
		return examplev1.Operation_OPERATION_DELETED

	default:
		return examplev1.Operation_OPERATION_UNSPECIFIED
	}
}

// resourceTypeFromDomain converts a resource type into the equivalent Protobuf enum variant.
//
// Unknown values are not rejected under the assumption that invalid results will be rejected by
// Protobuf validation.
func resourceTypeFromDomain(input entities.ResourceType) string {
	switch input {
	case entities.ResourceTypeResource:
		return string((&examplev1.Resource{}).ProtoReflect().Descriptor().FullName())

	default:
		return "unknown"
	}
}
