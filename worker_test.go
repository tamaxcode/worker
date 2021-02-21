package worker

import (
	"sync"
	"testing"
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

	cfg := &Config{}
	cfg.SetNoOfWorkers(2)
	collector := NewCollector(cfg)

	c := 0
	mtx := sync.Mutex{}
	for i := 0; i < 100; i++ {
		collector.AddWork(Work{
			Handler: func() {
				mtx.Lock()
				defer mtx.Unlock()

				if c == 2 {
					collector.Stop()
					return
				}

				c++
			},
		})
	}

	collector.Wait()

	if c != 2 {
		t.Errorf("Count Want: 2, Got: %d\n", c)
	}
}

func Benchmark_Run1000(b *testing.B) {
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
