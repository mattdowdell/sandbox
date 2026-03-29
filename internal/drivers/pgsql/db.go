// Package pgsql provides a PostgreSQL implementation of [sql.DB].
package pgsql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/stdlib"
	"go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/mattdowdell/sandbox/internal/drivers/config/secret"
)

type Cleanup func() error

// Config contains configuration for connecting to a PostgreSQL database. It is intended to be
// populated by internal/drivers/config.Load.
type Config struct {
	// The hostname of the database server.
	Hostname string `koanf:"hostname"`

	// The port the database server is listening on. Defaults to 5432.
	Port int `koanf:"port" default:"5432"`

	// The username to authenticate with.
	Username string `koanf:"username"`

	// The password to authenticate with. This should be omitted when UseIAMAuth is set to true.
	Password secret.String `koanf:"password"`

	// Enables the use of IAM authentication. For use in AWS environments only.
	UseIAMAuth bool `koanf:"useiamauth"`

	// The AWS region of the database server. For use in AWS environments only.
	Region string `koanf:"region"`

	// The name of the database to connect to.
	Name string `koanf:"name"`

	// The SSL mode to connect with. Defaults to verify-full.
	SSLMode string `koanf:"sslmode" default:"verify-full"`

	// The maximum time a database connection can be idle before being closed. Defaults to 5
	// minutes. Set to 0 to disable.
	//
	// This value should be less than or equal to than MaxLifetime, as it otherwise has no effect.
	// It should also be less than any external idle timeouts. This notably includes the idle
	// timeout implemented by Istio, which defaults to 60 minutes.
	MaxIdleTime time.Duration `koanf:"maxidletime" default:"5m"`

	// The maximum lifetime for a database connection. Defaults to 5 minutes. Set to 0 to disable.
	//
	// Setting lower values improves availability when server updates occur, such as DNS updates and
	// failovers. However, opening a new connection is somewhat expensive. Therefore, closing it
	// prematurely can increase latency during large, sudden bursts of requests, especially if those
	// bursts occur less frequently than the specified max lifetime.
	MaxLifetime time.Duration `koanf:"maxlifetime" default:"5m"`

	// The maximum number of idle connections in the database connection pool. Defaults to no limit.
	MaxIdleConns int `koanf:"maxidleconns"`

	// The maximum number of open connections in the database connection pool. Defaults to no limit.
	//
	// It is strongly recommended to set this value when sharing a single database server across
	// multiple clients to avoid accidentally starving other clients of connections.
	MaxOpenConns int `koanf:"maxopenconns"`
}

func (c *Config) toOptions() []Option {
	options := []Option{
		WithMaxIdleTime(c.MaxIdleTime),
		WithMaxLifetime(c.MaxLifetime),
		WithMaxIdleConns(c.MaxIdleConns),
		WithMaxOpenConns(c.MaxOpenConns),
	}

	if c.Password != "" {
		options = append(options, WithPassword(c.Password.String()))
	}

	if c.UseIAMAuth {
		endpoint := net.JoinHostPort(c.Hostname, strconv.Itoa(c.Port))
		options = append(options, WithIAMAuth(endpoint, c.Region, c.Username))
	}

	return options
}

// NewFromConfig creates a new [sql.DB] using the given configuration.
//
//nolint:gocritic // called once, little gain from passing Config by pointer
func NewFromConfig(ctx context.Context, conf Config) (*sql.DB, Cleanup, error) {
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
//
// The returned cleanup function can be used to close the database and unregister metrics when the
// databse is no longer needed. If the caller calls New once on application start and keep the
// connection pool until exit, cleanup does not need to be called.
func New(
	ctx context.Context,
	host string,
	port int,
	user string,
	name string,
	sslmode string,
	options ...Option,
) (*sql.DB, Cleanup, error) {
	opts := defaultOptions()
	for _, option := range options {
		if err := option.apply(ctx, opts); err != nil {
			return nil, nil, fmt.Errorf("failed to create db: %w", err)
		}
	}

	if err := opts.validate(); err != nil {
		return nil, nil, err
	}

	dsn := makeDSN(host, port, user, opts.password, name, sslmode)
	conf, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}

	opts.apply(db)

	reg, err := otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
	))
	if err != nil {
		return nil, db.Close, err
	}

	return db, func() error {
		return errors.Join(db.Close(), reg.Unregister())
	}, nil
}

func makeDSN(
	host string,
	port int,
	user string,
	password string,
	name string,
	sslmode string,
) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		name,
		sslmode,
	)
}
