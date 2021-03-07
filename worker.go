package worker

import (
	"context"
	"sync"
)

// Work will bring task handler for worker to call
type Work struct {
	Handler WorkHandler
}

// WorkHandler function that will be called by the worker to process incoming work data
// will terminate all pending work if return false
// example
/*
	handler := func (data interface{}) bool {
		str, _ := data.(string)

		fmt.Println(str)

		return true
	}
*/
type WorkHandler func(data interface{}) bool

// Worker ...
type Worker struct {
	id  int
	ctx context.Context

	workHandler func(data interface{}) bool
	workChannel <-chan interface{}
	wg          *sync.WaitGroup

	stopChannel chan<- bool
}

// Start ...
func (w *Worker) Start() {

	go func() {
		for {
			var (
				workData interface{}
			)

			select {
			case <-w.ctx.Done():
				// context canceled, stop worker
				return
			case work, ok := <-w.workChannel:
				if !ok {
					// channel closed, stop worker
					return
				}
				workData = work
			}

			cont := true
			if w.workChannel != nil {
				cont = w.workHandler(workData)
			}

			if cont == false {
				w.stopChannel <- true
			}

			w.wg.Done()
		}
	}()
}
