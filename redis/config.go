package redis

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
		RedisCfg redisConfig    `yaml:"redis"`
		Client   []clientConfig `yaml:"client"`
	} `yaml:"database"`
}

func (a *appConfig) getClientConfigs() []clientConfig {
	if a == nil {
		return []clientConfig{}
	}

	clientConfigs := make([]clientConfig, 0, len(a.Database.Client))
	for _, v := range a.Database.Client {
		v.Wait = a.Database.RedisCfg.Wait

		if v.MaxIdle == 0 {
			v.MaxIdle = a.Database.RedisCfg.MaxIdle
		}

		if v.MaxActive == 0 {
			v.MaxActive = a.Database.RedisCfg.MaxActive
		}

		if v.IdleTimeout == 0 {
			v.IdleTimeout = a.Database.RedisCfg.IdleTimeout
		}

		if v.MaxConnLifetime == 0 {
			v.MaxConnLifetime = a.Database.RedisCfg.MaxConnLifetime
		}

		if v.Timeout == 0 {
			v.Timeout = a.Database.RedisCfg.Timeout
		}

		clientConfigs = append(clientConfigs, v)
	}

	return clientConfigs
}

type redisConfig struct {
	MaxIdle         int  `yaml:"max_idle"`
	MaxActive       int  `yaml:"max_active"`
	MaxConnLifetime int  `yaml:"max_conn_lifetime"`
	IdleTimeout     int  `yaml:"idle_timeout"`
	Timeout         int  `yaml:"timeout"`
	Wait            bool `yaml:"wait"`
}

type clientConfig struct {
	Name string `yaml:"name"`
	DSN  string `yaml:"dsn"`

	redisConfig `yaml:",inline"`
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
			redisConfig: redisConfig{
				MaxIdle:         2048,
				MaxActive:       0,
				IdleTimeout:     180000,
				MaxConnLifetime: 0,
				Timeout:         1000,
				Wait:            false,
			},
		}
		clientConfigMap[name] = c
	}

	return c
}
