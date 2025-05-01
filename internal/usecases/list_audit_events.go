package usecases

import (
	"context"
	"log/slog"

	"github.com/mattdowdell/sandbox/internal/domain/apperrors"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ListResources provides the business logic for listing multiple audit events.
type ListAuditEvents struct{}

// NewListResources creates a new ListResources.
func NewListAuditEvents() *ListAuditEvents {
	return &ListAuditEvents{}
}

// Execute lists multiple audit events.
//
// Any failure will cause ErrInternal to be returned.
func (u *ListAuditEvents) Execute(
	ctx context.Context,
	logger *slog.Logger,
	store repositories.AuditEvent,
) ([]*entities.AuditEvent, error) {
	events, err := store.ListAuditEvents(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "failed to list audit events", slogx.Err(err))
		return nil, apperrors.ErrInternal
	}

	return events, nil
}
