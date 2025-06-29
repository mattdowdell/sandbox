package datastore

import (
	"context"
	"fmt"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore/modelhelpers"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/public/model"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/public/table"
	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
func (d *Datastore) ListAuditEvents(ctx context.Context) ([]*entities.AuditEvent, error) {
	stmt := table.AuditEvents.
		SELECT(
			table.AuditEvents.ID,
			table.AuditEvents.Operation,
			table.AuditEvents.CreatedAt,
			table.AuditEvents.Summary,
			table.AuditEvents.ResourceID,
			table.AuditEvents.ResourceType,
		).
		ORDER_BY(table.Resources.ID.ASC())

	var events []*model.AuditEvents
	if err := stmt.QueryContext(ctx, d.db, &events); err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return modelhelpers.AuditEventsToDomain(events), nil
}
