package concurrency

import (
	"context"
	"fmt"
	"sync"
)

// Start run concurrenty
//
// Result is collection of all executors results and error
func Start(ctx context.Context, executors []Executor, concurrenty uint8) *Result {
	result := &Result{}
	limitCh := make(chan struct{}, concurrenty)
	var wg sync.WaitGroup
	wg.Add(len(executors))

	go func() {
		for _, v := range executors {
			limitCh <- struct{}{}
			go func(executor Executor) {
				defer func() {
					<-limitCh
					if e := recover(); e != nil {
						result.append(nil, fmt.Errorf("[PANIC]%v", e))
					}
					wg.Done()
				}()

				result.append(executor.Exec(ctx))
			}(v)
		}
	}()

	wg.Wait()
	return result
}
