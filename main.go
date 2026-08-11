package main

import (
	"fmt"
	"sync"
	"time"

	"go-semaphore/semaphore"
)

func main() {
	const (
		workers = 10
		pool    = 6
		weight  = 2
	)

	sem := semaphore.New(pool)
	var wg sync.WaitGroup

	for i := range workers {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sem.Acquire(weight)
			fmt.Printf("worker %d acquired\n", id)
			time.Sleep(500 * time.Millisecond)
			sem.Release(weight)
			fmt.Printf("worker %d released\n", id)
		}(i)
	}

	wg.Wait()
}
