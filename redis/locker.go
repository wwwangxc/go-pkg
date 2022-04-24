package redis

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

// Locker distributed lock provider
//go:generate mockgen -source=locker.go -destination=mockredis/locker_mock.go -package=mockredis
type Locker interface {

	// TryLock try get lock, if lock acquired will return lock uuid.
	//
	// Not block the current goroutine.
	// Return ErrLockNotAcquired when lock not acquired.
	// Will reentrant lock when UUID option not empty.
	// If Heartbeat option not empty and not a reentrant lock, will automatically
	// renewal until unlocked.
	TryLock(ctx context.Context, key string, opts ...LockOption) (uuid string, err error)

	// Lock try get lock first, if the lock is not acquired, the unlock event
	// will be subscribed until the context canceled or the lock is acquired.
	//
	// Will block the current goroutine.
	// Will reentrant lock when UUID option not empty.
	// If Heartbeat option not empty and not a reentrant lock, will automatically
	// renewal until unlocked.
	Lock(ctx context.Context, key string, opts ...LockOption) (uuid string, err error)

	// Unlock
	//
	// Return ErrLockNotExist if the key does not exist.
	// Return ErrNotOwnerOfKey if the uuid invalid.
	// Support reentrant unlock.
	Unlock(ctx context.Context, key, uuid string) error
}

type lockerImpl struct {
	cli ClientProxy
}

// NewLocker new locker proxy
func NewLocker(name string, opts ...ClientOption) Locker {
	return NewClientProxy(name, opts...).GetLocker()
}

func newLocker(cli ClientProxy) Locker {
	return &lockerImpl{
		cli: cli,
	}
}

// TryLock try get lock, if lock acquired will return lock uuid.
//
// Not block the current goroutine.
// Return ErrLockNotAcquired when lock not acquired.
// Will reentrant lock when UUID option not empty.
// If Heartbeat option not empty and not a reentrant lock, will automatically
// renewal until unlocked.
func (l *lockerImpl) TryLock(ctx context.Context, key string, opts ...LockOption) (string, error) {
	k := fmt.Sprintf("%s.lock", strings.TrimSuffix(key, ".lock"))
	options := newLockOptions(opts...)
	script := redigo.NewScript(1, luaScriptLock)
	conn := l.cli.GetConn()
	defer func() {
		if err := conn.Close(); err != nil {
			logErrorf("connect close fail. error:%v", err)
		}
	}()

	lockCount, err := Int(script.DoContext(ctx, conn, k, options.UUID, options.Expire.Milliseconds()))
	if err != nil {
		return "", err
	}

	if lockCount == 0 {
		return "", ErrLockNotAcquired
	}

	if lockCount == 1 && options.Heartbeat > 0 {
		go l.sendLockHeartbeat(key, options.Expire, options.Heartbeat)
	}

	return options.UUID, nil
}

// Lock try get lock first, if the lock is not acquired, the unlock event
// will be subscribed until the context canceled or the lock is acquired.
//
// Will block the current goroutine.
// Will reentrant lock when UUID option not empty.
// If Heartbeat option not empty and not a reentrant lock, will automatically
// renewal until unlocked.
func (l *lockerImpl) Lock(ctx context.Context, key string, opts ...LockOption) (string, error) {
	uuid, err := l.TryLock(ctx, key, opts...)
	if err != nil {
		if IsErrLockNotAcquired(err) {
			return l.waitUntilLock(ctx, key, opts...)
		}

		return "", err
	}

	return uuid, nil
}

// Unlock
//
// Return ErrLockNotExist if the key does not exist.
// Return ErrNotOwnerOfKey if the uuid invalid.
// Support reentrant unlock.
func (l *lockerImpl) Unlock(ctx context.Context, key, uuid string) error {
	k := fmt.Sprintf("%s.lock", strings.TrimSuffix(key, ".lock"))
	script := redigo.NewScript(1, luaScriptUnlock)
	conn := l.cli.GetConn()
	defer func() {
		if err := conn.Close(); err != nil {
			logErrorf("connect close fail. error:%v", err)
		}
	}()

	ret, err := Int(script.DoContext(ctx, conn, k, uuid))
	if err != nil {
		return err
	}

	switch ret {
	case 0:
		return ErrLockNotExist
	case 1:
		return ErrNotOwnerOfKey
	case 2:
		return errors.New("locker key delete fail")
	case 666:
		return nil
	}

	return errors.New("error unknown")
}

func (l *lockerImpl) sendLockHeartbeat(key string, expire, heartbeatInterval time.Duration) {
	conn := l.cli.GetConn()
	defer func() {
		if err := conn.Close(); err != nil {
			logErrorf("connect close fail. error:%v", err)
		}
	}()

	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for {
		<-ticker.C

		exist, err := Bool(conn.Do("EXISTS", key))
		if err != nil {
			logErrorf(err.Error())
			return
		}

		if !exist {
			return
		}

		_, err = Bool(conn.Do("PEXPIRE", key, expire.Milliseconds()))
		if err != nil {
			logErrorf(err.Error())
			return
		}
	}
}

func (l *lockerImpl) waitUntilLock(ctx context.Context, key string, opts ...LockOption) (string, error) {
	k := fmt.Sprintf("%s.lock", strings.TrimSuffix(key, ".lock"))
	psc := redigo.PubSubConn{Conn: l.cli.GetConn()}
	defer func() {
		if err := psc.Close(); err != nil {
			logErrorf("pub/sub connect close fail. error:%v", err)
		}
	}()

	if err := psc.Subscribe(k); err != nil {
		return "", fmt.Errorf("subscribe: %s fail. error:%v", k, err)
	}

	ch := make(chan interface{})
	go func() {
		ch <- psc.Receive()
	}()

	select {
	case <-ctx.Done():
		return "", errors.New("lock fail: context canceled")
	case <-ch:
		return l.Lock(ctx, key, opts...)
	}
}
