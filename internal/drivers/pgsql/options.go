package pgsql

import (
	"database/sql"
	"errors"
	"time"
)

// ...
type Option interface {
	apply(*dbOptions)
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

func (o passwordOpt) apply(opts *dbOptions) {
	opts.password = string(o)
}

// TODO: support RDS IAM postgres auth
// see https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.Connecting.Go.html

type maxIdleTimeOpt time.Duration

// ...
func WithMaxIdleTime(d time.Duration) Option {
	return maxIdleTimeOpt(d)
}

func (o maxIdleTimeOpt) apply(opts *dbOptions) {
	opts.maxIdleTime = time.Duration(o)
}

type maxLifetimeOpt time.Duration

// ...
func WithMaxLifetime(d time.Duration) Option {
	return maxLifetimeOpt(d)
}

func (o maxLifetimeOpt) apply(opts *dbOptions) {
	opts.maxLifetime = time.Duration(o)
}

type maxIdleConnsOpt int

// ...
func WithMaxIdleConns(count int) Option {
	return maxIdleConnsOpt(count)
}

func (o maxIdleConnsOpt) apply(opts *dbOptions) {
	opts.maxIdleConns = int(o)
}

type maxOpenConnsOpt int

// ...
func WithMaxOpenConns(count int) Option {
	return maxOpenConnsOpt(count)
}

func (o maxOpenConnsOpt) apply(opts *dbOptions) {
	opts.maxOpenConns = int(o)
}
