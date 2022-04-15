package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	mysqlCfg, exist := clientConfigMap["db_mysql"]
	assert.True(t, exist, "db_mysql should exist")
	assert.Equal(t, "db_mysql", mysqlCfg.Name)
	assert.Equal(t, "root:root@tcp(127.0.0.1:3306)/db1?charset=utf8&parseTime=True", mysqlCfg.DSN)
	assert.Equal(t, "mysql", mysqlCfg.Driver)
	assert.Equal(t, 111, mysqlCfg.MaxIdle)
	assert.Equal(t, 222, mysqlCfg.MaxOpen)
	assert.Equal(t, 333, mysqlCfg.MaxIdleTime)

	postgresqlCfg, exist := clientConfigMap["db_postgresql"]
	assert.True(t, exist, "db_postgresql should exist")
	assert.Equal(t, "db_postgresql", postgresqlCfg.Name)
	assert.Equal(t, "user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai", postgresqlCfg.DSN)
	assert.Equal(t, "postgresql", postgresqlCfg.Driver)
	assert.Equal(t, 11, postgresqlCfg.MaxIdle)
	assert.Equal(t, 22, postgresqlCfg.MaxOpen)
	assert.Equal(t, 33, postgresqlCfg.MaxIdleTime)
}
