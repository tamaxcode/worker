package worker

import (
	"sync"
	"testing"
	"time"
)

func Test_Run(t *testing.T) {
	collector := NewCollector()

	c := 0
	mtx := sync.Mutex{}
	for i := 0; i < 1000; i++ {
		collector.AddWork(Work{
			Handler: func() {
				mtx.Lock()
				c++
				mtx.Unlock()
			},
		})
	}

	collector.Wait()

	if c != 1000 {
		t.Errorf("Count Want: 1000, Got: %d\n", c)
	}
}

func Test_Stop(t *testing.T) {

	n := 1
	collector := NewCollector(&Config{
		NoOfWorkers: &n,
	})

	c := 0
	mtx := sync.Mutex{}
	for i := 0; i < 100; i++ {
		collector.AddWork(Work{
			Handler: func() {
				mtx.Lock()
				c++
				mtx.Unlock()

				// simulate heavy work
				time.Sleep(10 * time.Second)
			},
		})
	}

	go func() {
		time.Sleep(1 * time.Second)
		collector.Stop()
	}()

	collector.Wait()

	if c != 1 {
		t.Errorf("Count Want: 1, Got: %d\n", c)
	}
}

func Benchmark_Run(b *testing.B) {
	collector := NewCollector()

	c := 0
	mtx := sync.Mutex{}
	for i := 0; i < 1000; i++ {
		collector.AddWork(Work{
			Handler: func() {
				mtx.Lock()
				c++
				mtx.Unlock()
			},
		})
	}

	collector.Wait()
}
