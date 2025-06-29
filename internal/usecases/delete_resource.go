package usecases

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gofrs/uuid/v5"

	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// DeleteResource provides the business logic for deleting a resource.
type DeleteResource struct{}

// NewDeleteResource creates a new DeleteResource.
func NewDeleteResource() *DeleteResource {
	return &DeleteResource{}
}

// Execute deletes a resource.
//
// If the resource ID does not exist, ErrNotFound is returned. Any other failure will cause
// ErrInternal to be returned.
func (u *DeleteResource) Execute(
	ctx context.Context,
	logger *slog.Logger,
	store repositories.Resource,
	id uuid.UUID,
) error {
	if err := store.DeleteResource(ctx, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			logger.InfoContext(ctx, "resource not found", slogx.Err(err))
			return domain.ErrNotFound
		}

		logger.ErrorContext(ctx, "failed to delete resource", slogx.Err(err))
		return domain.ErrInternal
	}

	logger.InfoContext(ctx, "deleted resource")
	return nil
}
