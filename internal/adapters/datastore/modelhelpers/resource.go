package modelhelpers

import (
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/public/model"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ResourcesToDomain converts multiple database resource models to the equivalent domain
// representation.
func ResourcesToDomain(inputs []*model.Resources) []*entities.Resource {
	outputs := make([]*entities.Resource, 0, len(inputs))

	for _, input := range inputs {
		outputs = append(outputs, ResourceToDomain(input))
	}

	return outputs
}

// ResourceToDomain converts a database resource model to the equivalent domain representation.
func ResourceToDomain(input *model.Resources) *entities.Resource {
	return &entities.Resource{
		ID:        toGofrsUUID(input.ID),
		Name:      input.Name,
		CreatedAt: input.CreatedAt,
		UpdatedAt: input.UpdatedAt,
	}
}

// ResourceToDomain converts a domain resource model to the equivalent database representation.
func ResourceFromDomain(input *entities.Resource) *model.Resources {
	return &model.Resources{
		ID:        toGoogleUUID(input.ID),
		Name:      input.Name,
		CreatedAt: input.CreatedAt,
		UpdatedAt: input.UpdatedAt,
	}
}
