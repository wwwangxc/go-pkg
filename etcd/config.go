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
	serviceConfigMap = map[string]serviceConfig{}
	serviceConfigRW  sync.RWMutex
)

func init() {
	c, err := loadAppConfig()
	if err != nil {
		log.Errorf("config load fail. error:%v", err)
		return
	}

	for _, v := range c.Client.Service {
		registerServiceConfig(v)
	}
}

type appConfig struct {
	Client struct {
		Service []serviceConfig `yaml:"service"`
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

type serviceConfig struct {
	Name        string `yaml:"name"`
	Target      string `yaml:"target"`
	Timeout     int    `yaml:"timeout"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	TLSKeyPath  string `yaml:"tls_key"`
	TLSCertPath string `yaml:"tls_cert"`
	CACertPath  string `yaml:"ca_cert"`
}

func newServiceConfig() *serviceConfig {
	return &serviceConfig{
		Name:        "",
		Target:      "",
		Timeout:     3000,
		Username:    "",
		Password:    "",
		TLSKeyPath:  "",
		TLSCertPath: "",
		CACertPath:  "",
	}
}

func (s *serviceConfig) etcdConfig() (*clientv3.Config, error) {
	tlsConfig, err := s.tlsConfig()
	if err != nil {
		return nil, err
	}

	cfg := newServiceConfig()
	cfg.Name = s.Name
	cfg.Target = s.Target
	cfg.Username = s.Username
	cfg.Password = s.Password
	cfg.TLSKeyPath = s.TLSKeyPath
	cfg.TLSCertPath = s.TLSCertPath
	cfg.CACertPath = s.CACertPath

	if s.Timeout > 0 {
		cfg.Timeout = s.Timeout
	}

	return &clientv3.Config{
		Endpoints:   strings.Split(s.Target, ","),
		DialTimeout: time.Millisecond * time.Duration(s.Timeout),
		Username:    s.Username,
		Password:    s.Password,
		TLS:         tlsConfig,
	}, nil
}

func (s *serviceConfig) tlsConfig() (*tls.Config, error) {
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

func registerServiceConfig(c serviceConfig) {
	serviceConfigRW.Lock()
	defer serviceConfigRW.Unlock()
	serviceConfigMap[c.Name] = c
}

func getServiceConfig(name string) serviceConfig {
	serviceConfigRW.RLock()
	defer serviceConfigRW.RUnlock()

	c, exist := serviceConfigMap[name]
	if !exist {
		c = serviceConfig{
			Name: name,
		}
		serviceConfigMap[name] = c
	}

	return c
}
