package workermanager

// WorkerManager acts as a centralized worker pool manager and
// entrypoint for application to communicate with pool of workers.
//
// The manager holds a pool of workers, the capacity of which is determined when the manager is initialized.
// The manager dispatchs workers when it's starting.
//
// The manager holds a queue of tasks, the capacity of which is determined when the manager is initialized.
// caller can add task(s) to the queue using [WorkerManager#Assign()].
//
// The manager controls the lifecycle of its workers. It will stop all workers when the manager is being stopped.
// It will also drain any queued tasks.
//
// worker struct is not exposed to the caller by design.
type WorkerManager struct {
	workersPool chan *worker
	tasksQueue  chan Task

	errorHandler          errorHandler    // used when task return any error
	resultHandler         resultHandler   // used when task return nil error
	workerIDGeneratorFunc idGeneratorFunc // used by worker

	stopCh     chan bool
	isStopping bool
}

type workerManagerOption func(*WorkerManager)

// WithErrorHandler set specified error handler to the worker manager.
func WithErrorHandler(errHandler errorHandler) workerManagerOption {
	return func(wm *WorkerManager) {
		wm.errorHandler = errHandler
	}
}

// WithResultHandler set specified result handler to the worker manager.
func WithResultHandler(resHandler resultHandler) workerManagerOption {
	return func(wm *WorkerManager) {
		wm.resultHandler = resHandler
	}
}

func WithWorkerIDGenerator(idGenFunc idGeneratorFunc) workerManagerOption {
	return func(wm *WorkerManager) {
		wm.workerIDGeneratorFunc = idGenFunc
	}
}

// NewWorkerManager initialized new worker manager.
// There are no default value for maximum workers and tasks queue, caller needs to define it.
// Outside of those two configuration, all default value is usable.
// Any additional configuration can be specified via opts.
func NewWorkerManager(maxWorkers, maxTasksQueue int, opts ...workerManagerOption) *WorkerManager {
	wm := &WorkerManager{
		workersPool: make(chan *worker, maxWorkers),
		tasksQueue:  make(chan Task, maxTasksQueue),
		stopCh:      make(chan bool, 1),
	}

	for _, opt := range opts {
		opt(wm)
	}

	return wm
}

// Start start the worker manager.
// it will spawn worker(s) as much as the capacity of the workers pool, and start them all.
// Start will block until caller stop the manager using Stop().
func (wm *WorkerManager) Start() {
	workers := make([]*worker, cap(wm.workersPool))
	workerOpts := wm.workerOptions()
	for i := 0; i < len(workers); i++ {
		workers[i] = NewWorker(workerOpts...)
		workers[i].Start(wm.workersPool)
	}

Loop:
	for {
		select {
		case task := <-wm.tasksQueue:
			select {
			case w := <-wm.workersPool:
				w.Assign(task)
			case <-wm.stopCh:
				break Loop
			}
		case <-wm.stopCh:
			break Loop
		}
	}

	for _, w := range workers {
		w.Stop()
	}

	for _, w := range workers {
		<-w.DoneCh()
	}

	// drain the task queue
	for range wm.tasksQueue {
		if len(wm.tasksQueue) == 0 {
			break
		}
	}

	wm.isStopping = false
}

// Stop stops the manager.
func (wm *WorkerManager) Stop() {
	if !wm.isStopping {
		wm.isStopping = true
		wm.stopCh <- true
	}
}

// AddTask add task into tasks queue. Will blocked if tasksQueue is full.
func (wm *WorkerManager) AddTask(task Task) {
	if !wm.isStopping {
		wm.tasksQueue <- task
	}
}

func (wm *WorkerManager) workerOptions() []workerOption {
	workerOpts := []workerOption{}
	if wm.errorHandler != nil {
		workerOpts = append(workerOpts, workerWithErrorHandler(wm.errorHandler))
	}
	if wm.resultHandler != nil {
		workerOpts = append(workerOpts, workerWithResultHandler(wm.resultHandler))
	}
	if wm.workerIDGeneratorFunc != nil {
		workerOpts = append(workerOpts, workerWithIDGenerator(wm.workerIDGeneratorFunc))
	}

	return workerOpts
}
