# go-pkg Redis Package

[![Go Report Card](https://goreportcard.com/badge/github.com/wwwangxc/go-pkg/redis)](https://goreportcard.com/report/github.com/wwwangxc/go-pkg/redis)
[![GoDoc](https://pkg.go.dev/badge/github.com/wwwangxc/go-pkg/redis?status.svg)](https://pkg.go.dev/github.com/wwwangxc/go-pkg/redis)

go-pkg/redis is an componentized redis package, it provides an easy way to create and manage redis pool.

Based off [gomodule/redigo](https://github.com/gomodule/redigo).

## Install

```sh
go get github.com/wwwangxc/go-pkg/redis
```

## Quick Start

**main.go**

```go
package main

import (
        "context"
        "fmt"

        // go-pkg/redis will automatically read configuration
        // files (./go-pkg.yaml) when package loaded
        "github.com/wwwangxc/go-pkg/redis"
)

func main() {
        cli := redis.NewClientProxy("client_name",
                redis.WithDSN("dsn"),             // set dsn, default use database.client.dsn
                redis.WithMaxIdle(20),            // set max idel. default 2048
                redis.WithMaxActive(100),         // set max active. default 0
                redis.WithIdleTimeout(180000),    // set idle timeout. unit millisecond, default 180000
                redis.WithTimeout(1000),          // set command timeout. unit millisecond, default 1000
                redis.WithMaxConnLifetime(10000), // set max conn life time, default 0
                redis.WithWait(true),             // set wait
        )

        // Exec GET command
        cli.Do(context.Background(), "GET", "foo")

        // Pipeline
        // get a redis connection
        c := cli.GetConn()
        defer c.Close()

        c.Send("SET", "foo", "bar")
        c.Send("GET", "foo")
        c.Flush()
        c.Receive()           // reply from SET
        v, err := c.Receive() // reply from GET
        fmt.Sprintf("reply: %s", v)
        fmt.Sprintf("error: %v", err)
}
```

**go-pkg.yaml**

```yaml
database:
  redis:
    max_idle: 20
    max_active: 100
    max_conn_lifetime: 1000
    idle_timeout: 180000
    timeout: 1000
    wait: true
  client:
    - name: redis_1
      dsn: redis://username:password@127.0.0.1:6379/1?timeout=1000ms
    - name: redis_2
      dsn: redis://username:password@127.0.0.1:6379/2?timeout=1000ms
      max_idle: 22
      max_active: 111
      max_conn_lifetime: 2000
      idle_timeout: 200000
      timeout: 2000

```

## How To Mock

```go
package tests

import (
    "testing"
    
    "github.com/agiledragon/gomonkey"
    "github.com/golang/mock/gomock"

    "github.com/wwwangxc/go-pkg/redis"
    "github.com/wwwangxc/go-pkg/redis/mockredis"
)

func TestMockClientProxy(t *testing.T){
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockCli := mockredis.NewMockClientProxy(ctrl)
    mockCli.EXPECT().Do(gomock.Any(), gomock.Any(), gomock.Any()).
        Return("reply", nil).AnyTimes()
    
    patches := gomonkey.ApplyFunc(redis.NewClientProxy,
        func(string, ...redis.Option) (redis.ClientProxy, error) {
            return mockCli, nil
        })
    defer patches.Reset()
}

func TestMockRedisConn(t *testing.T){
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockConn := mockredis.NewMockConn(ctrl)
    mockConn.EXPECT().Close().Return(nil).AnyTimes()
    mockConn.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
    mockConn.EXPECT().Flush().Return(nil).AnyTimes()
    mockConn.EXPECT().Receive().Return(nil, nil).AnyTimes()

    mockCli := mockredis.NewMockClientProxy(ctrl)
    mockCli.EXPECT().GetConn().Return(mockConn).AnyTimes()
    
    patches := gomonkey.ApplyFunc(redis.NewClientProxy,
        func(string, ...redis.Option) (redis.ClientProxy, error) {
            return mockCli, nil
        })
    defer patches.Reset()
}
```
