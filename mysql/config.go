package mysql

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/wwwangxc/go-pkg/config"
)

var (
	clientConfigMap = map[string]clientConfig{}
	clientConfigMu  sync.Mutex

	dbs  = map[string]*sql.DB{}
	dbRW sync.RWMutex
)

func init() {
	c, err := loadAppConfig()
	if err != nil {
		logErrorf("config load fail. error:%v", err)
		return
	}

	for _, v := range c.getClientConfigs() {
		registerClientConfig(v)
	}
}

type appConfig struct {
	Database struct {
		MaxIdle     int            `yaml:"max_idle"`
		MaxOpen     int            `yaml:"max_open"`
		MaxIdleTime int            `yaml:"max_idle_time"`
		MySQL       []clientConfig `yaml:"mysql"`
	} `yaml:"database"`
}

func (a *appConfig) getClientConfigs() []clientConfig {
	if a == nil {
		return []clientConfig{}
	}

	clientConfigs := make([]clientConfig, 0, len(a.Database.MySQL))
	for _, v := range a.Database.MySQL {
		if v.MaxIdle == 0 {
			v.MaxIdle = a.Database.MaxIdle
		}

		if v.MaxOpen == 0 {
			v.MaxOpen = a.Database.MaxOpen
		}

		if v.MaxIdleTime == 0 {
			v.MaxIdleTime = a.Database.MaxIdleTime
		}

		clientConfigs = append(clientConfigs, v)
	}

	return clientConfigs
}

type clientConfig struct {
	Name        string `yaml:"name"`
	DSN         string `yaml:"dsn"`
	MaxIdle     int    `yaml:"max_idle"`
	MaxOpen     int    `yaml:"max_open"`
	MaxIdleTime int    `yaml:"max_idle_time"`
}

func (c *clientConfig) buildDB() (*sql.DB, error) {
	dbRW.RLock()
	db, ok := dbs[c.Name]
	dbRW.RUnlock()
	if ok {
		return db, nil
	}

	dbRW.Lock()
	defer dbRW.Unlock()

	db, ok = dbs[c.Name]
	if ok {
		return db, nil
	}

	db, err := sql.Open("mysql", c.DSN)
	if err != nil {
		return nil, fmt.Errorf("mysql open fail. error:%v", err)
	}

	if c.MaxIdle > 0 {
		db.SetMaxIdleConns(c.MaxIdle)
	}

	if c.MaxOpen > 0 {
		db.SetMaxOpenConns(c.MaxOpen)
	}

	if c.MaxIdleTime > 0 {
		db.SetConnMaxIdleTime(time.Duration(c.MaxIdleTime) * time.Millisecond)
	}

	dbs[c.Name] = db
	return db, nil
}

func loadAppConfig() (*appConfig, error) {
	configure, err := config.Load("./go-pkg.yaml")
	if err != nil {
		return &appConfig{}, nil
	}

	c := &appConfig{}
	if err = configure.Unmarshal(c); err != nil {
		return nil, fmt.Errorf("config unmarshal fail. error:%v", err)
	}

	return c, nil
}

func registerClientConfig(c clientConfig) {
	clientConfigMu.Lock()
	defer clientConfigMu.Unlock()
	clientConfigMap[c.Name] = c
}

func getClientConfig(name string) clientConfig {
	clientConfigMu.Lock()
	defer clientConfigMu.Unlock()

	c, exist := clientConfigMap[name]
	if !exist {
		c = clientConfig{
			Name: name,
		}
		clientConfigMap[name] = c
	}

	return c
}
