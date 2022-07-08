package worker

import (
	"context"
	"runtime"
	"sync"
)

// Config represent configuration setting for collector to send received task to worker
type Config[T any] struct {
	// Number of worker spawn to handle concurrent work. Default to "runtime.NumCPU"
	// giving value <= 0 will fallback to default value
	NoOfWorkers int
	Handler     func(data T) bool
}

// Collector will communicate with outer layer to receive handler for worker to process task
type Collector[T any] struct {
	workers []*Worker[T]
	work    chan<- T
	wg      *sync.WaitGroup
	ctx     context.Context

	stopChan <-chan bool
	stopper  func()
}

// NewCollector ...
func NewCollector[T any](config Config[T]) *Collector[T] {

	ctx, canc := context.WithCancel(context.Background())

	noOfWorkers := config.NoOfWorkers
	if noOfWorkers <= 0 {
		noOfWorkers = runtime.NumCPU()
	}

	workers := make([]*Worker[T], noOfWorkers)

	inputChannel := make(chan T, noOfWorkers)
	wg := sync.WaitGroup{}

	stopChan := make(chan bool)

	for i := 0; i < noOfWorkers; i++ {
		w := Worker[T]{
			id:  i,
			ctx: ctx,

			workHandler: config.Handler,
			workChannel: inputChannel,
			wg:          &wg,

			stopChannel: stopChan,
		}
		w.Start()

		workers[i] = &w
	}

	return &Collector[T]{
		work:    inputChannel,
		workers: workers,
		wg:      &wg,

		ctx:      ctx,
		stopChan: stopChan,
		stopper:  canc,
	}
}

// Add ...
// Example
/*
	Collector.Add(somedata)
*/
func (c *Collector[T]) Add(data T) {
	c.wg.Add(1)

	go func() {
		c.work <- data
	}()
}

// Wait until all task completed
func (c *Collector[T]) Wait() {
	done := make(chan bool)
	go func() {
		c.wg.Wait()
		done <- true
	}()

	select {

	case <-c.ctx.Done():
		// context cancelled
	case <-c.stopChan:
		// worker abort
		c.stopper()
	case <-done:
		// wait group complete
	}
}

// Stop all ongoing task
func (c *Collector[T]) Stop() {
	c.stopper()
}
