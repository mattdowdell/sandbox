package models

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
var (
	ErrParseResourceID = errors.New("failed to parse resource id")
)

// ...
func ResourceCreateToDomain(create *examplev1.ResourceCreate) *entities.Resource {
	return &entities.Resource{
		Name: create.GetName(),
	}
}

// ...
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

// ...
func ResourcesFromDomain(input []*entities.Resource) []*examplev1.Resource {
	output := make([]*examplev1.Resource, 0, len(input))

	for i := range input {
		output = append(output, ResourceFromDomain(input[i]))
	}

	return output
}

// ...
func ResourceFromDomain(input *entities.Resource) *examplev1.Resource {
	return &examplev1.Resource{
		Id:        input.ID.String(),
		Name:      input.Name,
		CreatedAt: timestamppb.New(input.CreatedAt),
		UpdatedAt: timestamppb.New(input.UpdatedAt),
	}
}

// ...
func ParseID(msg interface{ GetId() string }) (uuid.UUID, error) {
	id, err := uuid.Parse(msg.GetId())
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w %q: %w", ErrParseResourceID, msg.GetId(), err)
	}

	return id, nil
}
