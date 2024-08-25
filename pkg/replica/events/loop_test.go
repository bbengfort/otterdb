package events_test

import (
	"sync"
	"testing"

	. "github.com/bbengfort/otterdb/pkg/replica/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Generates a fixed number of events for testing event loops.
func Generator(wg *sync.WaitGroup, events chan<- Event) {
	defer wg.Done()
	defer close(events)

	// Send a burst of 128 write ahead events
	for i := 0; i < 128; i++ {
		events <- &MockWriteAhead{}
	}

	// Send a heartbeat event
	events <- &MockHeartbeat{}

	// Send a burst of 1448 write ahead events
	for i := 0; i < 1448; i++ {
		events <- &MockWriteAhead{}
	}

	// Alternate between write ahead and heartbeat events
	for i := 0; i < 128; i++ {
		if i%2 == 0 {
			events <- &MockHeartbeat{}
		} else {
			events <- &MockWriteAhead{}
		}
	}

	// Send another burst of write ahead events
	for i := 0; i < 128; i++ {
		events <- &MockWriteAhead{}
	}

	// Send a heartbeat event
	events <- &MockHeartbeat{}
}

func TestLoop(t *testing.T) {
	wg := new(sync.WaitGroup)
	events := make(chan Event, BufferSize)
	mock := &MockHandler{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := Loop(events, mock)
		assert.NoError(t, err)
	}()

	wg.Add(1)
	go Generator(wg, events)

	wg.Wait()
	require.Equal(t, 1768, mock.events[WriteAhead])
	require.Equal(t, 0, mock.events[AggregatedWriteAhead])
	require.Equal(t, 66, mock.events[HeartbeatTimeout])
}

func TestAggregatingLoop(t *testing.T) {
	wg := new(sync.WaitGroup)
	events := make(chan Event, BufferSize)
	mock := &MockHandler{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := AggregatingLoop(events, mock)
		assert.NoError(t, err)
	}()

	wg.Add(1)
	go Generator(wg, events)

	wg.Wait()
	require.Equal(t, 63, mock.events[WriteAhead])
	require.GreaterOrEqual(t, mock.events[AggregatedWriteAhead], 5)
	require.LessOrEqual(t, mock.events[AggregatedWriteAhead], 7)
	require.Equal(t, 66, mock.events[HeartbeatTimeout])
}

type MockHandler struct {
	events map[EventType]int
}

func (d *MockHandler) Handle(e Event) error {
	if d.events == nil {
		d.events = make(map[EventType]int)
	}

	d.events[e.Event()]++
	return nil
}

type MockWriteAhead struct{}

func (m *MockWriteAhead) Event() EventType {
	return WriteAhead
}

type MockHeartbeat struct{}

func (m *MockHeartbeat) Event() EventType {
	return HeartbeatTimeout
}
