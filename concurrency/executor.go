package concurrency

import "context"

// Executor concurrenty executor
type Executor interface {
	Exec(ctx context.Context) (interface{}, error)
}
