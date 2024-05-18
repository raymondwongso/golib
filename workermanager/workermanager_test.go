package workermanager_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/raymondwongso/golib/workermanager"
	"github.com/stretchr/testify/require"
)

type testTask struct {
	isError bool
}

type testResult struct{}

func (t *testTask) Execute() (any, error) {
	if t.isError {
		return nil, errors.New("some error")
	}

	return &testResult{}, nil
}

var (
	testErrHandler = func(err error) {}
	testResHandler = func(res any) {}
)

func Test_New(t *testing.T) {
	opts := []workermanager.WorkerManagerOption{
		workermanager.WithErrorHandler(testErrHandler),
		workermanager.WithResultHandler(testResHandler),
	}

	require.NotPanics(t, func() {
		_ = workermanager.New(5, 5, opts...)
	})
}

// The strategy to test plain worker manager:
// 1. Initialized counter variable
// 2. Initialized worker manager with result handler that increment the counter
// 3. After the counter reach targetThreshold, stop the worker manager
// 4. if all goes well, the test should end without any panic or error
// 5. done channel is used to verify the continuation of the test
// 6. Timeout or panic means there is a possibility of bug.
//
// This test also verify that result handler is executed accordingly.
func Test_PlainWorkerManager_Start(t *testing.T) {
	var (
		wm                *workermanager.WorkerManager
		counter           = 0
		targetThreshold   = 5
		pool              = 1
		tasksQueue        = 10
		mut               = new(sync.Mutex)
		resCounterHandler = func(res any) {
			mut.Lock()
			counter += 1
			mut.Unlock()
			if counter == targetThreshold {
				wm.Stop()
			}
		}
		opts = []workermanager.WorkerManagerOption{
			workermanager.WithResultHandler(resCounterHandler),
		}
	)

	require.NotPanics(t, func() {
		wm = workermanager.New(pool, tasksQueue, opts...)

		done := make(chan bool, 1)
		go func() {
			wm.Start()
			done <- true
		}()

		for i := 0; i < tasksQueue; i++ {
			wm.AddTask(&testTask{})
		}

		<-done
	})
}

// The strategy to test error worker manager:
// 1. Initialized error handler that will stop the worker manager if threshold reached.
// 2. task will return error.
// 3. start worker manager with error handler and using the "vaulty" task
// 4. test should reach done without timeout or panics.
//
// This test verifies that error handler is executed accordingly.
func Test_WithErrorWorkerManager_Start(t *testing.T) {
	var (
		wm                *workermanager.WorkerManager
		counter           = 0
		targetError       = 5
		mut               = new(sync.Mutex)
		errCounterHandler = func(err error) {
			mut.Lock()
			counter += 1
			mut.Unlock()
			if counter == targetError {
				wm.Stop()
			}
		}
		opts = []workermanager.WorkerManagerOption{
			workermanager.WithErrorHandler(errCounterHandler),
		}
	)

	wm = workermanager.New(1, targetError, opts...)

	require.NotPanics(t, func() {
		done := make(chan bool, 1)
		go func() {
			wm.Start()
			done <- true
		}()

		for i := 0; i < targetError; i++ {
			wm.AddTask(&testTask{isError: false})
		}

		<-done
	})
}
