package debug

import (
	"log"
	"sync/atomic"
	"time"
)

// Event represents an action inside a collector
type Event struct {
	// Type is the type of the event
	Type string
	// RequestId identifies the HTTP request of the Event
	RequestId int32
	// Values contains the event's key-value pairs. Different type of events
	// can return different key-value pairs
	Values map[string]string
}

// Debugger is an interface for different type of debugging backends
type Debugger interface {
	// Init initializes the backend
	Init() error
	// Event receives a new collector event.
	Event(e *Event)
}

// LogDebugger is the simplest debugger which prints log messages to the STDERR
type LogDebugger struct {
	counter int32
	start   time.Time
}

// Init initializes the LogDebugger
func (l *LogDebugger) Init() error {
	l.counter = 0
	l.start = time.Now()
	return nil
}

// Event receives Collector events and prints them to STDERR
func (l *LogDebugger) Event(e *Event) {
	i := atomic.AddInt32(&l.counter, 1)
	log.Printf("[%06d] [%6d - %s] %q (%s)\n", i, e.RequestId, e.Type, e.Values, time.Since(l.start))
}
