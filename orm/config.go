package orm

import (
	"fmt"
	"sync"

	"github.com/wwwangxc/go-pkg/config"
	"github.com/wwwangxc/go-pkg/orm/driver"
	"gorm.io/gorm"
)

var (
	serviceConfigMap = map[string]serviceConfig{}
	serviceConfigMu  sync.Mutex
)

func init() {
	c, err := loadAppConfig()
	if err != nil {
		logErrorf("config load fail. error:%v", err)
		return
	}

	for _, v := range c.getServiceConfigs() {
		registerServiceConfig(v)
	}
}

type appConfig struct {
	Client struct {
		MySQL      dbConfig        `yaml:"mysql"`
		PostgreSQL dbConfig        `yaml:"postgresql"`
		SQLite     dbConfig        `yaml:"sqlite"`
		SQLServer  dbConfig        `yaml:"sqlserver"`
		Clickhouse dbConfig        `yaml:"clickhouse"`
		Service    []serviceConfig `yaml:"service"`
	} `yaml:"client"`
}

func (a *appConfig) getServiceConfigs() []serviceConfig {
	if a == nil {
		return []serviceConfig{}
	}

	serviceConfigs := make([]serviceConfig, 0, len(a.Client.Service))
	for _, v := range a.Client.Service {
		dbCfg := a.getDBConfig(v.Driver)
		if v.MaxIdle == 0 {
			v.MaxIdle = dbCfg.MaxIdle
		}

		if v.MaxOpen == 0 {
			v.MaxOpen = dbCfg.MaxOpen
		}

		if v.MaxIdleTime == 0 {
			v.MaxIdleTime = dbCfg.MaxIdleTime
		}

		serviceConfigs = append(serviceConfigs, v)
	}

	return serviceConfigs
}

func (a *appConfig) getDBConfig(driverName string) dbConfig {
	switch driverName {
	case driver.NameMySQL:
		return a.Client.MySQL
	case driver.NamePostgreSQL, driver.NamePostgreSQLSimple:
		return a.Client.PostgreSQL
	case driver.NameSQLite:
		return a.Client.SQLite
	case driver.NameSQLServer:
		return a.Client.SQLServer
	case driver.NameClickhouse:
		return a.Client.Clickhouse
	default:
		return dbConfig{}
	}
}

type dbConfig struct {
	MaxIdle     int `yaml:"max_idle"`
	MaxOpen     int `yaml:"max_open"`
	MaxIdleTime int `yaml:"max_idle_time"`
}

type serviceConfig struct {
	Name     string `yaml:"name"`
	DSN      string `yaml:"dsn"`
	Driver   string `yaml:"driver"`
	dbConfig `yaml:",inline"`

	gormConfig *gorm.Config `yaml:"-"`
}

func loadAppConfig() (*appConfig, error) {
	configure, err := config.Load("./app.yaml")
	if err != nil {
		return &appConfig{}, nil
	}

	c := &appConfig{}
	if err = configure.Unmarshal(c); err != nil {
		return nil, fmt.Errorf("config unmarshal fail. error:%v", err)
	}

	return c, nil
}

func registerServiceConfig(c serviceConfig) {
	serviceConfigMu.Lock()
	defer serviceConfigMu.Unlock()
	serviceConfigMap[c.Name] = c
}

func getServiceConfig(name string) serviceConfig {
	serviceConfigMu.Lock()
	defer serviceConfigMu.Unlock()

	c, exist := serviceConfigMap[name]
	if !exist {
		c = serviceConfig{
			Name:   name,
			Driver: "mysql",
		}
		serviceConfigMap[name] = c
	}

	return c
}
