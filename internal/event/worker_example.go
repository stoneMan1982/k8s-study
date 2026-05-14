package event

import (
	"fmt"
	"log"
	"time"
)

// Example demonstrates how to use the WorkerPool
func Example() {
	// Define a custom processor
	processor := func(event Event) error {
		fmt.Printf("[%s] Processing event: %s from %s\n",
			time.Now().Format("15:04:05"),
			event.Type,
			event.From)

		// Simulate some work
		time.Sleep(100 * time.Millisecond)

		// Simulate occasional errors
		if event.Type == "error-event" {
			return fmt.Errorf("failed to process error event")
		}

		return nil
	}

	// Create worker pool with custom configuration
	config := WorkerPoolConfig{
		WorkerNum:     3,  // 3 concurrent workers
		QueueSize:     50, // Queue can hold 50 events
		ErrorChanSize: 10, // Error channel buffer
		Processor:     processor,
	}

	pool := NewWorkerPool(config)

	// Start error monitoring goroutine
	go func() {
		for err := range pool.Errors() {
			log.Printf("Worker error: %v", err)
		}
	}()

	// Start the pool
	pool.Start()

	// Submit some events
	events := []Event{
		{Type: "user.created", From: "api", Payload: map[string]string{"userId": "123"}},
		{Type: "user.updated", From: "api", Payload: map[string]string{"userId": "123"}},
		{Type: "user.deleted", From: "api", Payload: map[string]string{"userId": "123"}},
		{Type: "order.placed", From: "web", Payload: map[string]string{"orderId": "456"}},
		{Type: "error-event", From: "test", Payload: nil},
	}

	for _, event := range events {
		if err := pool.Submit(event); err != nil {
			log.Printf("Failed to submit event: %v", err)
		}
	}

	// Submit with timeout
	urgentEvent := Event{
		Type:    "urgent.notification",
		From:    "system",
		Payload: "Critical alert",
	}
	if err := pool.SubmitWithTimeout(urgentEvent, 1*time.Second); err != nil {
		log.Printf("Failed to submit urgent event: %v", err)
	}

	// Let workers process
	time.Sleep(2 * time.Second)

	// Check queue length
	fmt.Printf("Remaining events in queue: %d\n", pool.QueueLen())

	// Graceful shutdown
	fmt.Println("Shutting down worker pool...")
	pool.Stop()
	// Or with timeout: pool.StopWithTimeout(5 * time.Second)

	fmt.Println("Worker pool stopped")
}

// ExampleWithDefaultConfig shows using default configuration
func ExampleWithDefaultConfig() {
	// Use default configuration
	config := DefaultConfig()
	config.Processor = func(event Event) error {
		fmt.Printf("Processing: %s\n", event.Type)
		return nil
	}

	pool := NewWorkerPool(config)
	pool.Start()
	defer pool.Stop()

	// Submit events
	pool.Submit(Event{Type: "test.event", From: "example"})
}
