package usecasefacades

import (
	"context"
	"log/slog"

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
	logger *slog.Logger,
) ([]*entities.AuditEvent, error) {
	return a.lister.Execute(ctx, logger, a.provider.Datastore())
}

// ...
func (a *AuditEvent) Watch(
	ctx context.Context,
	logger *slog.Logger,
) <-chan *entities.AuditEvent {
	return a.watcher.Execute(ctx, logger, a.provider.Datastore())
}
