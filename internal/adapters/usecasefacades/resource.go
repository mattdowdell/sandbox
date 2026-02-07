package usecasefacades

import (
	"context"
	"log/slog"

	"github.com/gofrs/uuid/v5"

	"github.com/mattdowdell/sandbox/internal/adapters/txn"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
)

// ...
type Resource struct {
	provider txn.Provider
	creator  ResourceCreator
	getter   ResourceGetter
	lister   ResourceLister
	updater  ResourceUpdater
	deleter  ResourceDeleter
}

// ...
func NewResource(
	provider txn.Provider,
	creator ResourceCreator,
	getter ResourceGetter,
	lister ResourceLister,
	updater ResourceUpdater,
	deleter ResourceDeleter,
) *Resource {
	return &Resource{
		provider: provider,
		creator:  creator,
		getter:   getter,
		lister:   lister,
		updater:  updater,
		deleter:  deleter,
	}
}

// ...
func (r *Resource) Create(
	ctx context.Context,
	logger *slog.Logger,
	input *entities.Resource,
) (*entities.Resource, error) {
	return txn.Value(ctx, logger, r.provider, func(ds txn.Datastore) (*entities.Resource, error) {
		return r.creator.Execute(ctx, logger, ds, input)
	})
}

// ...
func (r *Resource) Get(
	ctx context.Context,
	logger *slog.Logger,
	input uuid.UUID,
) (*entities.Resource, error) {
	return r.getter.Execute(ctx, logger, r.provider.Datastore(), input)
}

// ...
func (r *Resource) List(
	ctx context.Context,
	logger *slog.Logger,
	pager repositories.Pager,
) (*repositories.Paged[*entities.Resource], error) {
	return r.lister.Execute(ctx, logger, r.provider.Datastore(), pager)
}

// ...
func (r *Resource) Update(
	ctx context.Context,
	logger *slog.Logger,
	input *entities.Resource,
) (*entities.Resource, error) {
	return txn.Value(ctx, logger, r.provider, func(ds txn.Datastore) (*entities.Resource, error) {
		return r.updater.Execute(ctx, logger, ds, input)
	})
}

// ...
func (r *Resource) Delete(
	ctx context.Context,
	logger *slog.Logger,
	input uuid.UUID,
) error {
	return txn.Func(ctx, logger, r.provider, func(ds txn.Datastore) error {
		return r.deleter.Execute(ctx, logger, ds, input)
	})
}
