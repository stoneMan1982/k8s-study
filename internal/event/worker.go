package event

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Processor defines the function signature for processing events
type Processor func(event Event) error

// WorkerPool manages a pool of workers that process events concurrently
type WorkerPool struct {
	eventChan  chan Event
	processor  Processor
	workerNum  int
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	once       sync.Once
	errorChan  chan error
	started    bool
	mu         sync.RWMutex
}

// WorkerPoolConfig holds configuration for the worker pool
type WorkerPoolConfig struct {
	WorkerNum      int           // Number of concurrent workers
	QueueSize      int           // Size of the event queue buffer
	ErrorChanSize  int           // Size of the error channel buffer
	Processor      Processor     // Event processor function
}

// DefaultConfig returns default worker pool configuration
func DefaultConfig() WorkerPoolConfig {
	return WorkerPoolConfig{
		WorkerNum:     5,
		QueueSize:     100,
		ErrorChanSize: 10,
	}
}

// NewWorkerPool creates a new worker pool with the given configuration
func NewWorkerPool(config WorkerPoolConfig) *WorkerPool {
	if config.WorkerNum <= 0 {
		config.WorkerNum = 5
	}
	if config.QueueSize <= 0 {
		config.QueueSize = 100
	}
	if config.ErrorChanSize <= 0 {
		config.ErrorChanSize = 10
	}
	if config.Processor == nil {
		config.Processor = defaultProcessor
	}

	ctx, cancel := context.WithCancel(context.Background())
	
	return &WorkerPool{
		eventChan: make(chan Event, config.QueueSize),
		processor: config.Processor,
		workerNum: config.WorkerNum,
		errorChan: make(chan error, config.ErrorChanSize),
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Start begins processing events with the worker pool
func (wp *WorkerPool) Start() {
	wp.mu.Lock()
	if wp.started {
		wp.mu.Unlock()
		return
	}
	wp.started = true
	wp.mu.Unlock()

	// Start worker goroutines
	for i := 0; i < wp.workerNum; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker is the main loop for each worker goroutine
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	
	for {
		select {
		case event, ok := <-wp.eventChan:
			if !ok {
				// Channel closed, exit worker
				return
			}
			
			// Record when we start handling the event
			event.HandleAt = time.Now()
			
			// Process the event
			if err := wp.processor(event); err != nil {
				// Send error to error channel (non-blocking)
				select {
				case wp.errorChan <- fmt.Errorf("worker %d failed to process event %s: %w", id, event.Type, err):
				default:
					// Error channel is full, skip
				}
			}
			
		case <-wp.ctx.Done():
			// Context cancelled, exit worker
			return
		}
	}
}

// Submit adds an event to the processing queue
func (wp *WorkerPool) Submit(event Event) error {
	wp.mu.RLock()
	if !wp.started {
		wp.mu.RUnlock()
		return errors.New("worker pool not started")
	}
	wp.mu.RUnlock()

	event.SendAt = time.Now()
	
	select {
	case wp.eventChan <- event:
		return nil
	case <-wp.ctx.Done():
		return errors.New("worker pool is shutting down")
	default:
		return errors.New("event queue is full")
	}
}

// SubmitWithTimeout adds an event with a timeout
func (wp *WorkerPool) SubmitWithTimeout(event Event, timeout time.Duration) error {
	wp.mu.RLock()
	if !wp.started {
		wp.mu.RUnlock()
		return errors.New("worker pool not started")
	}
	wp.mu.RUnlock()

	event.SendAt = time.Now()
	
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	
	select {
	case wp.eventChan <- event:
		return nil
	case <-timer.C:
		return errors.New("submit timeout")
	case <-wp.ctx.Done():
		return errors.New("worker pool is shutting down")
	}
}

// Stop gracefully stops the worker pool
func (wp *WorkerPool) Stop() {
	wp.once.Do(func() {
		wp.mu.Lock()
		wp.started = false
		wp.mu.Unlock()
		
		// Close the event channel to signal workers to finish current work
		close(wp.eventChan)
		
		// Wait for all workers to complete
		wp.wg.Wait()
		
		// Cancel context
		wp.cancel()
		
		// Close error channel
		close(wp.errorChan)
	})
}

// StopWithTimeout stops the worker pool with a timeout
func (wp *WorkerPool) StopWithTimeout(timeout time.Duration) error {
	done := make(chan struct{})
	go func() {
		wp.Stop()
		close(done)
	}()
	
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	
	select {
	case <-done:
		return nil
	case <-timer.C:
		wp.cancel() // Force cancel
		return errors.New("stop timeout, forced shutdown")
	}
}

// Errors returns the error channel for monitoring worker errors
func (wp *WorkerPool) Errors() <-chan error {
	return wp.errorChan
}

// QueueLen returns the current number of events in the queue
func (wp *WorkerPool) QueueLen() int {
	return len(wp.eventChan)
}

// IsStarted returns whether the worker pool is started
func (wp *WorkerPool) IsStarted() bool {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.started
}

// defaultProcessor is a no-op processor used when none is provided
func defaultProcessor(event Event) error {
	// Default implementation does nothing
	return nil
}
