package usecasefacades

import (
	"context"

	"github.com/mattdowdell/sandbox/internal/adapters/common"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
type AuditEvent struct {
	provider common.Provider
	lister   AuditEventLister
	watcher  AuditEventWatcher
}

// ...
func NewAuditEvent(
	provider common.Provider,
	lister AuditEventLister,
	watcher AuditEventWatcher,
) *AuditEvent {
	return &AuditEvent{
		provider: provider,
		lister:   lister,
		watcher:  watcher,
	}
}

// ...
func (a *AuditEvent) List(
	ctx context.Context,
) ([]*entities.AuditEvent, error) {
	return a.lister.Execute(ctx, a.provider.Datastore())
}

// ...
func (a *AuditEvent) Watch(
	ctx context.Context,
) <-chan *entities.AuditEvent {
	return a.watcher.Execute(ctx, a.provider.Datastore())
}
