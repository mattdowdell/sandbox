package usecases

import (
	"context"
	"errors"
	"log/slog"

	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ...
type CreateResource struct {
	clock   repositories.Clock
	uuidgen repositories.UUIDGenerator
}

// ...
func NewCreateResource(
	clock repositories.Clock,
	uuidgen repositories.UUIDGenerator,
) *CreateResource {
	return &CreateResource{
		clock:   clock,
		uuidgen: uuidgen,
	}
}

// ...
func (u *CreateResource) Execute(
	ctx context.Context,
	store repositories.Resource,
	resource *entities.Resource,
) (*entities.Resource, error) {
	id, err := u.uuidgen.NewV7()
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate id", slogx.Err(err))
		return nil, errors.New("internal error")
	}

	resource.Init(id, u.clock.Now())

	// TODO: handle conflict
	if err := store.CreateResource(ctx, resource); err != nil {
		slog.ErrorContext(ctx, "failed to create resource", slogx.Err(err))
		return nil, errors.New("internal error")
	}

	slog.InfoContext(ctx, "created resource")

	return resource, nil
}
