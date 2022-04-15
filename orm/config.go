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
		MaxIdle     int            `yaml:"max_idle"`
		MaxOpen     int            `yaml:"max_open"`
		MaxIdleTime int            `yaml:"max_idle_time"`
		Client      []clientConfig `yaml:"client"`
	} `yaml:"database"`
}

func (a *appConfig) getClientConfigs() []clientConfig {
	if a == nil {
		return []clientConfig{}
	}

	clientConfigs := make([]clientConfig, 0, len(a.Database.Client))
	for _, v := range a.Database.Client {
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
	Driver      string `yaml:"driver"`
	MaxIdle     int    `yaml:"max_idle"`
	MaxOpen     int    `yaml:"max_open"`
	MaxIdleTime int    `yaml:"max_idle_time"`
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
