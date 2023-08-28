package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	errCount := int64(0)
	var extErr error
	taskCh := make(chan Task)

	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for taskFunc := range taskCh {
				err := taskFunc()
				if err != nil {
					atomic.AddInt64(&errCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		taskCh <- task
		if int(atomic.LoadInt64(&errCount)) >= m && m > 0 {
			extErr = ErrErrorsLimitExceeded
			break
		}
	}
	close(taskCh)
	wg.Wait()

	return extErr
}
