package concurrency

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	type args struct {
		ctx         context.Context
		executors   []Executor
		concurrenty uint8
	}
	tests := []struct {
		name string
		args args
		want *Result
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				executors: []Executor{
					&executorNormal{},
					&executorNormal{},
					&executorNormal{},
					&executorNormal{},
					&executorNormal{},
				},
				concurrenty: 3,
			},
			want: &Result{
				failed: false,
				resultSet: []*singleResult{
					{result: "success", err: nil},
					{result: "success", err: nil},
					{result: "success", err: nil},
					{result: "success", err: nil},
					{result: "success", err: nil},
				},
			},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				executors: []Executor{
					&executorError{},
					&executorError{},
					&executorError{},
					&executorError{},
					&executorError{},
				},
				concurrenty: 3,
			},
			want: &Result{
				failed: true,
				resultSet: []*singleResult{
					{result: nil, err: fmt.Errorf("error message")},
					{result: nil, err: fmt.Errorf("error message")},
					{result: nil, err: fmt.Errorf("error message")},
					{result: nil, err: fmt.Errorf("error message")},
					{result: nil, err: fmt.Errorf("error message")},
				},
			},
		},
		{
			name: "panic",
			args: args{
				ctx: context.Background(),
				executors: []Executor{
					&executorPanic{},
					&executorPanic{},
					&executorPanic{},
					&executorPanic{},
					&executorPanic{},
				},
				concurrenty: 3,
			},
			want: &Result{
				failed: true,
				resultSet: []*singleResult{
					{result: nil, err: fmt.Errorf("[PANIC]panic message")},
					{result: nil, err: fmt.Errorf("[PANIC]panic message")},
					{result: nil, err: fmt.Errorf("[PANIC]panic message")},
					{result: nil, err: fmt.Errorf("[PANIC]panic message")},
					{result: nil, err: fmt.Errorf("[PANIC]panic message")},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Start(tt.args.ctx, tt.args.executors, tt.args.concurrenty)
			assert.Equal(t, tt.want.failed, got.failed)
			assert.Equal(t, len(tt.want.resultSet), len(got.resultSet))
			for i, v := range got.resultSet {
				wantRet := tt.want.resultSet[i]
				assert.Equal(t, wantRet.result, v.result)
				assert.True(t, errors.Is(wantRet.err, v.err) || wantRet.err.Error() == v.err.Error())
			}
		})
	}
}

type executorPanic struct{}

func (e *executorPanic) Exec(ctx context.Context) (interface{}, error) {
	panic("panic message")
}

type executorError struct{}

func (e *executorError) Exec(ctx context.Context) (interface{}, error) {
	return nil, fmt.Errorf("error message")
}

type executorNormal struct{}

func (e *executorNormal) Exec(ctx context.Context) (interface{}, error) {
	return "success", nil
}
