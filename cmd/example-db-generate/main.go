package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/go-jet/jet/v2/generator/postgres"

	"github.com/mattdowdell/sandbox/internal/drivers/config/secret"
	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/internal/drivers/pgsql"
	"github.com/mattdowdell/sandbox/pkg/exit"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

const (
	defaultHost    = "localhost"
	defaultPort    = 5432
	defaultName    = "postgres"
	defaultSchema  = "public"
	defaultSSLMode = "disable"
)

var defaultOutput = filepath.Join("internal", "adapters", "datastore", "schema")

var (
	output   = flag.String("output", defaultOutput, "The directory to output to")
	host     = flag.String("host", defaultHost, "The hostname of the database server")
	port     = flag.Int("port", defaultPort, "The port the database server is lietening on")
	username = flag.String("username", "", "The username to authenticate with")
	password = flag.String("password", "", "The password to authenticate with")
	name     = flag.String("name", defaultName, "The name of the database")
	schema   = flag.String("schema", defaultSchema, "The database schema to use")
	sslmode  = flag.String("sslmode", defaultSSLMode, "The SSL mode to connect with")
)

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	flag.Parse()

	logger := logging.New(slog.LevelInfo)

	db, _, err := pgsql.NewFromConfig(ctx, pgsql.Config{
		Hostname: *host,
		Port:     *port,
		Username: *username,
		Password: secret.String(*password),
		Name:     *name,
		SSLMode:  *sslmode,
	})
	if err != nil {
		logger.ErrorContext(ctx, "failed to create db connection pool", slogx.Err(err))
		return exit.Failure
	}

	if err := postgres.GenerateDB(db, *schema, *output, updateTemplate()); err != nil {
		logger.ErrorContext(ctx, "failed to generate database models", slogx.Err(err))
		return exit.Failure
	}

	return exit.Success
}
