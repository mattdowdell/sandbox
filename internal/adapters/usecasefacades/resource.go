package usecasefacades

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/mattdowdell/sandbox/internal/adapters/common"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
type Resource struct {
	provider common.Provider
	creator  ResourceCreator
	getter   ResourceGetter
	lister   ResourceLister
	updater  ResourceUpdater
	deleter  ResourceDeleter
}

// ...
func NewResource(
	provider common.Provider,
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
	return common.TxValue(ctx, logger, r.provider, func(ds common.Datastore) (*entities.Resource, error) {
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
) ([]*entities.Resource, error) {
	return r.lister.Execute(ctx, logger, r.provider.Datastore())
}

// ...
func (r *Resource) Update(
	ctx context.Context,
	logger *slog.Logger,
	input *entities.Resource,
) (*entities.Resource, error) {
	return common.TxValue(ctx, logger, r.provider, func(ds common.Datastore) (*entities.Resource, error) {
		return r.updater.Execute(ctx, logger, ds, input)
	})
}

// ...
func (r *Resource) Delete(
	ctx context.Context,
	logger *slog.Logger,
	input uuid.UUID,
) error {
	return common.TxFunc(ctx, logger, r.provider, func(ds common.Datastore) error {
		return r.deleter.Execute(ctx, logger, ds, input)
	})
}
