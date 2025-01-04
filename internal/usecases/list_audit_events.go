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
type ListAuditEvents struct {
	store repositories.AuditEvent
}

// ...
func NewListAuditEvents(
	store repositories.AuditEvent,
) *ListAuditEvents {
	return &ListAuditEvents{
		store: store,
	}
}

// ...
func (u *ListAuditEvents) Execute(
	ctx context.Context,
) ([]*entities.AuditEvent, error) {
	events, err := u.store.ListAuditEvents(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list audit events", slogx.Err(err))
		return nil, errors.New("internal error")
	}

	return events, nil
}
