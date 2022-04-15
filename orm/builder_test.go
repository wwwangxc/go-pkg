package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_gormBuilder_build(t *testing.T) {
	type fields struct {
		dbConfig   clientConfig
		gormConfig gorm.Config
	}
	tests := []struct {
		name     string
		fields   fields
		wantErr  bool
		openErr  error
		getDBErr error
	}{
		{
			name:    "invalid driver",
			wantErr: true,
			fields: fields{
				dbConfig: clientConfig{
					Driver: "test driver",
				},
			},
		},
		{
			name:    "open fail",
			wantErr: true,
			fields: fields{
				dbConfig: clientConfig{
					Driver: "mysql",
				},
			},
			openErr: fmt.Errorf(""),
		},
		{
			name:    "get db fail",
			wantErr: true,
			fields: fields{
				dbConfig: clientConfig{
					Driver: "mysql",
				},
			},
			getDBErr: fmt.Errorf(""),
		},
		{
			name:    "normal",
			wantErr: false,
			fields: fields{
				dbConfig: clientConfig{
					Name:   "test",
					Driver: "mysql",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := gomonkey.ApplyFunc(gorm.Open,
				func(gorm.Dialector, ...gorm.Option) (*gorm.DB, error) {
					return &gorm.DB{}, tt.openErr
				})
			defer patches.Reset()

			var db *gorm.DB
			patches.ApplyMethod(reflect.TypeOf(db), "DB",
				func(*gorm.DB) (*sql.DB, error) {
					return &sql.DB{}, tt.getDBErr
				})

			g := &gormBuilder{
				dbConfig:   tt.fields.dbConfig,
				gormConfig: tt.fields.gormConfig,
			}
			_, err := g.build()
			if (err != nil) != tt.wantErr {
				t.Errorf("gormBuilder.build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, 1, len(dbs))
				_, exist := dbs[tt.fields.dbConfig.Name]
				assert.True(t, exist)
			}
		})
	}
}
