package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/go-jet/jet/v2/generator/postgres"

	"github.com/mattdowdell/sandbox/internal/drivers/exit"
	"github.com/mattdowdell/sandbox/internal/drivers/pgsql"
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
	output   = flag.String("output", defaultOutput, "TODO")
	host     = flag.String("host", defaultHost, "TODO")
	port     = flag.Int("port", defaultPort, "TODO")
	username = flag.String("username", "", "TODO")
	password = flag.String("password", "", "TODO")
	name     = flag.String("name", defaultName, "TODO")
	schema   = flag.String("schema", defaultSchema, "TODO")
	sslmode  = flag.String("sslmode", defaultSSLMode, "TODO")
)

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	flag.Parse()

	db, err := pgsql.NewFromConfig(ctx, pgsql.Config{
		Hostname: *host,
		Port:     *port,
		Username: *username,
		Password: *password,
		Name:     *name,
		SSLMode:  *sslmode,
	})
	if err != nil {
		log.Print(err)
		return exit.Failure
	}

	if err := postgres.GenerateDB(db, *schema, *output, updateTemplate()); err != nil {
		log.Print(err)
		return exit.Failure
	}

	return exit.Success
}
