// Package workermanager implements goroutine-based, distributed worker for tasks execution.
//
// The workermanager package only expose WorkerManager to the caller by design.
// Example of the usage can be seen in [Example Worker Manager].
//
// [Example Worker Manager]: https://github.com/raymondwongso/golib/example/workermanager/main.go]
package workermanager

// define error handler func.
type errorHandler func(error)

// define result handler func.
type resultHandler func(any)
