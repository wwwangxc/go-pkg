package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func Test_clientProxyImpl_Transaction(t *testing.T) {
	errBeginTx := errors.New("begin transaction fail")
	errTxFunc := errors.New("exec tx func fail")
	errCommit := errors.New("commit fail")
	errRollback := errors.New("rollback fail")
	type fields struct {
		c  *clientConfig
		db *sql.DB
	}
	type args struct {
		ctx  context.Context
		f    TxFunc
		opts []TxOption
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		want        error
		beginTxErr  error
		commitErr   error
		rollbackErr error
	}{
		{
			name:    "begin transaction fail",
			wantErr: true,
			fields: fields{
				db: &sql.DB{},
			},
			args: args{
				ctx: context.Background(),
			},
			beginTxErr: errBeginTx,
			want:       errBeginTx,
		},
		{
			name:    "exec tx func fail",
			wantErr: true,
			fields: fields{
				db: &sql.DB{},
			},
			args: args{
				ctx: context.Background(),
				f: func(t *sql.Tx) error {
					return errTxFunc
				},
			},
			want: errTxFunc,
		},
		{
			name:    "rollback fail affter tx func",
			wantErr: true,
			fields: fields{
				db: &sql.DB{},
			},
			args: args{
				ctx: context.Background(),
				f: func(t *sql.Tx) error {
					return errTxFunc
				},
			},
			rollbackErr: errRollback,
			want:        errRollback,
		},
		{
			name:    "commit fail",
			wantErr: true,
			fields: fields{
				db: &sql.DB{},
			},
			args: args{
				ctx: context.Background(),
				f: func(t *sql.Tx) error {
					return nil
				},
			},
			commitErr: errCommit,
			want:      errCommit,
		},
		{
			name:    "rollback fail after commit",
			wantErr: true,
			fields: fields{
				db: &sql.DB{},
			},
			args: args{
				ctx: context.Background(),
				f: func(t *sql.Tx) error {
					return nil
				},
			},
			commitErr:   errCommit,
			rollbackErr: errRollback,
			want:        errRollback,
		},
		{
			name:    "normal",
			wantErr: false,
			fields: fields{
				db: &sql.DB{},
			},
			args: args{
				ctx: context.Background(),
				f: func(t *sql.Tx) error {
					return nil
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := gomonkey.ApplyMethod(reflect.TypeOf(tt.fields.db), "BeginTx",
				func(*sql.DB, context.Context, *sql.TxOptions) (*sql.Tx, error) {
					return &sql.Tx{}, tt.beginTxErr
				})
			defer patches.Reset()

			var tx *sql.Tx
			patches.ApplyMethod(reflect.TypeOf(tx), "Commit",
				func(*sql.Tx) error {
					return tt.commitErr
				})

			patches.ApplyMethod(reflect.TypeOf(tx), "Rollback",
				func(*sql.Tx) error {
					return tt.rollbackErr
				})

			c := &clientProxyImpl{
				c:  tt.fields.c,
				db: tt.fields.db,
			}

			err := c.Transaction(tt.args.ctx, tt.args.f, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("clientProxyImpl.Transaction() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.want, err)
		})
	}
}

func Test_clientProxyImpl_Query(t *testing.T) {
	type fields struct {
		c  *clientConfig
		db *sql.DB
	}
	type args struct {
		ctx   context.Context
		f     ScanFunc
		query string
		args  []interface{}
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantErr         bool
		hasNext         bool
		queryContextErr error
	}{
		{
			name:    "query context fail",
			wantErr: true,
			fields: fields{
				db: &sql.DB{},
			},
			args: args{
				ctx: context.Background(),
			},
			queryContextErr: fmt.Errorf(""),
		},
		{
			name:    "exec scan func fail",
			wantErr: true,
			fields: fields{
				db: &sql.DB{},
			},
			args: args{
				ctx: context.Background(),
				f: func(r *sql.Rows) error {
					return fmt.Errorf("")
				},
			},
			hasNext: true,
		},
		{
			name:    "normal",
			wantErr: false,
			fields: fields{
				db: &sql.DB{},
			},
			args: args{
				ctx: context.Background(),
				f: func(r *sql.Rows) error {
					return nil
				},
			},
			hasNext: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := gomonkey.ApplyMethod(reflect.TypeOf(tt.fields.db), "QueryContext",
				func(*sql.DB, context.Context, string, ...interface{}) (*sql.Rows, error) {
					return &sql.Rows{}, tt.queryContextErr
				})
			defer patches.Reset()

			i := 1
			var rows *sql.Rows
			patches.ApplyMethod(reflect.TypeOf(rows), "Next",
				func(*sql.Rows) bool {
					i--
					return tt.hasNext && i >= 0
				})

			patches.ApplyMethod(reflect.TypeOf(rows), "Close",
				func(*sql.Rows) error {
					return nil
				})

			c := &clientProxyImpl{
				c:  tt.fields.c,
				db: tt.fields.db,
			}
			if err := c.Query(tt.args.ctx, tt.args.f, tt.args.query, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("clientProxyImpl.Query() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_clientProxyImpl_QueryRow(t *testing.T) {
	type fields struct {
		c  *clientConfig
		db *sql.DB
	}
	type args struct {
		ctx   context.Context
		dest  []interface{}
		query string
		args  []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		scanErr error
	}{
		{
			name:    "scan fail",
			wantErr: true,
			fields: fields{
				db: &sql.DB{},
			},
			args: args{
				ctx: context.Background(),
			},
			scanErr: fmt.Errorf(""),
		},
		{
			name:    "normal",
			wantErr: false,
			fields: fields{
				db: &sql.DB{},
			},
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := gomonkey.ApplyMethod(reflect.TypeOf(tt.fields.db), "QueryRowContext",
				func(*sql.DB, context.Context, string, ...interface{}) *sql.Row {
					return &sql.Row{}
				})
			defer patches.Reset()

			var row *sql.Row
			patches.ApplyMethod(reflect.TypeOf(row), "Scan",
				func(*sql.Row, ...interface{}) error {
					return tt.scanErr
				})

			c := &clientProxyImpl{
				c:  tt.fields.c,
				db: tt.fields.db,
			}
			if err := c.QueryRow(tt.args.ctx, tt.args.dest, tt.args.query, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("clientProxyImpl.QueryRow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_clientProxyImpl_Select(t *testing.T) {
	type fields struct {
		c  *clientConfig
		db *sql.DB
	}
	type args struct {
		ctx   context.Context
		dest  interface{}
		query string
		args  []interface{}
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantErr         bool
		queryContextErr error
	}{
		{
			name:    "query context fail",
			wantErr: true,
			fields: fields{
				db: &sql.DB{},
			},
			args: args{
				ctx: context.Background(),
			},
			queryContextErr: fmt.Errorf(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := gomonkey.ApplyMethod(reflect.TypeOf(tt.fields.db), "QueryContext",
				func(*sql.DB, context.Context, string, ...interface{}) (*sql.Rows, error) {
					return &sql.Rows{}, tt.queryContextErr
				})
			defer patches.Reset()

			var rows *sql.Rows
			patches.ApplyMethod(reflect.TypeOf(rows), "Close",
				func(*sql.Rows) error {
					return nil
				})

			c := &clientProxyImpl{
				c:  tt.fields.c,
				db: tt.fields.db,
			}
			if err := c.Select(tt.args.ctx, tt.args.dest, tt.args.query, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("clientProxyImpl.Select() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
