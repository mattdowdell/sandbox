package usecases

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// GetResource provides the business logic for getting a single resource.
type GetResource struct{}

// NewGetResource creates a new GetResource.
func NewGetResource() *GetResource {
	return &GetResource{}
}

// Execute gets a single resource.
//
// If the resource ID does not exist, ErrNotFound is returned. Any other failure will cause
// ErrInternal to be returned.
func (u *GetResource) Execute(
	ctx context.Context,
	logger *slog.Logger,
	store repositories.Resource,
	id uuid.UUID,
) (*entities.Resource, error) {
	resource, err := store.GetResource(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			logger.InfoContext(ctx, "resource not found", slogx.Err(err))
			return nil, domain.ErrNotFound
		}

		logger.ErrorContext(ctx, "failed to get resource", slogx.Err(err))
		return nil, domain.ErrInternal
	}

	return resource, nil
}
