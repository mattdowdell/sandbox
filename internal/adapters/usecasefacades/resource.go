package usecasefacades

import (
	"context"

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
	input *entities.Resource,
) (*entities.Resource, error) {
	return common.TxValue(ctx, r.provider, func(ds common.Datastore) (*entities.Resource, error) {
		return r.creator.Execute(ctx, ds, input)
	})
}

// ...
func (r *Resource) Get(
	ctx context.Context,
	input uuid.UUID,
) (*entities.Resource, error) {
	return r.getter.Execute(ctx, r.provider.Datastore(), input)
}

// ...
func (r *Resource) List(
	ctx context.Context,
) ([]*entities.Resource, error) {
	return r.lister.Execute(ctx, r.provider.Datastore())
}

// ...
func (r *Resource) Update(
	ctx context.Context,
	input *entities.Resource,
) (*entities.Resource, error) {
	return common.TxValue(ctx, r.provider, func(ds common.Datastore) (*entities.Resource, error) {
		return r.updater.Execute(ctx, ds, input)
	})
}

// ...
func (r *Resource) Delete(
	ctx context.Context,
	input uuid.UUID,
) error {
	return common.TxFunc(ctx, r.provider, func(ds common.Datastore) error {
		return r.deleter.Execute(ctx, ds, input)
	})
}
