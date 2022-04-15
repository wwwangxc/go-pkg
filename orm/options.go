package orm

import "gorm.io/gorm"

// GORMProxyOption GORM DB proxy option
type GORMProxyOption func(*gormBuilder)

// WithDSN set dsn
func WithDSN(dsn string) GORMProxyOption {
	return func(g *gormBuilder) {
		g.dbConfig.DSN = dsn
	}
}

// WithMaxIdle set the maximum number of connections in the idle
// connection pool.
//
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns,
// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit.
func WithMaxIdle(maxIdel int) GORMProxyOption {
	return func(g *gormBuilder) {
		g.dbConfig.MaxIdle = maxIdel
	}
}

// WithMaxOpen set the maximum number of open connections to the database.
//
// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
// MaxIdleConns, then MaxIdleConns will be reduced to match the new
// MaxOpenConns limit.
func WithMaxOpen(maxOpen int) GORMProxyOption {
	return func(g *gormBuilder) {
		g.dbConfig.MaxOpen = maxOpen
	}
}

// WithMaxIdleTime set the maximum amount of time a connection may be reused.
//
// Expired connections may be closed lazily before reuse.
// Uint: milliseconds
func WithMaxIdleTime(maxIdelTime int) GORMProxyOption {
	return func(g *gormBuilder) {
		g.dbConfig.MaxIdleTime = maxIdelTime
	}
}

// WithGORMConfig set GORM config.
func WithGORMConfig(gormConfig *gorm.Config) GORMProxyOption {
	return func(g *gormBuilder) {
		g.gormConfig = *gormConfig
	}
}

// WithDriver set GORM driver
//
// Support: mysql, postgresql, sqlite, sqlserver, clickhouse
// - postgresql uses the extended protocol
// - disables implicit prepared statement use postgresql.simple
func WithDriver(driver string) GORMProxyOption {
	return func(g *gormBuilder) {
		g.dbConfig.Driver = driver
	}
}
