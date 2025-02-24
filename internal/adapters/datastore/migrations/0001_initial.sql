-- TODO: document

-- For storing resource data.
CREATE TABLE resources (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	created_at TIMESTAMPTZ (0) NOT NULL,
	updated_at TIMESTAMPTZ (0) NOT NULL
);

-- For storing audit events for the service.
--
-- TODO: document use for resource id/type.
CREATE TABLE audit_events (
	id UUID PRIMARY KEY,
	operation TEXT NOT NULL,
	created_at TIMESTAMPTZ (0) NOT NULL,
	summary TEXT NOT NULL,
	resource_id UUID NOT NULL,
	resource_type TEXT NOT NULL
);
