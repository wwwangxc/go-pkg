package driver

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func init() {
	register(&sqlserverDriver{})
}

type sqlserverDriver struct{}

// Open return GORM sqlserver dialector
func (s *sqlserverDriver) Open(dsn string) gorm.Dialector {
	return sqlserver.Open(dsn)
}

func (s *sqlserverDriver) name() string {
	return "sqlserver"
}
