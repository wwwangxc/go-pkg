package redis

import (
	"context"
	"errors"

	redigo "github.com/gomodule/redigo/redis"
)

// Fetcher object fetcher
//go:generate mockgen -source=fetcher.go -destination=mockredis/fetcher_mock.go -package=mockredis
type Fetcher interface {

	// Fetch data and storing the result into the struct pointed at by dest.
	//
	// Use json decode
	Fetch(ctx context.Context, key string, dest interface{}, opts ...FetchOption) error
}

type fetcherImpl struct {
	name string
	opts []ClientOption
}

// NewFetcher new object fetcher
func NewFetcher(name string, opts ...ClientOption) Fetcher {
	return &fetcherImpl{
		name: name,
		opts: opts,
	}
}

// Fetch data and storing the result into the struct pointed at by dest.
//
// Use json decode
func (f *fetcherImpl) Fetch(ctx context.Context, key string, dest interface{}, opts ...FetchOption) error {
	options := newFetchOptions(opts...)
	conn := f.getConn()
	defer func() {
		if err := conn.Close(); err != nil {
			logErrorf("conn close fail. error:%v", err)
		}
	}()

	data, err := Bytes(redigo.DoContext(conn, ctx, "GET", key))
	if err != nil && !errors.Is(redigo.ErrNil, err) {
		return err
	}

	if errors.Is(redigo.ErrNil, err) && options.Callback != nil {
		val, err := options.Callback()
		if err != nil {
			return err
		}

		data, err = options.Marshal(val)
		if err != nil {
			return err
		}

		_, err = redigo.DoContext(conn, ctx, "PSETEX", key, options.Expire.Milliseconds(), data)
		if err != nil {
			return err
		}
	}

	return options.Unmarshal(data, dest)
}

func (f *fetcherImpl) getConn() redigo.Conn {
	return getRedisPool(f.name, f.opts...).Get()
}
