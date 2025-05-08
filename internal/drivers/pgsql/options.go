package pgsql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
)

// Option customises the behaviour when creating a [sql.DB].
type Option interface {
	apply(context.Context, *dbOptions) error
}

type dbOptions struct {
	password     string
	maxIdleTime  time.Duration
	maxLifetime  time.Duration
	maxIdleConns int
	maxOpenConns int
}

//nolint:mnd // defaults are documented in the relevant options
func defaultOptions() *dbOptions {
	return &dbOptions{
		maxIdleTime: time.Minute * 5,
		maxLifetime: time.Minute * 5,
	}
}

func (o *dbOptions) validate() error {
	if o.password == "" {
		return errors.New("missing database password")
	}

	return nil
}

func (o *dbOptions) apply(db *sql.DB) {
	if o.maxIdleTime > 0 {
		db.SetConnMaxIdleTime(o.maxIdleTime)
	}

	if o.maxLifetime > 0 {
		db.SetConnMaxLifetime(o.maxLifetime)
	}

	if o.maxIdleConns > 0 {
		db.SetMaxIdleConns(o.maxIdleConns)
	}

	if o.maxOpenConns > 0 {
		db.SetMaxOpenConns(o.maxOpenConns)
	}
}

type passwordOpt string

// WithPassword sets the password to authenticate with.
func WithPassword(password string) Option {
	return passwordOpt(password)
}

func (o passwordOpt) apply(_ context.Context, opts *dbOptions) error {
	opts.password = string(o)
	return nil
}

type iamAuthOpt struct {
	endpoint string
	region   string
	user     string
}

// WithAIMAuth enables the use of AWS IAM for database authentication.
//
// Based on [AWS docs] (untested).
//
// [AWS docs]: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.Connecting.Go.html
func WithIAMAuth(endpoint, region, user string) Option {
	return iamAuthOpt{
		endpoint: endpoint,
		region:   region,
		user:     user,
	}
}

func (o iamAuthOpt) apply(ctx context.Context, opts *dbOptions) error {
	conf, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	token, err := auth.BuildAuthToken(ctx, o.endpoint, o.region, o.user, conf.Credentials)
	if err != nil {
		return fmt.Errorf("failed to create auth token: %w", err)
	}

	opts.password = token
	return nil
}

type maxIdleTimeOpt time.Duration

// WithMaxIdleTime sets the maximum idle time for a database connection. Defaults to 5 minutes. Set
// to 0 to disable.
//
// The maximum idle time should be set to avoid the effects of idle connections being closed. For
// example, Istio will close idle connections after 60 minutes, which is only observable when trying
// to use the now closed connection. By setting a timeout below this value, Istio's timeout will
// never be observed and the availability impact never felt.
func WithMaxIdleTime(d time.Duration) Option {
	return maxIdleTimeOpt(d)
}

func (o maxIdleTimeOpt) apply(_ context.Context, opts *dbOptions) error {
	opts.maxIdleTime = time.Duration(o)
	return nil
}

type maxLifetimeOpt time.Duration

// WithMaxLifetime sets the maximum lifetime for a database connection. Defaults to 5 minutes. Set
// to 0 to disable.
//
// The 5 minute default value allows some connection issues to be automatically recovered from in a
// relatively short period of time. For example, if a deadlock is created, it can last no longer
// than 5 minutes. The chosen value is by no means perfect, but is intended to be a good start to
// refine from.
func WithMaxLifetime(d time.Duration) Option {
	return maxLifetimeOpt(d)
}

func (o maxLifetimeOpt) apply(_ context.Context, opts *dbOptions) error {
	opts.maxLifetime = time.Duration(o)
	return nil
}

type maxIdleConnsOpt int

// WithMaxIdleConns limits the number of idle connections in the database connection pool.
func WithMaxIdleConns(count int) Option {
	return maxIdleConnsOpt(count)
}

func (o maxIdleConnsOpt) apply(_ context.Context, opts *dbOptions) error {
	opts.maxIdleConns = int(o)
	return nil
}

type maxOpenConnsOpt int

// WithMaxOpenConns limits the number of open connections in the database connection pool.
func WithMaxOpenConns(count int) Option {
	return maxOpenConnsOpt(count)
}

func (o maxOpenConnsOpt) apply(_ context.Context, opts *dbOptions) error {
	opts.maxOpenConns = int(o)
	return nil
}
