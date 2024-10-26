package dispatcher

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockHandler struct {
	called    bool
	callCount int
	mu        sync.Mutex
	wg        *sync.WaitGroup
}

func (m *MockHandler) Handle(_ context.Context, event string, payload []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.called = true
	m.callCount++
	if m.wg != nil {
		m.wg.Done() // Signal that this handler has completed
	}
}

func TestEventDispatcher_RegisterAndDispatch(t *testing.T) {
	// Create the dispatcher
	dispatcher := NewSimpleEventDispatcher(false)

	// Create a WaitGroup to wait for handlers to complete
	var wg sync.WaitGroup

	// Create mock handlers
	handler1 := &MockHandler{wg: &wg}
	handler2 := &MockHandler{wg: &wg}

	// Register handlers for the same event
	dispatcher.Register(context.Background(), "testEvent", handler1)
	dispatcher.Register(context.Background(), "testEvent", handler2)

	// Set the WaitGroup counter to the number of handlers
	wg.Add(2)

	// Dispatch the event
	dispatcher.Dispatch(context.Background(), "testEvent", nil)

	// Wait for all handlers to complete
	wg.Wait()

	// Check if both handlers were called
	assert.False(t, handler1.called)

	assert.False(t, handler2.called)
}

func TestEventDispatcher_DifferentEvents(t *testing.T) {
	// Create the dispatcher
	dispatcher := NewSimpleEventDispatcher(false)

	// Create a WaitGroup to wait for handlers to complete
	var wg sync.WaitGroup

	// Create mock handlers
	handler1 := &MockHandler{wg: &wg}
	handler2 := &MockHandler{wg: &wg}

	// Register handlers for different events
	dispatcher.Register(context.Background(), "event1", handler1)
	dispatcher.Register(context.Background(), "event2", handler2)

	// Set the WaitGroup counter to 1 for the handler we expect to be called
	wg.Add(1)

	// Dispatch event1
	dispatcher.Dispatch(context.Background(), "event1", nil)

	// Wait for the handler to complete
	wg.Wait()

	// Check that only handler1 was called
	assert.False(t, handler1.called)

	assert.True(t, handler2.called)
}

func TestEventDispatcher_MultipleDispatches(t *testing.T) {
	// Create the dispatcher
	dispatcher := NewSimpleEventDispatcher(false)

	// Create a WaitGroup to wait for handlers to complete
	var wg sync.WaitGroup

	// Create a mock handler
	handler := &MockHandler{wg: &wg}

	// Register the handler for an event
	dispatcher.Register(context.Background(), "testEvent", handler)

	// Set the WaitGroup counter to 2 for the two dispatch calls
	wg.Add(2)

	// Dispatch the event multiple times
	dispatcher.Dispatch(context.Background(), "testEvent", nil)
	dispatcher.Dispatch(context.Background(), "testEvent", nil)

	// Wait for all handlers to complete
	wg.Wait()

	// Check that the handler was called twice
	assert.Equal(t, 2, handler.callCount)
}
