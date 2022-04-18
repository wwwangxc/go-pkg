package example

import (
	"context"

	// go-pkg/redis will automatically read configuration
	// files (./go-pkg.yaml) when package loaded
	"github.com/wwwangxc/go-pkg/redis"
)

func ExampleNewClientProxy() {
	cli := redis.NewClientProxy("client_name",
		redis.WithDSN("dsn"),             // set dsn, default use database.client.dsn
		redis.WithMaxIdle(20),            // set max idel. default 2048
		redis.WithMaxActive(100),         // set max active. default 0
		redis.WithIdleTimeout(180000),    // set idle timeout. unit millisecond, default 180000
		redis.WithTimeout(1000),          // set command timeout. unit millisecond, default 1000
		redis.WithMaxConnLifetime(10000), // set max conn life time, default 0
		redis.WithWait(true),             // set wait
	)

	cli.Do(context.Background(), "GET", "foo")
	// do something...
}

func ExampleClientProxy_GetConn() {
	c := redis.NewClientProxy("client_name").GetConn()
	defer c.Close()

	c.Send("SET", "foo", "bar")
	c.Send("GET", "foo")
	c.Flush()
	c.Receive() // reply from SET
	c.Receive() // reply from GET
}
