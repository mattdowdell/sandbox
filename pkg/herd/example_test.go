package herd_test

import (
	"context"
	"database/sql"
	"embed"
	"log/slog"

	"github.com/mattdowdell/sandbox/pkg/herd"
)

//go:embed testdata/migrations/*.sql
var migrationFS embed.FS

func Example() {
	ctx := context.Background()

	migrations, _ := herd.CollectFileMigrationsFromFS(migrationFS)
	migrator, _ := herd.New(migrations)

	db, _ := sql.Open("example", "dsn")

	result, err := migrator.Migrate(ctx, db)
	if err != nil {
		slog.Error("failed to run migrate", slog.Any("error", err))
		return
	}

	slog.Info("migration completed", slog.Any("result", result))
}
