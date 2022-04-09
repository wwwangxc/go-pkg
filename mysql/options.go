package mysql

import "database/sql"

// TxOption transaction option
type TxOption func(*sql.TxOptions)

// WithIsolation set transaction isolation level
//
// Isolation is the transaction isolation level.
// If zero, the driver or database's default level is used.
func WithIsolation(isolation sql.IsolationLevel) TxOption {
	return func(options *sql.TxOptions) {
		options.Isolation = isolation
	}
}

// WithReadOnly set transaction readonly
func WithReadOnly(readOnly bool) TxOption {
	return func(options *sql.TxOptions) {
		options.ReadOnly = readOnly
	}
}

// ClientProxyOption mysql client proxy option
type ClientProxyOption func(*clientConfig)

// WithDSN set dsn
func WithDSN(dsn string) ClientProxyOption {
	return func(c *clientConfig) {
		c.DSN = dsn
	}
}

// WithMaxIdle sets the maximum number of connections in the idle
// connection pool.
//
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns,
// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit.
func WithMaxIdle(maxIdel int) ClientProxyOption {
	return func(c *clientConfig) {
		c.MaxIdle = maxIdel
	}
}

// WithMaxOpen sets the maximum number of open connections to the database.
//
// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
// MaxIdleConns, then MaxIdleConns will be reduced to match the new
// MaxOpenConns limit.
func WithMaxOpen(maxOpen int) ClientProxyOption {
	return func(c *clientConfig) {
		c.MaxOpen = maxOpen
	}
}

// WithMaxIdleTime sets the maximum amount of time a connection may be reused.
//
// Expired connections may be closed lazily before reuse.
// Uint: milliseconds
func WithMaxIdleTime(maxIdelTime int) ClientProxyOption {
	return func(c *clientConfig) {
		c.MaxIdleTime = maxIdelTime
	}
}
