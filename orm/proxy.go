package orm

import (
	"gorm.io/gorm"
)

// NewGORMProxy new GORM DB proxy
func NewGORMProxy(name string, opts ...GORMProxyOption) (*gorm.DB, error) {
	return newGORMBuilder(name, opts...).build()
}
