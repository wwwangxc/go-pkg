package example_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/wwwangxc/go-pkg/concurrency"
)

func ExanpleStart(t *testing.T) {
	executors := []concurrency.Executor{
		&ExecutorImpl{arg1: "123", arg2: "456"},
		&ExecutorImpl{arg1: "123", arg2: "456"},
		&ExecutorImpl{arg1: "123", arg2: "456"},
		&ExecutorImpl{arg1: "123", arg2: "456"},
		&ExecutorImpl{arg1: "123", arg2: "456"},
	}

	result := concurrency.Start(context.Background(), executors, 2)

	// return true when no executor exec failed
	result.Succeed()

	// return true when there are executor exec failed
	result.Failed()

	// return merged error
	//
	// Format like:
	// 2 errors occurred:
	//     * error message ...
	//     * [PANIC]panic message ...
	result.MergedError()

	// return collection of all errors
	result.Errors()
}

type ExecutorImpl struct {
	arg1 string
	arg2 string
}

func (e *ExecutorImpl) Exec(ctx context.Context) (interface{}, error) {
	fmt.Println(e.arg1)
	fmt.Println(e.arg2)
	return nil, nil
}
