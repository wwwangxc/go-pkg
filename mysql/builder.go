package mysql

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
)

var (
	dbs  = map[string]*sql.DB{}
	dbRW sync.RWMutex
)

type mysqlBuilder struct {
	cliConfig clientConfig
}

func newMySQLBuilder(name string, opts ...Option) *mysqlBuilder {
	builder := &mysqlBuilder{
		cliConfig: getClientConfig(name),
	}

	for _, opt := range opts {
		opt(builder)
	}

	return builder
}

func (m *mysqlBuilder) build() (*sql.DB, error) {
	dbRW.RLock()
	db, ok := dbs[m.cliConfig.Name]
	dbRW.RUnlock()
	if ok {
		return db, nil
	}

	dbRW.Lock()
	defer dbRW.Unlock()

	db, ok = dbs[m.cliConfig.Name]
	if ok {
		return db, nil
	}

	db, err := sql.Open("mysql", m.cliConfig.DSN)
	if err != nil {
		return nil, fmt.Errorf("mysql open fail. error:%v", err)
	}

	if m.cliConfig.MaxIdle > 0 {
		db.SetMaxIdleConns(m.cliConfig.MaxIdle)
	}

	if m.cliConfig.MaxOpen > 0 {
		db.SetMaxOpenConns(m.cliConfig.MaxOpen)
	}

	if m.cliConfig.MaxIdleTime > 0 {
		db.SetConnMaxIdleTime(time.Duration(m.cliConfig.MaxIdleTime) * time.Millisecond)
	}

	dbs[m.cliConfig.Name] = db
	return db, nil
}
