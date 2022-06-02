package orm

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/wwwangxc/go-pkg/orm/driver"
)

var (
	dbs  = map[string]*gorm.DB{}
	dbRW sync.RWMutex
)

type gormBuilder struct {
	dbConfig   serviceConfig
	gormConfig gorm.Config
}

func newGORMBuilder(name string, opts ...GORMProxyOption) *gormBuilder {
	builder := &gormBuilder{
		dbConfig: getServiceConfig(name),
	}

	for _, opt := range opts {
		opt(builder)
	}

	return builder
}

func (g *gormBuilder) build() (*gorm.DB, error) {
	dbRW.RLock()
	db, ok := dbs[g.dbConfig.Name]
	dbRW.RUnlock()
	if ok {
		return db, nil
	}

	dbRW.Lock()
	defer dbRW.Unlock()

	db, ok = dbs[g.dbConfig.Name]
	if ok {
		return db, nil
	}

	d, exist := driver.Get(g.dbConfig.Driver)
	if !exist {
		return nil, fmt.Errorf("invalid driver:%s", g.dbConfig.Driver)
	}

	db, err := gorm.Open(d.Open(g.dbConfig.DSN), &g.gormConfig)
	if err != nil {
		return nil, fmt.Errorf("gorm open fail. error:%v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB fail. error:%v", err)
	}

	if g.dbConfig.MaxIdle > 0 {
		sqlDB.SetMaxIdleConns(g.dbConfig.MaxIdle)
	}

	if g.dbConfig.MaxOpen > 0 {
		sqlDB.SetMaxOpenConns(g.dbConfig.MaxOpen)
	}

	if g.dbConfig.MaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(g.dbConfig.MaxIdleTime) * time.Millisecond)
	}

	dbs[g.dbConfig.Name] = db
	return db, nil
}
