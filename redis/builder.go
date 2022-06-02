package redis

import (
	"context"
	"net"
	"sync"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

var (
	pools   = map[string]*redigo.Pool{}
	poolsRW sync.RWMutex
)

type redisBuilder struct {
	cliConfig serviceConfig
}

func newRedisBuilder(name string, opts ...ClientOption) *redisBuilder {
	builder := &redisBuilder{
		cliConfig: getserviceConfig(name),
	}

	for _, opt := range opts {
		opt(builder)
	}

	return builder
}

func (r *redisBuilder) build() *redigo.Pool {
	poolsRW.RLock()
	pool, ok := pools[r.cliConfig.Name]
	poolsRW.RUnlock()
	if ok {
		return pool
	}

	poolsRW.Lock()
	defer poolsRW.Unlock()

	pool, ok = pools[r.cliConfig.Name]
	if ok {
		return pool
	}

	timeout := time.Duration(r.cliConfig.Timeout) * time.Millisecond
	pool = &redigo.Pool{
		MaxIdle:         r.cliConfig.MaxIdle,
		MaxActive:       r.cliConfig.MaxActive,
		IdleTimeout:     time.Duration(r.cliConfig.IdleTimeout) * time.Millisecond,
		MaxConnLifetime: time.Duration(r.cliConfig.MaxConnLifetime) * time.Millisecond,
		Dial: func() (redigo.Conn, error) {
			dialOpts := []redigo.DialOption{
				redigo.DialWriteTimeout(timeout),
				redigo.DialReadTimeout(timeout),
				redigo.DialConnectTimeout(timeout),
				redigo.DialContextFunc(func(ctx context.Context, network, addr string) (net.Conn, error) {
					dialer := &net.Dialer{
						Timeout: timeout,
					}
					return dialer.DialContext(ctx, network, addr)
				}),
			}

			c, err := redigo.DialURL(r.cliConfig.DSN, dialOpts...)
			if err != nil {
				return nil, err
			}

			return c, nil
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		Wait: r.cliConfig.Wait,
	}

	pools[r.cliConfig.Name] = pool
	return pool
}
