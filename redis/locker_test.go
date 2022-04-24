package redis

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock/v3"
)

func Test_lockerImpl_TryLock(t *testing.T) {
	type fields struct {
		cli ClientProxy
	}
	type args struct {
		ctx  context.Context
		key  string
		opts []LockOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
		intRet  int
		intErr  error
	}{
		{
			name:    "int fail",
			wantErr: true,
			intErr:  fmt.Errorf(""),
		},
		{
			name:    "lock not acquired",
			wantErr: true,
			intRet:  0,
		},
		{
			name:    "normal process",
			wantErr: false,
			intRet:  2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cli *clientProxyImpl
			patches := gomonkey.ApplyMethod(reflect.TypeOf(cli), "GetConn",
				func(*clientProxyImpl) redigo.Conn {
					return redigomock.NewConn()
				})
			defer patches.Reset()

			var script *redigo.Script
			patches.ApplyMethod(reflect.TypeOf(script), "DoContext",
				func(*redigo.Script, context.Context, redigo.Conn, ...interface{}) (interface{}, error) {
					return nil, nil
				})

			patches.ApplyFunc(Int,
				func(interface{}, error) (int, error) {
					return tt.intRet, tt.intErr
				})

			l := &lockerImpl{
				cli: NewClientProxy("client_name"),
			}
			_, err := l.TryLock(tt.args.ctx, tt.args.key, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("lockerImpl.TryLock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_lockerImpl_Lock(t *testing.T) {
	type fields struct {
		cli ClientProxy
	}
	type args struct {
		ctx  context.Context
		key  string
		opts []LockOption
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       string
		wantErr    bool
		tryLockErr error
	}{
		{
			name:       "try lock fail",
			wantErr:    true,
			tryLockErr: fmt.Errorf(""),
		},
		{
			name:    "normal process",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &lockerImpl{
				cli: tt.fields.cli,
			}

			patches := gomonkey.ApplyMethod(reflect.TypeOf(l), "TryLock",
				func(*lockerImpl, context.Context, string, ...LockOption) (string, error) {
					return "", tt.tryLockErr
				})
			defer patches.Reset()

			got, err := l.Lock(tt.args.ctx, tt.args.key, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("lockerImpl.Lock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("lockerImpl.Lock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lockerImpl_Unlock(t *testing.T) {
	type fields struct {
		cli ClientProxy
	}
	type args struct {
		ctx  context.Context
		key  string
		uuid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		intRet  int
		intErr  error
	}{
		{
			name:    "convert to int fail",
			wantErr: true,
			intErr:  fmt.Errorf(""),
		},
		{
			name:    "lock not exist",
			wantErr: true,
			intRet:  0,
		},
		{
			name:    "not owner of key",
			wantErr: true,
			intRet:  1,
		},
		{
			name:    "key delete fail",
			wantErr: true,
			intRet:  2,
		},
		{
			name:    "error unknown",
			wantErr: true,
			intRet:  3,
		},
		{
			name:    "normal process",
			wantErr: false,
			intRet:  666,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cli *clientProxyImpl
			patches := gomonkey.ApplyMethod(reflect.TypeOf(cli), "GetConn",
				func(*clientProxyImpl) redigo.Conn {
					return redigomock.NewConn()
				})
			defer patches.Reset()

			var script *redigo.Script
			patches.ApplyMethod(reflect.TypeOf(script), "DoContext",
				func(*redigo.Script, context.Context, redigo.Conn, ...interface{}) (interface{}, error) {
					return nil, nil
				})

			patches.ApplyFunc(Int,
				func(interface{}, error) (int, error) {
					return tt.intRet, tt.intErr
				})
			l := &lockerImpl{
				cli: NewClientProxy("client_name"),
			}

			err := l.Unlock(tt.args.ctx, tt.args.key, tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("lockerImpl.Unlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
