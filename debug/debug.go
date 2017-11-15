package debug

// Event represents an action inside a collector
type Event struct {
	// Type is the type of the event
	Type string
	// RequestId identifies the HTTP request of the Event
	RequestId int32
	// CollectorId identifies the collector of the Event
	CollectorId int32
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
