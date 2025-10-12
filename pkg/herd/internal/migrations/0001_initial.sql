-- Initial migration for recording which migrations have been applied via herd.

-- Create a table for recording internal/system migrations. This is mostly a way to track what
-- schema version is in use.
CREATE TABLE herd_system_migrations (
	migration_version BIGINT PRIMARY KEY,
	migrated_at TIMESTAMPTZ (0) NOT NULL,
	code_version TEXT NOT NULL,
	code_revision TEXT NOT NULL,
);

-- Create a table for recording user-defined migrations. Also include the herd_version which is the
-- highest herd_system_migrations.migration_version value.
CREATE TABLE herd_user_migrations (
	migration_version BIGINT PRIMARY KEY,
	migrated_at TIMESTAMPTZ (0) NOT NULL,
	code_version TEXT NOT NULL,
	code_revision TEXT NOT NULL,
	herd_version BIGINT NOT NULL,
);
