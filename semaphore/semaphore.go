package semaphore

import "sync"


type Semaphore struct {
	pool uint32
	cond *sync.Cond
	mu *sync.Mutex
}

func New(pool uint32) *Semaphore {
	var mu sync.Mutex
	return &Semaphore{
		pool: pool,
		cond: sync.NewCond(&mu),
		mu: &mu,
	}
}

func (s *Semaphore) Acquire(weight uint8) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for s.pool < uint32(weight) {
		s.cond.Wait()
	}

	s.pool -= uint32(weight);
}


// why are we using broadcasd instead of signal is that say there is a pool of 10 tokens, now some goroutine because of some heavy work acquire '5-5' tokens out of the pool ....now there might be the case that other goroutines that are sleeping are just waiting for 2 or 1 token .. so we do a broadcast to wake up all the waiting routines .
func (s *Semaphore) Release(weight uint8) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pool += uint32(weight)
	s.cond.Broadcast()
}