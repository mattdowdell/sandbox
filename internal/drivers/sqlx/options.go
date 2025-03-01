package sqlx

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

// ...
func WithPassword(password string) Option {
	return nil
}

// TODO: support RDS IAM postgres auth
// see https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.Connecting.Go.html

// ...
func WithMaxIdleTime(d time.Duration) Option {
	return nil
}

// ...
func WithMaxLifetime(d time.Duration) Option {
	return nil
}

// ...
func WithMaxIdleConns(count int) Option {
	return nil
}

// ...
func WithMaxOpenConns(count int) Option {
	return nil
}
