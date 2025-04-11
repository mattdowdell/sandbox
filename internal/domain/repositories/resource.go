package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// Resource implementations support CRUDL operations for a resource.
type Resource interface {
	// CreateResource creates a resource.
	//
	// If the resource name is already in use, ErrAlreadyExists shall be returned. Any other
	// failure shall cause ErrInternal to be returned. These errors must be wrapped with additional
	// context to enable debugging.
	CreateResource(context.Context, *entities.Resource) error

	// GetResource gets a single resource.
	//
	// If the resource ID does not exist, ErrNotFound shall be returned. Any other failure shall
	// cause ErrInternal to be returned. These errors must be wrapped with additional context to
	// enable debugging.
	GetResource(context.Context, uuid.UUID) (*entities.Resource, error)

	// ListResources lists multiple resources.
	//
	// Any failure shall cause ErrInternal to be returned. This error must be wrapped with
	// additional context to enable debugging.
	ListResources(context.Context) ([]*entities.Resource, error)

	// UpdateResource updates a resource and returns the latest version. Only updates to the Name
	// and UpdatedAt fields are permitted.
	//
	// If the resource ID does not exist, ErrNotFound shall be returned. If the resource name is
	// already in use, ErrAlreadyExists shall be returned. Any other failure shall cause ErrInternal
	// to be returned. These errors must be wrapped with additional context to enable debugging.
	UpdateResource(context.Context, *entities.Resource) (*entities.Resource, error)

	// DeleteResource deletes a resource.
	//
	// If the resource ID does not exist, ErrNotFound shall be returned. Any other failure shall
	// cause ErrInternal to be returned. These errors must be wrapped with additional context to
	// enable debugging.
	DeleteResource(context.Context, uuid.UUID) error
}
