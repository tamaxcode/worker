package worker

import (
	"context"
	"sync"
)

// Collector will communicate with outer layer to receive handler for worker to process task
type Collector struct {
	workers []*Worker
	work    chan<- Work
	wg      *sync.WaitGroup
	ctx     context.Context

	stopper func()
}

// NewCollector ...
func NewCollector(configs ...*Config) *Collector {

	ctx, canc := context.WithCancel(context.Background())
	cfg := mergeOrDefault(configs)

	wg := sync.WaitGroup{}

	inputChannel := make(chan Work)
	workers := make([]*Worker, 0)

	for i := 0; i < *cfg.NoOfWorkers; i++ {
		w := createWorker(ctx, inputChannel, &wg)
		w.Start()

		workers = append(workers, w)
	}

	return &Collector{
		work:    inputChannel,
		workers: workers,
		wg:      &wg,

		ctx:     ctx,
		stopper: canc,
	}
}

// AddWork ...
// Example
/*
	Collector.AddWork(Work{
		Handler: func() {
			SomeHandler(someData)
		}
	})
*/
func (c *Collector) AddWork(work Work) {
	c.wg.Add(1)

	go func() {
		c.work <- work
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
	case <-done:
		// wait group complete
	}
}

// Stop all ongoing task
func (c *Collector) Stop() {
	c.stopper()
}
