package worker

import (
	"sync"
	"testing"
)

func Test_Run(t *testing.T) {
	c := 0
	mtx := sync.Mutex{}
	fn := func(_ interface{}) bool {
		mtx.Lock()
		c++
		mtx.Unlock()

		return true
	}

	collector := NewCollector(Config{
		Handler: fn,
	})

	for i := 0; i < 1000; i++ {
		collector.Add(nil)
	}

	collector.Wait()

	if c != 1000 {
		t.Errorf("Count Want: 1000, Got: %d\n", c)
	}
}

func Test_Concurrent(t *testing.T) {
	want := make([]int, 1000)
	got := make([]int, 1000)

	type payload struct {
		n int
		i int
	}

	collector := NewCollector(Config{
		Handler: func(data interface{}) bool {
			p, _ := data.(payload)

			got[p.i] = p.n

			return true
		},
	})

	for i := 0; i < 1000; i++ {
		want[i] = i

		collector.Add(payload{
			n: i,
			i: i,
		})
	}

	collector.Wait()

	for i, n := range want {
		if got[i] != n {
			t.Errorf("At index %d, WANT: %d, GOT: %d", i, n, got[i])
		}
	}
}

func Test_Worker_Stop(t *testing.T) {

	c := 0
	mtx := sync.Mutex{}
	fn := func(data interface{}) bool {
		mtx.Lock()
		defer mtx.Unlock()

		if c == 2 {
			return false
		}

		c++

		return true
	}

	collector := NewCollector(Config{
		NoOfWorkers: 2,
		Handler:     fn,
	})

	for i := 0; i < 100; i++ {
		collector.Add(nil)
	}

	collector.Wait()

	if c != 2 {
		t.Errorf("Count Want: 2, Got: %d\n", c)
	}
}

func Benchmark_Run1000(b *testing.B) {
	c := 0
	mtx := sync.Mutex{}

	collector := NewCollector(Config{
		Handler: func(data interface{}) bool {
			mtx.Lock()
			c++
			mtx.Unlock()

			return true
		},
	})

	for i := 0; i < 1000; i++ {
		collector.Add(nil)
	}

	collector.Wait()
}
