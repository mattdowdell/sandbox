package usecases

import (
	"context"
	"errors"
	"log/slog"

	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// UpdateResource provides the business logic for updating a resource.
type UpdateResource struct {
	clock repositories.Clock
}

// NewUpdateResource creates a new UpdateResource.
func NewUpdateResource(
	clock repositories.Clock,
) *UpdateResource {
	return &UpdateResource{
		clock: clock,
	}
}

// Execute gets a single resource.
//
// If the resource ID does not exist, ErrNotFound is returned. If the resource name is already in
// use, ErrAlreadyExists is returned. Any other failure will cause ErrInternal to be returned.
func (u *UpdateResource) Execute(
	ctx context.Context,
	logger *slog.Logger,
	store repositories.Resource,
	changes *entities.Resource,
) (*entities.Resource, error) {
	changes.Update(u.clock.Now())

	resource, err := store.UpdateResource(ctx, changes)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			logger.InfoContext(ctx, "resource not found", slogx.Err(err))
			return nil, domain.ErrNotFound

		case errors.Is(err, domain.ErrAlreadyExists):
			logger.InfoContext(ctx, "resource exists", slogx.Err(err))
			return nil, domain.ErrAlreadyExists

		default:
			logger.ErrorContext(ctx, "failed to update resource", slogx.Err(err))
			return nil, domain.ErrInternal
		}
	}

	logger.InfoContext(ctx, "updated resource")
	return resource, nil
}
