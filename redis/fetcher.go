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
	cli ClientProxy
}

// NewFetcher new object fetcher
func NewFetcher(name string, opts ...ClientOption) Fetcher {
	return NewClientProxy(name, opts...).GetFetcher()
}

func newFetcher(cli ClientProxy) Fetcher {
	return &fetcherImpl{
		cli: cli,
	}
}

// Fetch data and storing the result into the struct pointed at by dest.
//
// Use json decode
func (f *fetcherImpl) Fetch(ctx context.Context, key string, dest interface{}, opts ...FetchOption) error {
	options := newFetchOptions(opts...)

	data, err := Bytes(f.cli.Do(ctx, "GET", key))
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

		_, err = f.cli.Do(ctx, "PSETEX", key, options.Expire.Milliseconds(), data)
		if err != nil {
			return err
		}
	}

	return options.Unmarshal(data, dest)
}
