package redis

// Option redis proxy option
type Option func(*redisBuilder)

// WithDSN set dsn
func WithDSN(dsn string) Option {
	return func(b *redisBuilder) {
		b.cliConfig.DSN = dsn
	}
}

// WithMaxIdle set max idle
//
// Maximum number of connections in the idle connection pool.
// Default 2048
func WithMaxIdle(maxIdle int) Option {
	return func(b *redisBuilder) {
		b.cliConfig.MaxIdle = maxIdle
	}
}

// WithMaxActive set max active
//
// Maximum number of connections allocated by the pool at a given time.
// When zero, there is no limit on the number of connections in the pool.
// Default 0
func WithMaxActive(maxActive int) Option {
	return func(b *redisBuilder) {
		b.cliConfig.MaxActive = maxActive
	}
}

// WithIdleTimeout set idle timeout
//
// Close connections after remaining idle for this duration. If the value
// is zero, then idle connections are not closed. Applications should set
// the timeout to a value less than the server's timeout.
// Unit millisecond, default 180000
func WithIdleTimeout(idleTimeout int) Option {
	return func(b *redisBuilder) {
		b.cliConfig.IdleTimeout = idleTimeout
	}
}

// WithMaxConnLifetime set max conn lifetime
//
// Close connections older than this duration. If the value is zero, then
// the pool does not close connections based on age.
// Unit millisecond, default 0
func WithMaxConnLifetime(maxConnLifetime int) Option {
	return func(b *redisBuilder) {
		b.cliConfig.MaxConnLifetime = maxConnLifetime
	}
}

// WithTimeout set timeout
//
// Write, read and connect timeout
// Unit millisecond, default 1000
func WithTimeout(timeout int) Option {
	return func(b *redisBuilder) {
		b.cliConfig.Timeout = timeout
	}
}

// WithWait set wait
//
// If Wait is true and the pool is at the MaxActive limit, then Get() waits
// for a connection to be returned to the pool before returning.
func WithWait(wait bool) Option {
	return func(b *redisBuilder) {
		b.cliConfig.Wait = wait
	}
}
