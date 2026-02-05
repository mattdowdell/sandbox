-- Initial database migration.

-- For storing resource data.
CREATE TABLE resources (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	created_at TIMESTAMPTZ (0) NOT NULL,
	updated_at TIMESTAMPTZ (0) NOT NULL
);

-- For storing audit events for the service.
--
-- The resource ID and type reference the resource the audit event targets. It is assumed that these
-- values uniquely identify any resource in the system.
CREATE TABLE audit_events (
	id UUID PRIMARY KEY,
	operation TEXT NOT NULL,
	created_at TIMESTAMPTZ (0) NOT NULL,
	summary TEXT NOT NULL,
	resource_id UUID NOT NULL,
	resource_type TEXT NOT NULL
);
