package redis

import "errors"

var (
	// ErrLockNotAcquired lock not acquired
	ErrLockNotAcquired = errors.New("lock not acquired")

	// ErrLockNotAcquired lock dose not exist
	ErrLockNotExist = errors.New("lock does not exist")

	// ErrNotOwnerOfKey not the owner of the key
	ErrNotOwnerOfKey = errors.New("not the owner of the key")
)

// IsErrLockNotAcquired is lock not acquired error
func IsErrLockNotAcquired(err error) bool {
	return errors.Is(err, ErrLockNotAcquired)
}
