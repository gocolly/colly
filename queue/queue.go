package queue

import (
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/gocolly/colly"
)

const stop = true

// Storage is the interface of the queue's storage backend
type Storage interface {
	// Init initializes the storage
	Init() error
	// AddRequest adds a serialized request to the queue
	AddRequest([]byte) error
	// GetRequest pops the next request from the queue
	// or returns error if the queue is empty
	GetRequest() ([]byte, error)
	// QueueSize returns with the size of the queue
	QueueSize() (int, error)
}

// Queue is a request queue which uses a Collector to consume
// requests in multiple threads
type Queue struct {
	// Threads defines the number of consumer threads
	Threads           int
	storage           Storage
	activeThreadCount int32
	threadChans       []chan bool
	lock              *sync.Mutex
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
	Request []byte
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
		Threads:     threads,
		storage:     s,
		lock:        &sync.Mutex{},
		threadChans: make([]chan bool, 0, threads),
	}, nil
}

// IsEmpty returns true if the queue is empty
func (q *Queue) IsEmpty() bool {
	s, _ := q.Size()
	return s == 0
}

// AddURL adds a new URL to the queue
func (q *Queue) AddURL(URL string) error {
	u, err := url.Parse(URL)
	if err != nil {
		return err
	}
	r := &colly.Request{
		URL:    u,
		Method: "GET",
	}
	d, err := r.Marshal()
	if err != nil {
		return err
	}
	return q.storage.AddRequest(d)
}

// AddRequest adds a new Request to the queue
func (q *Queue) AddRequest(r *colly.Request) error {
	d, err := r.Marshal()
	if err != nil {
		return err
	}
	if err := q.storage.AddRequest(d); err != nil {
		return err
	}
	q.lock.Lock()
	for _, c := range q.threadChans {
		c <- !stop
	}
	q.threadChans = make([]chan bool, 0, q.Threads)
	q.lock.Unlock()
	return nil
}

// Size returns the size of the queue
func (q *Queue) Size() (int, error) {
	return q.storage.QueueSize()
}

// Run starts consumer threads and calls the Collector
// to perform requests. Run blocks while the queue has active requests
func (q *Queue) Run(c *colly.Collector) error {
	wg := &sync.WaitGroup{}
	for i := 0; i < q.Threads; i++ {
		wg.Add(1)
		go func(c *colly.Collector, wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				if q.IsEmpty() {
					if q.activeThreadCount == 0 {
						break
					}
					ch := make(chan bool)
					q.lock.Lock()
					q.threadChans = append(q.threadChans, ch)
					q.lock.Unlock()
					action := <-ch
					if action == stop && q.IsEmpty() {
						break
					}
				}
				atomic.AddInt32(&q.activeThreadCount, 1)
				rb, err := q.storage.GetRequest()
				if err != nil || rb == nil {
					q.finish()
					continue
				}
				r, err := c.UnmarshalRequest(rb)
				if err != nil || r == nil {
					q.finish()
					continue
				}
				r.Do()
				q.finish()
			}
		}(c, wg)
	}
	wg.Wait()
	return nil
}

func (q *Queue) finish() {
	atomic.AddInt32(&q.activeThreadCount, -1)
	q.lock.Lock()
	for _, c := range q.threadChans {
		c <- stop
	}
	q.threadChans = make([]chan bool, 0, q.Threads)
	q.lock.Unlock()
}

// Init implements Storage.Init() function
func (q *InMemoryQueueStorage) Init() error {
	q.lock = &sync.RWMutex{}
	return nil
}

// AddRequest implements Storage.AddRequest() function
func (q *InMemoryQueueStorage) AddRequest(r []byte) error {
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
func (q *InMemoryQueueStorage) GetRequest() ([]byte, error) {
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
