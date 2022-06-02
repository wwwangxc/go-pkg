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
	cfg serviceConfig
}

func newMySQLBuilder(name string, opts ...Option) *mysqlBuilder {
	builder := &mysqlBuilder{
		cfg: getServiceConfig(name),
	}

	for _, opt := range opts {
		opt(builder)
	}

	return builder
}

func (m *mysqlBuilder) build() (*sql.DB, error) {
	dbRW.RLock()
	db, ok := dbs[m.cfg.Name]
	dbRW.RUnlock()
	if ok {
		return db, nil
	}

	dbRW.Lock()
	defer dbRW.Unlock()

	db, ok = dbs[m.cfg.Name]
	if ok {
		return db, nil
	}

	db, err := sql.Open("mysql", m.cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("mysql open fail. error:%v", err)
	}

	if m.cfg.MaxIdle > 0 {
		db.SetMaxIdleConns(m.cfg.MaxIdle)
	}

	if m.cfg.MaxOpen > 0 {
		db.SetMaxOpenConns(m.cfg.MaxOpen)
	}

	if m.cfg.MaxIdleTime > 0 {
		db.SetConnMaxIdleTime(time.Duration(m.cfg.MaxIdleTime) * time.Millisecond)
	}

	dbs[m.cfg.Name] = db
	return db, nil
}
