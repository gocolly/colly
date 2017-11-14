package debug

import (
	"log"
	"sync/atomic"
	"time"
)

type Event struct {
	Type   string
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

// Init implements the Init() function of the Debugger interface
func (l *LogDebugger) Init() error {
	l.counter = 0
	l.start = time.Now()
	return nil
}

// Event handles Collector events and prints them to STDERR
func (l *LogDebugger) Event(e *Event) {
	i := atomic.AddInt32(&l.counter, 1)
	log.Printf("[%6d] [%s] %q (%s)\n", i, e.Type, e.Values, time.Since(l.start))
}
