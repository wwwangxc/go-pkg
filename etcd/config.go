package etcd

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"crypto/tls"

	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/wwwangxc/go-pkg/config"
	"github.com/wwwangxc/go-pkg/etcd/log"
)

var (
	clientConfigMap = map[string]clientConfig{}
	clientConfigRW  sync.RWMutex
)

func init() {
	c, err := loadAppConfig()
	if err != nil {
		log.Errorf("config load fail. error:%v", err)
		return
	}

	c.registerClientConfig()
}

type appConfig struct {
	Client struct {
		Timeout int            `yaml:"timeout"`
		Service []clientConfig `yaml:"service"`
	} `yaml:"client"`
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

func (a *appConfig) registerClientConfig() {
	defaultTimeout := defaultClientConfig("").Timeout

	for _, v := range a.Client.Service {
		if v.Timeout < 1 {
			v.Timeout = a.Client.Timeout
		}

		if v.Timeout < 1 {
			v.Timeout = defaultTimeout
		}

		registerClientConfig(v)
	}
}

type clientConfig struct {
	Name        string `yaml:"name"`
	Target      string `yaml:"target"`
	Timeout     int    `yaml:"timeout"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	TLSKeyPath  string `yaml:"tls_key"`
	TLSCertPath string `yaml:"tls_cert"`
	CACertPath  string `yaml:"ca_cert"`
}

func defaultClientConfig(name string) clientConfig {
	return clientConfig{
		Name:        name,
		Target:      "",
		Timeout:     3000,
		Username:    "",
		Password:    "",
		TLSKeyPath:  "",
		TLSCertPath: "",
		CACertPath:  "",
	}
}

func (s *clientConfig) etcdConfig() (*clientv3.Config, error) {
	tlsConfig, err := s.tlsConfig()
	if err != nil {
		return nil, err
	}

	return &clientv3.Config{
		Endpoints:   strings.Split(s.Target, ","),
		DialTimeout: time.Millisecond * time.Duration(s.Timeout),
		Username:    s.Username,
		Password:    s.Password,
		TLS:         tlsConfig,
	}, nil
}

func (s *clientConfig) tlsConfig() (*tls.Config, error) {
	if s.TLSKeyPath == "" || s.TLSCertPath == "" || s.CACertPath == "" {
		return nil, nil
	}

	tlsInfo := &transport.TLSInfo{
		TrustedCAFile: s.CACertPath,
		CertFile:      s.TLSCertPath,
		KeyFile:       s.TLSKeyPath,
	}

	return tlsInfo.ClientConfig()
}

func registerClientConfig(c clientConfig) {
	clientConfigRW.Lock()
	defer clientConfigRW.Unlock()
	clientConfigMap[c.Name] = c
}

func getClientConfig(name string) clientConfig {
	clientConfigRW.RLock()
	defer clientConfigRW.RUnlock()

	c, exist := clientConfigMap[name]
	if !exist {
		clientConfigMap[name] = defaultClientConfig(name)
	}

	return c
}
