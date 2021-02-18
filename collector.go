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

	wgCount int
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
	c.wgCount++
	c.wg.Add(1)

	go func() {
		c.work <- work
	}()
}

// Wait until all task completed
func (c *Collector) Wait() {
	c.wg.Wait()
}

// Stop all ongoing task
func (c *Collector) Stop() {
	c.stopper()
	for i := 0; i < c.wgCount; i++ {
		c.wg.Done()
	}
}
