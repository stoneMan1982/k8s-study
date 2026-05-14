package event

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPool_BasicFunctionality(t *testing.T) {
	var processedCount int32
	processor := func(event Event) error {
		atomic.AddInt32(&processedCount, 1)
		return nil
	}

	config := WorkerPoolConfig{
		WorkerNum: 3,
		QueueSize: 10,
		Processor: processor,
	}

	pool := NewWorkerPool(config)
	pool.Start()
	defer pool.Stop()

	// Submit 10 events
	for i := 0; i < 10; i++ {
		event := Event{Type: "test", From: "test"}
		if err := pool.Submit(event); err != nil {
			t.Errorf("Failed to submit event: %v", err)
		}
	}

	// Wait for processing
	time.Sleep(200 * time.Millisecond)

	count := atomic.LoadInt32(&processedCount)
	if count != 10 {
		t.Errorf("Expected 10 events processed, got %d", count)
	}
}

func TestWorkerPool_ErrorHandling(t *testing.T) {
	processor := func(event Event) error {
		if event.Type == "error" {
			return errors.New("simulated error")
		}
		return nil
	}

	config := WorkerPoolConfig{
		WorkerNum:     2,
		QueueSize:     10,
		ErrorChanSize: 5,
		Processor:     processor,
	}

	pool := NewWorkerPool(config)
	pool.Start()
	defer pool.Stop()

	var errorCount int32
	go func() {
		for range pool.Errors() {
			atomic.AddInt32(&errorCount, 1)
		}
	}()

	// Submit events that will cause errors
	for i := 0; i < 3; i++ {
		pool.Submit(Event{Type: "error", From: "test"})
	}

	// Submit successful events
	for i := 0; i < 2; i++ {
		pool.Submit(Event{Type: "success", From: "test"})
	}

	time.Sleep(200 * time.Millisecond)

	count := atomic.LoadInt32(&errorCount)
	if count != 3 {
		t.Errorf("Expected 3 errors, got %d", count)
	}
}

func TestWorkerPool_ConcurrentSubmit(t *testing.T) {
	var processedCount int32
	processor := func(event Event) error {
		atomic.AddInt32(&processedCount, 1)
		time.Sleep(10 * time.Millisecond)
		return nil
	}

	config := WorkerPoolConfig{
		WorkerNum: 5,
		QueueSize: 100,
		Processor: processor,
	}

	pool := NewWorkerPool(config)
	pool.Start()
	defer pool.Stop()

	// Concurrently submit events from multiple goroutines
	var wg sync.WaitGroup
	numGoroutines := 10
	eventsPerGoroutine := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < eventsPerGoroutine; j++ {
				event := Event{Type: "concurrent", From: "test"}
				if err := pool.Submit(event); err != nil {
					t.Errorf("Failed to submit event: %v", err)
				}
			}
		}()
	}

	wg.Wait()
	time.Sleep(2 * time.Second)

	expected := int32(numGoroutines * eventsPerGoroutine)
	count := atomic.LoadInt32(&processedCount)
	if count != expected {
		t.Errorf("Expected %d events processed, got %d", expected, count)
	}
}

func TestWorkerPool_SubmitWithTimeout(t *testing.T) {
	// Create a pool with no workers to ensure queue stays full
	config := WorkerPoolConfig{
		WorkerNum: 0, // No workers to process events
		QueueSize: 2,
		Processor: func(event Event) error { return nil },
	}

	pool := NewWorkerPool(config)
	// Don't start the pool - no workers will be spawned
	pool.started = true // Manually mark as started to allow submits
	defer pool.Stop()

	// Fill the queue completely
	pool.Submit(Event{Type: "test1", From: "test"})
	pool.Submit(Event{Type: "test2", From: "test"})

	// Now queue is full, timeout should occur
	event := Event{Type: "test3", From: "test"}
	err := pool.SubmitWithTimeout(event, 50*time.Millisecond)

	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
	if err != nil && err.Error() != "submit timeout" {
		t.Errorf("Expected 'submit timeout' error, got: %v", err)
	}
}

func TestWorkerPool_Stop(t *testing.T) {
	var processedCount int32
	processor := func(event Event) error {
		atomic.AddInt32(&processedCount, 1)
		time.Sleep(50 * time.Millisecond)
		return nil
	}

	config := WorkerPoolConfig{
		WorkerNum: 2,
		QueueSize: 10,
		Processor: processor,
	}

	pool := NewWorkerPool(config)
	pool.Start()

	// Submit some events
	for i := 0; i < 5; i++ {
		pool.Submit(Event{Type: "test", From: "test"})
	}

	// Stop immediately
	pool.Stop()

	// Try to submit after stop
	err := pool.Submit(Event{Type: "after-stop", From: "test"})
	if err == nil {
		t.Error("Expected error when submitting to stopped pool")
	}

	// Verify pool is stopped
	if pool.IsStarted() {
		t.Error("Pool should not be started after Stop()")
	}
}

func TestWorkerPool_StopWithTimeout(t *testing.T) {
	processor := func(event Event) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	config := WorkerPoolConfig{
		WorkerNum: 1,
		QueueSize: 10,
		Processor: processor,
	}

	pool := NewWorkerPool(config)
	pool.Start()

	// Submit events
	for i := 0; i < 3; i++ {
		pool.Submit(Event{Type: "test", From: "test"})
	}

	// Stop with short timeout (should timeout)
	err := pool.StopWithTimeout(50 * time.Millisecond)
	if err == nil {
		t.Error("Expected timeout error")
	}
}

func TestWorkerPool_QueueLen(t *testing.T) {
	processor := func(event Event) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	config := WorkerPoolConfig{
		WorkerNum: 1,
		QueueSize: 10,
		Processor: processor,
	}

	pool := NewWorkerPool(config)
	pool.Start()
	defer pool.Stop()

	// Submit events
	for i := 0; i < 5; i++ {
		pool.Submit(Event{Type: "test", From: "test"})
	}

	// Check queue length
	queueLen := pool.QueueLen()
	if queueLen == 0 {
		t.Error("Queue length should be greater than 0")
	}
}

func TestWorkerPool_DefaultProcessor(t *testing.T) {
	config := WorkerPoolConfig{
		WorkerNum: 2,
		QueueSize: 10,
		Processor: nil, // Will use default processor
	}

	pool := NewWorkerPool(config)
	pool.Start()
	defer pool.Stop()

	// Should not panic with nil processor
	err := pool.Submit(Event{Type: "test", From: "test"})
	if err != nil {
		t.Errorf("Failed to submit with default processor: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
}

func TestWorkerPool_MultipleStops(t *testing.T) {
	config := DefaultConfig()
	config.Processor = func(event Event) error { return nil }

	pool := NewWorkerPool(config)
	pool.Start()

	// Stop multiple times should not panic
	pool.Stop()
	pool.Stop()
	pool.Stop()
}

func BenchmarkWorkerPool_Submit(b *testing.B) {
	processor := func(event Event) error {
		return nil
	}

	config := WorkerPoolConfig{
		WorkerNum: 5,
		QueueSize: 1000,
		Processor: processor,
	}

	pool := NewWorkerPool(config)
	pool.Start()
	defer pool.Stop()

	event := Event{Type: "benchmark", From: "test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.Submit(event)
	}
}

func BenchmarkWorkerPool_SubmitParallel(b *testing.B) {
	processor := func(event Event) error {
		return nil
	}

	config := WorkerPoolConfig{
		WorkerNum: 5,
		QueueSize: 10000,
		Processor: processor,
	}

	pool := NewWorkerPool(config)
	pool.Start()
	defer pool.Stop()

	event := Event{Type: "benchmark", From: "test"}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Submit(event)
		}
	})
}

func BenchmarkWorkerPool_WithProcessing(b *testing.B) {
	processor := func(event Event) error {
		// Simulate light processing
		_ = event.Type + event.From
		return nil
	}

	config := WorkerPoolConfig{
		WorkerNum: 5,
		QueueSize: 1000,
		Processor: processor,
	}

	pool := NewWorkerPool(config)
	pool.Start()
	defer pool.Stop()

	event := Event{Type: "benchmark", From: "test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.Submit(event)
	}

	pool.Stop()
}

func BenchmarkWorkerPool_CPUIntensive(b *testing.B) {
	processor := func(event Event) error {
		// Simulate CPU-intensive work
		sum := 0
		for i := 0; i < 1000; i++ {
			sum += i
		}
		return nil
	}

	config := WorkerPoolConfig{
		WorkerNum: 5,
		QueueSize: 1000,
		Processor: processor,
	}

	pool := NewWorkerPool(config)
	pool.Start()
	defer pool.Stop()

	event := Event{Type: "benchmark", From: "test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.Submit(event)
	}
}

func BenchmarkWorkerPool_DifferentWorkerCounts(b *testing.B) {
	processor := func(event Event) error {
		sum := 0
		for i := 0; i < 100; i++ {
			sum += i
		}
		return nil
	}

	workerCounts := []int{1, 2, 4, 8, 16}
	for _, workers := range workerCounts {
		b.Run(fmt.Sprintf("Workers_%d", workers), func(b *testing.B) {
			config := WorkerPoolConfig{
				WorkerNum: workers,
				QueueSize: 1000,
				Processor: processor,
			}

			pool := NewWorkerPool(config)
			pool.Start()
			defer pool.Stop()

			event := Event{Type: "benchmark", From: "test"}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				pool.Submit(event)
			}
		})
	}
}
