package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	glworkermanager "github.com/raymondwongso/golib/workermanager"
)

var (
	wm *glworkermanager.WorkerManager
)

type someTask struct {
	ID int
}

type someResult struct {
	ID    int
	Value int
}

func (st *someTask) Execute() (any, error) {
	if st.ID%2 == 0 {
		return nil, errors.New("task error")
	}
	return &someResult{st.ID, st.ID * 2}, nil
}

func someErrorHandler(err error) {
	if err != nil {
		fmt.Printf("encountered error: %s\n", err)
		// Example when you want to stop the worker manager:
		// wm.Stop()
	}
}

func someResultHandler(res any) {
	if res, ok := res.(*someResult); ok {
		fmt.Printf("result: %#v\n", res)
	}
}

func main() {
	opts := []glworkermanager.WorkerManagerOption{
		glworkermanager.WithErrorHandler(someErrorHandler),
		glworkermanager.WithResultHandler(someResultHandler),
	}
	wm = glworkermanager.NewWorkerManager(5, 5, opts...)

	intr := make(chan os.Signal, 1)
	signal.Notify(intr, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)
	go func() {
		wm.Start()
		done <- true
	}()

	for i := 0; i < 15; i++ {
		wm.AddTask(&someTask{ID: i})
	}

	<-intr
	wm.Stop()
	<-done

	fmt.Println("worker manager done")
}
