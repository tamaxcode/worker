package worker

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

// Work will bring task handler for worker to call
type Work struct {
	Handler WorkHandler
}

// WorkHandler function that will be called by the worker
// example
/*
	handler := func () {
		Handler(someData)
	}
*/
type WorkHandler func()

// Worker ...
type Worker struct {
	id  string
	ctx context.Context

	task <-chan Work
	wg   *sync.WaitGroup
}

func createWorker(ctx context.Context, taskChannel <-chan Work, wg *sync.WaitGroup) *Worker {
	id := uuid.NewString()

	return &Worker{
		id:   id,
		ctx:  ctx,
		task: taskChannel,
		wg:   wg,
	}
}

// Start ...
func (w *Worker) Start() {

	go func() {
		for {
			var handler WorkHandler

			select {
			case <-w.ctx.Done():
				// context canceled, stop worker
				return
			case work, ok := <-w.task:
				if !ok {
					// channel closed, stop worker
					return
				}
				handler = work.Handler
			}

			handler()

			w.wg.Done()
		}
	}()
}
