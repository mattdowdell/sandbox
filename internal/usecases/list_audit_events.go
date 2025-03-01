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
type ListAuditEvents struct {}

// ...
func NewListAuditEvents() *ListAuditEvents {
	return &ListAuditEvents{}
}

// ...
func (u *ListAuditEvents) Execute(
	ctx context.Context,
	store repositories.AuditEvent,
) ([]*entities.AuditEvent, error) {
	events, err := store.ListAuditEvents(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list audit events", slogx.Err(err))
		return nil, errors.New("internal error")
	}

	return events, nil
}
