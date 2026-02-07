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

// CreateResource provides the business logic for creating a resource.
type CreateResource struct {
	clock   repositories.Clock
	uuidgen repositories.UUIDGenerator
}

// NewCreateResource creates a new CreateResource.
func NewCreateResource(
	clock repositories.Clock,
	uuidgen repositories.UUIDGenerator,
) *CreateResource {
	return &CreateResource{
		clock:   clock,
		uuidgen: uuidgen,
	}
}

// Execute creates a resource.
//
// If the resource name is already in use, ErrAlreadyExists is returned. Any other failure will
// cause ErrInternal to be returned.
func (u *CreateResource) Execute(
	ctx context.Context,
	logger *slog.Logger,
	store repositories.Resource,
	resource *entities.Resource,
) (*entities.Resource, error) {
	id, err := u.uuidgen.NewV7()
	if err != nil {
		logger.ErrorContext(ctx, "failed to generate id", slogx.Err(err))
		return nil, domain.ErrInternal
	}

	logger = logger.With(slogx.ResourceID(id))
	resource.Init(id, u.clock.UTCNow())

	if err := store.CreateResource(ctx, resource); err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			logger.InfoContext(ctx, "resource exists", slogx.Err(err))
			return nil, domain.ErrAlreadyExists
		}

		logger.ErrorContext(ctx, "failed to create resource", slogx.Err(err))
		return nil, domain.ErrInternal
	}

	logger.InfoContext(ctx, "created resource")
	return resource, nil
}
