package usecases

import (
	"context"
	"errors"
	"log/slog"

	"github.com/mattdowdell/sandbox/internal/domain/apperrors"
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
	store repositories.Resource,
	resource *entities.Resource,
) (*entities.Resource, error) {
	id, err := u.uuidgen.NewV7()
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate id", slogx.Err(err))

		return nil, apperrors.ErrInternal
	}

	resource.Init(id, u.clock.Now())

	if err := store.CreateResource(ctx, resource); err != nil {
		if errors.Is(err, apperrors.ErrAlreadyExists) {
			slog.InfoContext(ctx, "resource exists", slogx.Err(err))

			return nil, apperrors.ErrAlreadyExists
		}

		slog.ErrorContext(ctx, "failed to create resource", slogx.Err(err))

		return nil, apperrors.ErrInternal
	}

	slog.InfoContext(ctx, "created resource")
	return resource, nil
}
