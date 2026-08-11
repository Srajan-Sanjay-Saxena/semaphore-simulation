package semaphore

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestAcquireRelease(t *testing.T) {
	s := New(5)
	s.Acquire(3)
	if s.pool != 2 {
		t.Fatalf("expected pool=2, got %d", s.pool)
	}
	s.Release(3)
	if s.pool != 5 {
		t.Fatalf("expected pool=5, got %d", s.pool)
	}
}

func TestMaxConcurrency(t *testing.T) {
	const pool = 4
	s := New(pool)

	var active atomic.Int32
	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Acquire(2)
			cur := active.Add(1)
			if cur > pool/2 {
				t.Errorf("concurrency exceeded: %d active (max 2)", cur)
			}
			active.Add(-1)
			s.Release(2)
		}()
	}
	wg.Wait()
}

func TestWeightedAcquire(t *testing.T) {
	s := New(10)
	s.Acquire(5)
	s.Acquire(5)
	if s.pool != 0 {
		t.Fatalf("expected pool=0, got %d", s.pool)
	}
	s.Release(5)
	s.Release(5)
	if s.pool != 10 {
		t.Fatalf("expected pool=10, got %d", s.pool)
	}
}

// goroutine blocked on Acquire unblocks once enough tokens are released
func TestBlocksUntilRelease(t *testing.T) {
	s := New(4)
	s.Acquire(4) // drain pool

	done := make(chan struct{})
	go func() {
		s.Acquire(2)
		close(done)
		s.Release(2)
	}()

	s.Release(4) // unblock the goroutine

	<-done
}

func TestConcurrentMixedWeights(t *testing.T) {
	s := New(10)
	var wg sync.WaitGroup

	for i := range 20 {
		wg.Add(1)
		w := uint8(i%3 + 1) // weights: 1, 2, 3
		go func(weight uint8) {
			defer wg.Done()
			s.Acquire(weight)
			s.Release(weight)
		}(w)
	}
	wg.Wait()

	if s.pool != 10 {
		t.Fatalf("expected pool restored to 10, got %d", s.pool)
	}
}
