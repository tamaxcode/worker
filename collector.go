package worker

import (
	"context"
	"runtime"
	"sync"
)

// Config represent configuration setting for collector to send received task to worker
type Config struct {
	// Number of worker spawn to handle concurrent work. Default to "runtime.NumCPU"
	// giving value <= 0 will fallback to default value
	NoOfWorkers int
	Handler     func(data interface{}) bool
}

// Collector will communicate with outer layer to receive handler for worker to process task
type Collector struct {
	workers []*Worker
	work    chan<- interface{}
	wg      *sync.WaitGroup
	ctx     context.Context

	stopChan <-chan bool
	stopper  func()
}

// NewCollector ...
func NewCollector(config Config) *Collector {

	ctx, canc := context.WithCancel(context.Background())

	noOfWorkers := config.NoOfWorkers
	if noOfWorkers <= 0 {
		noOfWorkers = runtime.NumCPU()
	}

	workers := make([]*Worker, noOfWorkers)

	inputChannel := make(chan interface{}, noOfWorkers)
	wg := sync.WaitGroup{}

	stopChan := make(chan bool)

	for i := 0; i < noOfWorkers; i++ {
		w := Worker{
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

	return &Collector{
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
func (c *Collector) Add(data interface{}) {
	c.wg.Add(1)

	go func() {
		c.work <- data
	}()
}

// Wait until all task completed
func (c *Collector) Wait() {
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
func (c *Collector) Stop() {
	c.stopper()
}
