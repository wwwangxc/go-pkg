package redis

import (
	"context"

	redigo "github.com/gomodule/redigo/redis"
)

// ClientProxy Redis client proxy
//go:generate mockgen -source=client.go -destination=mockredis/client_mock.go -package=mockredis
type ClientProxy interface {

	// Do sends a command to server and returns the received reply.
	// min(ctx,DialReadTimeout()) will be used as the deadline.
	// The connection will be closed if DialReadTimeout() timeout or ctx timeout or ctx canceled when this function is running.
	// DialReadTimeout() timeout return err can be checked by strings.Contains(e.Error(), "io/timeout").
	// ctx timeout return err context.DeadlineExceeded.
	// ctx canceled return err context.Canceled.
	Do(ctx context.Context, cmd string, args ...interface{}) (interface{}, error)

	// GetConn gets a connection. The application must close the returned connection.
	// This method always returns a valid connection so that applications can defer
	// error handling to the first use of the connection. If there is an error
	// getting an underlying connection, then the connection Err, Do, Send, Flush
	// and Receive methods return that error.
	GetConn() redigo.Conn
}

type clientProxyImpl struct {
	pool *redigo.Pool
}

// NewClientProxy new redis client proxy
func NewClientProxy(name string, opts ...Option) ClientProxy {
	return &clientProxyImpl{
		pool: newRedisBuilder(name, opts...).build(),
	}
}

// Do sends a command to server and returns the received reply.
//
// min(ctx,DialReadTimeout()) will be used as the deadline.
// The connection will be closed if DialReadTimeout() timeout or ctx timeout or ctx canceled when this function is running.
// DialReadTimeout() timeout return err can be checked by strings.Contains(e.Error(), "io/timeout").
// ctx timeout return err context.DeadlineExceeded.
// ctx canceled return err context.Canceled.
func (c *clientProxyImpl) Do(ctx context.Context, cmd string, args ...interface{}) (interface{}, error) {
	conn := c.pool.Get()
	defer func() {
		if err := conn.Close(); err != nil {
			logErrorf("connect close fail. error:%v", err)
		}
	}()
	defer conn.Close()

	return redigo.DoContext(conn, ctx, cmd, args...)
}

// GetConn gets a connection. The application must close the returned connection.
// This method always returns a valid connection so that applications can defer
// error handling to the first use of the connection. If there is an error
// getting an underlying connection, then the connection Err, Do, Send, Flush
// and Receive methods return that error.
func (c *clientProxyImpl) GetConn() redigo.Conn {
	return c.pool.Get()
}
