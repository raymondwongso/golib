package workermanager

// Task defines an interface for an object to be considered as executable by worker.
type Task interface {
	Execute() (any, error)
}

// Worker defines singular worker that execute the task.
type Worker struct {
	ID     string `json:"id,omitempty"`
	taskCh chan Task
	stopCh chan bool

	resultHandler resultHandler // typically passed from manager
	errorHandler  errorHandler  // typically passed from manager

	doneCh chan bool
}

type workerOption func(*Worker)

func workerWithResultHandler(resHandler resultHandler) workerOption {
	return func(w *Worker) {
		w.resultHandler = resHandler
	}
}

func workerWithErrorHandler(errHandler errorHandler) workerOption {
	return func(w *Worker) {
		w.errorHandler = errHandler
	}
}

// NewWorker creates worker with various options.
// These are available options:
//   - resultHandler
//   - errorHandler
func NewWorker(opts ...workerOption) *Worker {
	w := &Worker{
		taskCh: make(chan Task, 1),
		stopCh: make(chan bool, 1),
		doneCh: make(chan bool, 1),
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

// Assign assigns task to the task channel of the worker.
func (w *Worker) Assign(task Task) {
	w.taskCh <- task
}

// Start makes worker register itself into specified workersPool
// and then listen to the task channel.
// worker will execute task received from the channel and will invoke
// errorHandler or resultHandler if any.
func (w *Worker) Start(workersPool chan *Worker) {
	go func() {
	Loop:
		for {
			workersPool <- w

			select {
			case task := <-w.taskCh:
				res, err := task.Execute()
				if err != nil && w.errorHandler != nil {
					w.errorHandler(err)
				} else if w.resultHandler != nil {
					w.resultHandler(res)
				}
			case <-w.stopCh:
				w.doneCh <- true
				break Loop
			}
		}
	}()
}

// DoneCh return doneChannel of the worker
func (w *Worker) DoneCh() chan bool {
	return w.doneCh
}

// Stop stop the worker and notify the stop channel.
func (w *Worker) Stop() {
	w.stopCh <- true
}
