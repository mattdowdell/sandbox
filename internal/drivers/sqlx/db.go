package sqlx

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/XSAM/otelsql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
)

// ...
type Config struct {
	Hostname     string        `koanf:"host"`
	Port         uint16        `koanf:"port" default:"5432"`
	User         string        `koanf:"username"`
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

type DB struct {
	db *sql.DB
}

// ...
func NewFromConfig(conf Config, options ...Option) (*DB, error) {
	options = append(options, conf.toOptions()...)

	return New(
		conf.Hostname,
		conf.Port,
		conf.User,
		conf.Name,
		conf.SSLMode,
		options...,
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
) (*DB, error) {
	opts := defaultOptions()
	for _, option := range options {
		option.apply(opts)
	}

	if err := opts.validate(); err != nil {
		return nil, err
	}

	dsn := makeDSN(host, port, user, opts.password, name, sslmode)

	db, err := otelsql.Open("pgx/v5", dsn, otelsql.WithAttributes(
		semconv.DBSystemMySQL,
	))
	if err != nil {
		return nil, err
	}

	opts.apply(db)

	return nil, nil
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
