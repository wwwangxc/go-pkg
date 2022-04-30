package concurrency

import (
	"bytes"
	"fmt"
	"sync"
)

// Result set
type Result struct {
	resultSet []*singleResult
	failed    bool

	sync.RWMutex
}

// Succeed return true when no failed result
func (r *Result) Succeed() bool {
	r.RLock()
	defer r.RUnlock()

	return !r.failed
}

// Failed return true when there are failed result
func (r *Result) Failed() bool {
	r.RLock()
	defer r.RUnlock()

	return r.failed
}

// MergedError return merged error
//
// Format like:
// 	2 errors occurred:
// 		* error message ...
// 		* [PANIC]panic message ...
func (r *Result) MergedError() error {
	if r.Succeed() {
		return nil
	}

	r.RLock()
	defer r.RUnlock()

	errNum := 0
	errMsg := bytes.NewBufferString("")
	for _, v := range r.resultSet {
		if v.err == nil {
			continue
		}

		errNum++
		_, _ = errMsg.WriteString(fmt.Sprintf("\n    * %v", v.err))
	}

	return fmt.Errorf("%d errors occurred:%s", errNum, errMsg.String())
}

// Errors return all error
//
// return nil when no failed result
func (r *Result) Errors() []error {
	if r.Succeed() {
		return nil
	}

	r.RLock()
	defer r.RUnlock()

	errs := make([]error, 0, len(r.resultSet))
	for _, v := range r.resultSet {
		if v.err == nil {
			continue
		}

		errs = append(errs, v.err)
	}

	return errs
}

func (r *Result) append(result interface{}, err error) {
	r.Lock()
	defer r.Unlock()

	if err != nil {
		r.failed = true
	}

	r.resultSet = append(r.resultSet, &singleResult{
		result: result,
		err:    err,
	})
}

type singleResult struct {
	result interface{}
	err    error
}
