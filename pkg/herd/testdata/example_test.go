package herd_test

import (
	"embed"

	"github.com/mattdowdell/sandbox/pkg/herd"
)

//go:embed testdata/*.sql
var migrationFS embed.FS

func Example() {
	migrations, _ := herd.CollectFileMigrations(migrationFS)

	// TODO: use migrations
	_ = migrations
}
