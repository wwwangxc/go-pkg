package driver

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	register(&sqliteDriver{})
}

type sqliteDriver struct{}

// Open return GORM sqlite sialector
func (s *sqliteDriver) Open(dsn string) gorm.Dialector {
	return sqlite.Open(dsn)
}

func (s *sqliteDriver) name() string {
	return "sqlite"
}
