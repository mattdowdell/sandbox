package usecases

import (
	"context"
	"log/slog"

	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ...
type WatchAuditEvents struct {
	store repositories.AuditEvent
}

// ...
func NewWatchAuditEvents(
	store repositories.AuditEvent,
) *WatchAuditEvents {
	return &WatchAuditEvents{
		store: store,
	}
}

// ...
func (u *WatchAuditEvents) Execute(ctx context.Context) <-chan *entities.AuditEvent {
	ch := make(chan *entities.AuditEvent, 1)

	go func() {
		if err := u.store.WatchAuditEvents(ctx, ch); err != nil {
			slog.ErrorContext(ctx, "failed to watch audit events", slogx.Err(err))
		}
	}()

	return ch
}
