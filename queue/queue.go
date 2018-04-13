package queue

import (
	"io"
	"sync"

	"github.com/gocolly/colly"
)

// Request is the container of HTTP request attributes
type Request struct {
	// URL of the request
	URL string
	// Method of the request
	Method string
	// Body is the optional request body for POST/PUT requests
	Body io.Reader
}

// Storage is the interface of the queue's storage backend
type Storage interface {
	// Init initializes the storage
	Init() error
	// AddRequest adds a request to the queue
	AddRequest(*Request) error
	// GetRequest pops the next request from the queue
	// or returns error if the queue is empty
	GetRequest() (*Request, error)
	// QueueSize returns with the size of the queue
	QueueSize() (int, error)
}

// Queue is a request queue which uses a Collector to consume
// requests in multiple threads
type Queue struct {
	// Threads defines the number of consumer threads
	Threads int
	storage Storage
}

// InMemoryQueueStorage is the default implementation of the Storage interface.
// InMemoryQueueStorage holds the request queue in memory.
type InMemoryQueueStorage struct {
	// MaxSize defines the capacity of the queue.
	// New requests are discarded if the queue size reaches MaxSize
	MaxSize int
	lock    *sync.RWMutex
	size    int
	first   *inMemoryQueueItem
	last    *inMemoryQueueItem
}

type inMemoryQueueItem struct {
	Request *Request
	Next    *inMemoryQueueItem
}

// New creates a new queue with a Storage specified in argument
// A standard InMemoryQueueStorage is used if Storage argument is nil.
func New(threads int, s Storage) (*Queue, error) {
	if s == nil {
		s = &InMemoryQueueStorage{MaxSize: 100000}
	}
	if err := s.Init(); err != nil {
		return nil, err
	}
	return &Queue{
		Threads: threads,
		storage: s,
	}, nil
}

// IsEmpty returns true if the queue is empty
func (q *Queue) IsEmpty() bool {
	s, _ := q.Size()
	return s == 0
}

// AddURL adds a new URL to the queue
func (q *Queue) AddURL(URL string) error {
	return q.storage.AddRequest(&Request{URL: URL, Method: "GET"})
}

// AddRequest adds a new Request to the queue
func (q *Queue) AddRequest(method, URL string, body io.Reader) error {
	return q.storage.AddRequest(&Request{URL: URL, Method: method, Body: body})
}

// Size returns the size of the queue
func (q *Queue) Size() (int, error) {
	return q.storage.QueueSize()
}

// Run starts consumer threads and calls the Collector
// to perform requests. Run blocks while the queue has items
func (q *Queue) Run(c *colly.Collector) error {
	wg := &sync.WaitGroup{}
	for i := 0; i < q.Threads; i++ {
		wg.Add(1)
		go func(c *colly.Collector, wg *sync.WaitGroup) {
			defer wg.Done()
			for !q.IsEmpty() {
				r, err := q.storage.GetRequest()
				if err != nil || r == nil {
					break
				}
				c.Request(r.Method, r.URL, r.Body, nil, nil)
			}
		}(c, wg)
	}
	wg.Wait()
	return nil
}

// Init implements Storage.Init() function
func (q *InMemoryQueueStorage) Init() error {
	q.lock = &sync.RWMutex{}
	return nil
}

// AddRequest implements Storage.AddRequest() function
func (q *InMemoryQueueStorage) AddRequest(r *Request) error {
	q.lock.Lock()
	defer q.lock.Unlock()
	// Discard URLs if size limit exceeded
	if q.MaxSize > 0 && q.size >= q.MaxSize {
		return nil
	}
	i := &inMemoryQueueItem{Request: r}
	if q.first == nil {
		q.first = i
	} else {
		q.last.Next = i
	}
	q.last = i
	q.size++
	return nil
}

// GetRequest implements Storage.GetRequest() function
func (q *InMemoryQueueStorage) GetRequest() (*Request, error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.size == 0 {
		return nil, nil
	}
	r := q.first.Request
	q.first = q.first.Next
	q.size--
	return r, nil
}

// QueueSize implements Storage.QueueSize() function
func (q *InMemoryQueueStorage) QueueSize() (int, error) {
	return q.size, nil
}
