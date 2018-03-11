// Copyright 2018 Adam Tauber
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package debug

import (
	"io"
	"log"
	"os"
	"sync/atomic"
	"time"
)

// LogDebugger is the simplest debugger which prints log messages to the STDERR
type LogDebugger struct {
	// Output is the log destination, anything can be used which implements them
	// io.Writer interface. Leave it blank to use STDERR
	Output io.Writer
	// Prefix appears at the beginning of each generated log line
	Prefix string
	// Flag defines the logging properties.
	Flag    int
	logger  *log.Logger
	counter int32
	start   time.Time
}

// Init initializes the LogDebugger
func (l *LogDebugger) Init() error {
	l.counter = 0
	l.start = time.Now()
	if l.Output == nil {
		l.Output = os.Stderr
	}
	l.logger = log.New(l.Output, l.Prefix, l.Flag)
	return nil
}

// Event receives Collector events and prints them to STDERR
func (l *LogDebugger) Event(e *Event) {
	i := atomic.AddInt32(&l.counter, 1)
	l.logger.Printf("[%06d] %d [%6d - %s] %q (%s)\n", i, e.CollectorID, e.RequestID, e.Type, e.Values, time.Since(l.start))
}
