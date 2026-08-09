# go-semaphore

A simulation of semaphore-based concurrency control in Go, demonstrating how to limit concurrent access to a shared resource using goroutines and channels.

## What it does

- Uses a buffered channel as a semaphore to cap the number of goroutines running concurrently
- Simulates multiple workers competing for a limited number of slots
- Shows how semaphores prevent race conditions and resource exhaustion

## Run

```bash
go run main.go
```

## How it works

A buffered channel of size `N` acts as the semaphore. Each worker acquires a slot before doing work and releases it after, ensuring at most `N` workers run at the same time.

```go
sem := make(chan struct{}, N)

// acquire
sem <- struct{}{}

// release
<-sem
```
