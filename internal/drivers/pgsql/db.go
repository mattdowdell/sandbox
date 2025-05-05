// Package pgsql provides a PostgreSQL implementation of [sql.DB].
package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/stdlib"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Config contains configuration for connecting to a PostgreSQL database. It is intended to be
// populated by internal/drivers/config.Load.
type Config struct {
<<<<<<< Updated upstream
	Hostname     string        `koanf:"hostname"`
	Port         string        `koanf:"port" default:"5432"`
	Username     string        `koanf:"username"`
	Password     string        `koanf:"password"`
	UseIAMAuth   bool          `koanf:"useiamauth"`
	Name         string        `koanf:"name"`
	SSLMode      string        `koanf:"sslmode" default:"verify-full"`
	Region       string        `koanf:"region"`
	MaxIdleTime  time.Duration `koanf:"maxidletime" default:"5m"`
	MaxLifetime  time.Duration `koanf:"maxlifetime" default:"5m"`
	MaxIdleConns int           `koanf:"maxidleconns"`
	MaxOpenConns int           `koanf:"maxopenconns"`
=======
	// The hostname of the database server.
	Hostname string `koanf:"hostname"`

	// The port the database server is listening on. Defaults to 5432.
	Port string `koanf:"port" default:"5432"`

	// The username to authenticate with.
	Username string `koanf:"username"`

	// The password to authenticate with. This should be omitted when UseIAMAuth is set to true.
	Password string `koanf:"password" json:"-"`

	// Enables the use of IAM authentication. For use in AWS environments only.
	UseIAMAuth bool `koanf:"useiamauth"`

	// The name of the database to connect to.
	Name string `koanf:"name"`

	// The SSL mode to connect with. Defaults to verify-full.
	SSLMode string `koanf:"sslmode" default:"verify-full"`

	// The AWS region of the database server. this should be omitted when UseIAMAuth is false.
	Region string `koanf:"region"`

	// The maximum time a database connection can be idle before being closed. Defaults to 5
	// minutes. Set to 0 to disable.
	MaxIdleTime time.Duration `koanf:"maxidletime" default:"5m"`

	// The maximum lifetime for a database connection. Defaults to 5 minutes. Set to 0 to disable.
	MaxLifetime time.Duration `koanf:"maxlifetime" default:"5m"`

	// The maximum number of idle connections in the database connection pool. Defaults to no limit.
	MaxIdleConns int `koanf:"maxidleconns"`

	// The maximum number of open connections in the database connection pool. Defaults to no limit.
	//
	// It is strongly recommended to set this value when sharing a single database server across
	// multiple clients to avoid accidentally starving other clients of connections.
	MaxOpenConns int `koanf:"maxopenconns"`
>>>>>>> Stashed changes
}

func (c *Config) toOptions() []Option {
	options := []Option{
		WithMaxIdleTime(c.MaxIdleTime),
		WithMaxLifetime(c.MaxLifetime),
		WithMaxIdleConns(c.MaxIdleConns),
		WithMaxOpenConns(c.MaxOpenConns),
	}

	if c.Password != "" {
		options = append(options, WithPassword(c.Password))
	}

	if c.UseIAMAuth {
		endpoint := net.JoinHostPort(c.Hostname, c.Port)
		options = append(options, WithIAMAuth(endpoint, c.Region, c.Username))
	}

	return options
}

// NewFromConfig creates a new [sql.DB] using the given configuration.
//
//nolint:gocritic // called once, little gain from passing Config by pointer
func NewFromConfig(ctx context.Context, conf Config) (*sql.DB, error) {
	return New(
		ctx,
		conf.Hostname,
		conf.Port,
		conf.Username,
		conf.Name,
		conf.SSLMode,
		conf.toOptions()...,
	)
}

// New creates a new [sql.DB] for the given PostgreSQL database.
//
// The returned database connection pool has the following features:
//
//   - Built-in OpenTelemetry tracing and metrics.
//   - Automatic closure of connections that observe a read-only transaction error to enable recovery
//     when a primary server moves to standby.
func New(
	ctx context.Context,
	host string,
	port string,
	user string,
	name string,
	sslmode string,
	options ...Option,
) (*sql.DB, error) {
	opts := defaultOptions()
	for _, option := range options {
		if err := option.apply(ctx, opts); err != nil {
			return nil, fmt.Errorf("failed to create db: %w", err)
		}
	}

	if err := opts.validate(); err != nil {
		return nil, err
	}

	dsn := makeDSN(host, port, user, opts.password, name, sslmode)
	conf, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	conf.OnPgError = func(_ *pgconn.PgConn, err *pgconn.PgError) bool {
		// automatically close on any fatal errors
		if strings.EqualFold(err.Severity, "FATAL") {
			return false
		}

		// this error is produced if a write is attempted in a readonly transaction.
		// it can mean that the database primary moved to standby and now only accepts reads.
		// closing the connection allows the server ip address to be refreshed and enables faster
		// self-healing that the configured max connection lifetime would allow.
		if err.Code == pgerrcode.ReadOnlySQLTransaction {
			return false
		}

		return true
	}

	conn := stdlib.RegisterConnConfig(conf)

	db, err := otelsql.Open("pgx/v5", conn, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
	))
	if err != nil {
		return nil, err
	}

	opts.apply(db)

	if err := otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
	)); err != nil {
		return nil, err
	}

	return db, nil
}

func makeDSN(
	host string,
	port string,
	user string,
	password string,
	name string,
	sslmode string,
) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		name,
		sslmode,
	)
}
