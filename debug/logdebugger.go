package debug

import (
	"log"
	"sync/atomic"
	"time"
)

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
