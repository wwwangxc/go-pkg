package orm

import (
	"fmt"
	"sync"

	"github.com/wwwangxc/go-pkg/config"
)

var (
	clientConfigMap = map[string]clientConfig{}
	clientConfigMu  sync.Mutex
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
		MySQL      dbConfig       `json:"mysql"`
		PostgreSQL dbConfig       `json:"postgresql"`
		SQLite     dbConfig       `json:"sqlite"`
		SQLServer  dbConfig       `json:"sqlserver"`
		Clickhouse dbConfig       `json:"clickhouse"`
		CliConfig  []clientConfig `yaml:"client" json:"cli_config"`
	} `yaml:"database" json:"database"`
}

func (a *appConfig) getClientConfigs() []clientConfig {
	if a == nil {
		return []clientConfig{}
	}

	clientConfigs := make([]clientConfig, 0, len(a.Database.CliConfig))
	for _, v := range a.Database.CliConfig {
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

		clientConfigs = append(clientConfigs, v)
	}

	return clientConfigs
}

func (a *appConfig) getDBConfig(driver string) dbConfig {
	switch driver {
	case "mysql":
		return a.Database.MySQL
	case "postgresql":
		return a.Database.PostgreSQL
	case "sqlite":
		return a.Database.SQLite
	case "sqlserver":
		return a.Database.SQLServer
	case "clickhouse":
		return a.Database.Clickhouse
	default:
		return dbConfig{}
	}
}

type dbConfig struct {
	MaxIdle     int `yaml:"max_idle"`
	MaxOpen     int `yaml:"max_open"`
	MaxIdleTime int `yaml:"max_idle_time"`
}

type clientConfig struct {
	Name   string `yaml:"name"`
	DSN    string `yaml:"dsn"`
	Driver string `yaml:"driver"`

	dbConfig `yaml:",inline"`
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
			Name:   name,
			Driver: "mysql",
		}
		clientConfigMap[name] = c
	}

	return c
}
