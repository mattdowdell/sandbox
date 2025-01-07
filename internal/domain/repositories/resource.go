package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
type Resource interface {
	// ...
	CreateResource(context.Context, *entities.Resource) error

	// ...
	GetResource(context.Context, uuid.UUID) (*entities.Resource, error)

	// ...
	ListResources(context.Context) ([]*entities.Resource, error)

	// ...
	UpdateResource(context.Context, *entities.Resource) error

	// ...
	DeleteResource(context.Context, uuid.UUID) error
}
