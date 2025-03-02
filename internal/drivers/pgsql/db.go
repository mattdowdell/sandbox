package pgsql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/stdlib"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
)

// ...
type Config struct {
	Hostname     string        `koanf:"hostname"`
	Port         uint16        `koanf:"port" default:"5432"`
	Username     string        `koanf:"username"`
	Password     string        `koanf:"password"`
	Name         string        `koanf:"name"`
	SSLMode      string        `koanf:"sslmode" default:"verify-full"`
	MaxIdleTime  time.Duration `koanf:"maxidletime" default:"5m"`
	MaxLifetime  time.Duration `koanf:"maxlifetime" default:"5m"`
	MaxIdleConns int           `koanf:"maxidleconns"`
	MaxOpenConns int           `koanf:"maxopenconns"`
}

// ...
func (c *Config) toOptions() []Option {
	var options []Option

	if c.Password != "" {
		options = append(options, WithPassword(c.Password))
	}

	return options
}

// ...
//
//nolint:gocritic // called once, little gain from passing by pointer
func NewFromConfig(conf Config) (*sql.DB, error) {
	return New(
		conf.Hostname,
		conf.Port,
		conf.Username,
		conf.Name,
		conf.SSLMode,
		conf.toOptions()...,
	)
}

// ...
func New(
	host string,
	port uint16,
	user string,
	name string,
	sslmode string,
	options ...Option,
) (*sql.DB, error) {
	opts := defaultOptions()
	for _, option := range options {
		option.apply(opts)
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

		// this error is produced if a write is attempted in a readonly transaction
		// it can mean that the database primary moved to standby and now only accepts reads
		// closing to refreshes the ip address and so enables self-healing
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

// ...
func makeDSN(
	host string,
	port uint16,
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
