package datastore

import (
	"context"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore/modelhelpers"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/postgres/public/model"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/postgres/public/table"
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

	var events []model.AuditEvents
	if err := stmt.QueryContext(ctx, d.db, &events); err != nil {
		return nil, err
	}

	return modelhelpers.AuditEventsToDomain(events), nil
}
