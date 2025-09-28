-- TODO: document

-- TODO: document
CREATE TABLE herd_system_migrations (
	id UUID PRIMARY KEY,
	migrated_at TIMESTAMPTZ (0) NOT NULL,
	migration_version TEXT NOT NULL,
	code_version TEXT NOT NULL,
);

-- TODO: document
CREATE TABLE herd_user_migrations (
	id UUID PRIMARY KEY,
	migrated_at TIMESTAMPTZ (0) NOT NULL,
	migration_version TEXT NOT NULL,
	code_version TEXT NOT NULL,
);
