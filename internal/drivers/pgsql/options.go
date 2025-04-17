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

// ...
type Option interface {
	apply(context.Context, *dbOptions) error
}

// ...
type dbOptions struct {
	password     string
	maxIdleTime  time.Duration
	maxLifetime  time.Duration
	maxIdleConns int
	maxOpenConns int
}

// ...
func defaultOptions() *dbOptions {
	return &dbOptions{}
}

// ...
func (o *dbOptions) validate() error {
	if o.password == "" {
		return errors.New("missing database password")
	}

	return nil
}

// ...
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

// ...
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

// ...
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

// ...
func WithMaxIdleTime(d time.Duration) Option {
	return maxIdleTimeOpt(d)
}

func (o maxIdleTimeOpt) apply(_ context.Context, opts *dbOptions) error {
	opts.maxIdleTime = time.Duration(o)
	return nil
}

type maxLifetimeOpt time.Duration

// ...
func WithMaxLifetime(d time.Duration) Option {
	return maxLifetimeOpt(d)
}

func (o maxLifetimeOpt) apply(_ context.Context, opts *dbOptions) error {
	opts.maxLifetime = time.Duration(o)
	return nil
}

type maxIdleConnsOpt int

// ...
func WithMaxIdleConns(count int) Option {
	return maxIdleConnsOpt(count)
}

func (o maxIdleConnsOpt) apply(_ context.Context, opts *dbOptions) error {
	opts.maxIdleConns = int(o)
	return nil
}

type maxOpenConnsOpt int

// ...
func WithMaxOpenConns(count int) Option {
	return maxOpenConnsOpt(count)
}

func (o maxOpenConnsOpt) apply(_ context.Context, opts *dbOptions) error {
	opts.maxOpenConns = int(o)
	return nil
}
