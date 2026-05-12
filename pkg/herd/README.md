# Herd

<img width="200" alt="Herd Logo" align="right" src="https://github.com/user-attachments/assets/b17a9fe1-af36-40fe-a016-fa4beb9aa1a9" />

Herd is a library for applying migrations to a SQL database. Currently, only PostgreSQL is
supported, but other database engines may be added in the future.

At this time, only up-migrations are supported. This means that once a migration has been applied,
it cannot be un-done. Problematic migrations should therefore be corrected by creating an additional
migration.

## Migrations

Migrations must implement the `herd.Migration` interface, providing `Version` and `Migrate` methods.

```go
type MyMigration struct {}

// The unique version number for the migration. Migrations are applied in ascending order.
func (m *MyMigration) Version() int64 { /* ... */ }

// Migrate is called to appky the migration.
func (m *MyMigration) Migrate(ctx context.Context, db herd.DB) error { /* ... */ }
```

Migrations made up of SQL files only can use the built-in `herd.FileMigration`.

```go
migrations := []herd.Migration{
	herd.NewFileMigration(1, "-- your migration here"),
	herd.NewFileMigration(2, "-- your migration here"),
	herd.NewFileMigration(3, "-- your migration here"),
}
```

Files in a directory and named in the format `{version}_{name}.sql` can be automatically gathered.

```go
//go:embed migrations/*.sql
var migrationFS embed.FS

migrations, err := herd.CollectFileMigrations(migrationFS)
if err != nil {
	// handle error
}
```

Version numbers must be positive integers, but do not need to be consecutive. This is to allow date
or time based versions to be used if desired. However, a new version should not be inserted between
2 existing versions. Versions less than the largest applied version are skipped.

When more complex migration logic is required, custom implementations of `herd.Migration` can be
used. These can be used alongside `herd.FileMigration` instances. Each migration must continue to
have a unique version number, regardless of how it is implemented.

## Applying Migrations

Once the migrations have been gathered, they can be executed via `herd.Herd`. This determines which
migrations have not yet been applied, and then executes them within a transaction. As PostgreSQL
supports schema changes within a transaction, this means the pending migrations are either all
applied, or none of them are.

On success, the version before and after the applied migrations is returned.

```go
herder := herd.New(migrations)

result, err := herder.Migrate(ctx, db)
if err != nil {
	// handle error
}

fmt.Println(result.Before, "->", result.After)
```

To migrate to a specific version and stop, use `herd.WithTargetVersion`. This causes any pending
migrations with a version greater than the target version to be skipped, such that they can be
applied at a later time.

```go
herder := herd.New(migrations, herd.WithTargetVersion(100))

result, err := herder.Migrate(ctx, db)
if err != nil {
	// handle error
}

fmt.Println(result.Before, "->", result.After)
```

## Migration Metadata

Herd uses 2 database tables for tracking migrations: `herd_system_migrations` and
`herd_user_migrations`. `herd_system_migrations` is for tracking the schema changes to these 2
tables, while `herd_user_migrations` is for tracking the application of the provided migrations.

```sql
-- TODO: add example for each table
```

The primary use of these tables is for auditing, track when each migration was applied.
Additionally, they track the version and commit of the application that executed them. This refers
to the application that is importing and using Herd, not Herd itself.


