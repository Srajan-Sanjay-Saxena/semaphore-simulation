package main

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	cond    *sync.Cond
	count   int
	maxSlots int
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		cond:     sync.NewCond(&sync.Mutex{}),
		maxSlots: n,
	}
}

func (s *Semaphore) Acquire() {
	s.cond.L.Lock()
	for s.count >= s.maxSlots {
		s.cond.Wait()
	}
	s.count++
	s.cond.L.Unlock()
}

func (s *Semaphore) Release() {
	s.cond.L.Lock()
	s.count--
	s.cond.Signal()
	s.cond.L.Unlock()
}

func main() {
	const (
		workers  = 10
		maxSlots = 3
	)

	sem := NewSemaphore(maxSlots)
	var wg sync.WaitGroup

	for i := range workers {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sem.Acquire()
			fmt.Printf("worker %d acquired slot\n", id)
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("worker %d releasing slot\n", id)
			sem.Release()
		}(i)
	}

	wg.Wait()
}
