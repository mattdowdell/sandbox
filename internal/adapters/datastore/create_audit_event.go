package datastore

import (
	"context"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore/modelhelpers"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/postgres/public/table"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
func (d *Datastore) CreateAuditEvent(ctx context.Context, event *entities.AuditEvent) error {
	m := modelhelpers.AuditEventFromDomain(event)

	stmt := table.AuditEvents.
		INSERT(
			table.AuditEvents.ID,
			table.AuditEvents.Operation,
			table.AuditEvents.CreatedAt,
			table.AuditEvents.Summary,
			table.AuditEvents.ResourceID,
			table.AuditEvents.ResourceType,
		).
		MODEL(m)

	if _, err := stmt.ExecContext(ctx, d.db); err != nil {
		return err
	}

	return nil
}
