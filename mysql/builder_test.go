package mysql

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
)

func Test_mysqlBuilder_build(t *testing.T) {
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

			_, err := newMySQLBuilder("client1").build()
			if (err != nil) != tt.wantErr {
				t.Errorf("mysqlBuilder.build() error = %v, wantErr %v", err, tt.wantErr)
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
