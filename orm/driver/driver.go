package driver

import (
	"sync"

	"gorm.io/gorm"
)

var (
	driverMap   = map[string]Driver{}
	driverMapRW sync.RWMutex
)

// Driver database driver
type Driver interface {

	// Open return GORM database dialector
	Open(dsn string) gorm.Dialector

	name() string
}

func register(d Driver) {
	driverMapRW.Lock()
	defer driverMapRW.Unlock()
	driverMap[d.name()] = d
}

// Get database driver
func Get(name string) (Driver, bool) {
	driverMapRW.RLock()
	defer driverMapRW.RUnlock()

	d, exist := driverMap[name]
	return d, exist
}
