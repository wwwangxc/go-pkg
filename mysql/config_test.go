package mysql

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	cli1, exist := clientConfigMap["client1"]
	assert.True(t, exist, "client1 should exist")
	assert.Equal(t, "client1", cli1.Name)
	assert.Equal(t, "root:root@tcp(127.0.0.1:3306)/db1?charset=utf8&parseTime=True", cli1.DSN)
	assert.Equal(t, 11, cli1.MaxIdle)
	assert.Equal(t, 22, cli1.MaxOpen)
	assert.Equal(t, 33, cli1.MaxIdleTime)

	cli2, exist := clientConfigMap["client2"]
	assert.True(t, exist, "client2 should exist")
	assert.Equal(t, "client2", cli2.Name)
	assert.Equal(t, "root:root@tcp(127.0.0.1:3306)/db2?charset=utf8&parseTime=True", cli2.DSN)
	assert.Equal(t, 111, cli2.MaxIdle)
	assert.Equal(t, 222, cli2.MaxOpen)
	assert.Equal(t, 333, cli2.MaxIdleTime)
}

func Test_clientConfig_buildDB(t *testing.T) {
	tests := []struct {
		name       string
		wantErr    bool
		sqlOpenErr error
	}{
		{
			name:       "sql open fail",
			wantErr:    true,
			sqlOpenErr: fmt.Errorf(""),
		},
		{
			name:    "normal",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbs = map[string]*sql.DB{}

			patches := gomonkey.ApplyFunc(sql.Open,
				func(string, string) (*sql.DB, error) {
					return &sql.DB{}, tt.sqlOpenErr
				})
			defer patches.Reset()

			c := getClientConfig("client1")
			_, err := c.buildDB()
			if (err != nil) != tt.wantErr {
				t.Errorf("clientConfig.buildDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.Equal(t, len(dbs), 1)
				_, exist := dbs["client1"]
				assert.True(t, exist)
			}
		})
	}
}
