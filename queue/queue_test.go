package queue

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gocolly/colly/v2"
)

func TestQueue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(serverHandler))
	defer server.Close()

	rng := rand.New(rand.NewSource(12387123712321232))
	var (
		items    uint32
		requests uint32
		success  uint32
		failure  uint32
	)
	storage := &InMemoryQueueStorage{MaxSize: 100000}
	q, err := New(10, storage)
	if err != nil {
		panic(err)
	}
	put := func() {
		t := time.Duration(rng.Intn(50)) * time.Microsecond
		url := server.URL + "/delay?t=" + t.String()
		atomic.AddUint32(&items, 1)
		q.AddURL(url)
	}
	for i := 0; i < 3000; i++ {
		put()
		storage.AddRequest([]byte("error request"))
	}
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	c.OnRequest(func(req *colly.Request) {
		atomic.AddUint32(&requests, 1)
	})
	c.OnResponse(func(resp *colly.Response) {
		if resp.StatusCode == http.StatusOK {
			atomic.AddUint32(&success, 1)
		} else {
			atomic.AddUint32(&failure, 1)
		}
		toss := rng.Intn(2) == 0
		if toss {
			put()
		}
	})
	c.OnError(func(resp *colly.Response, err error) {
		atomic.AddUint32(&failure, 1)
	})
	err = q.Run(c)
	if err != nil {
		t.Fatalf("Queue.Run() return an error: %v", err)
	}
	if items != requests || success+failure != requests || failure > 0 {
		t.Fatalf("wrong Queue implementation: "+
			"items = %d, requests = %d, success = %d, failure = %d",
			items, requests, success, failure)
	}
}

func serverHandler(w http.ResponseWriter, req *http.Request) {
	if !serverRoute(w, req) {
		shutdown(w)
	}
}

func serverRoute(w http.ResponseWriter, req *http.Request) bool {
	if req.URL.Path == "/delay" {
		return serveDelay(w, req) == nil
	}
	return false
}

func serveDelay(w http.ResponseWriter, req *http.Request) error {
	q := req.URL.Query()
	t, err := time.ParseDuration(q.Get("t"))
	if err != nil {
		return err
	}
	time.Sleep(t)
	w.WriteHeader(http.StatusOK)
	return nil
}

func shutdown(w http.ResponseWriter) {
	taker, ok := w.(http.Hijacker)
	if !ok {
		return
	}
	raw, _, err := taker.Hijack()
	if err != nil {
		return
	}
	raw.Close()
}
