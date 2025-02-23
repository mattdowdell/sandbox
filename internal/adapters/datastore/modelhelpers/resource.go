package modelhelpers

import (
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/postgres/public/model"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
func ResourcesToDomain(inputs []model.Resources) []*entities.Resource {
	outputs := make([]*entities.Resource, 0, len(inputs))

	for _, input := range inputs {
		outputs = append(outputs, ResourceToDomain(input))
	}

	return outputs
}

// ...
func ResourceToDomain(input model.Resources) *entities.Resource {
	return &entities.Resource{
		ID:        input.ID,
		Name:      input.Name,
		CreatedAt: input.CreatedAt,
		UpdatedAt: input.UpdatedAt,
	}
}

// ...
func ResourceFromDomain(input *entities.Resource) model.Resources {
	return model.Resources{
		ID:        input.ID,
		Name:      input.Name,
		CreatedAt: input.CreatedAt,
		UpdatedAt: input.UpdatedAt,
	}
}
