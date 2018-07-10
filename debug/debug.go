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

// Event represents an action inside a collector
type Event struct {
	// Type is the type of the event
	Type string
	// RequestID identifies the HTTP request of the Event
	RequestID uint32
	// CollectorID identifies the collector of the Event
	CollectorID uint32
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
