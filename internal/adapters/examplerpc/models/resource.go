package models

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ResourceCreateToDomain converts a ResourceCreate message into the equivalent domain
// representation.
func ResourceCreateToDomain(create *examplev1.ResourceCreate) *entities.Resource {
	return &entities.Resource{
		Name: create.GetName(),
	}
}

// ResourceUpdateToDomain converts a ResourceUpdate message into the equivalent domain
// representation.
func ResourceUpdateToDomain(update *examplev1.ResourceUpdate) (*entities.Resource, error) {
	id, err := ParseID(update)
	if err != nil {
		return nil, err
	}

	return &entities.Resource{
		ID:   id,
		Name: update.GetName(),
	}, nil
}

// ResourcesFromDomain converts multiple resources into the equivalent Protobuf messages.
func ResourcesFromDomain(input []*entities.Resource) []*examplev1.Resource {
	output := make([]*examplev1.Resource, 0, len(input))

	for i := range input {
		output = append(output, ResourceFromDomain(input[i]))
	}

	return output
}

// ResourcesFromDomain converts a resource into the equivalent Protobuf message.
func ResourceFromDomain(input *entities.Resource) *examplev1.Resource {
	return &examplev1.Resource{
		Id:        input.ID.String(),
		Name:      input.Name,
		CreatedAt: timestamppb.New(input.CreatedAt),
		UpdatedAt: timestamppb.New(input.UpdatedAt),
	}
}
