package colly

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const testDelay = 200 * time.Millisecond

func newTraceTestServer(delay time.Duration) *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(200)
	})
	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(500)
	})

	return httptest.NewServer(mux)
}

func TestTraceWithNoDelay(t *testing.T) {
	ts := newTraceTestServer(0)
	defer ts.Close()

	client := ts.Client()
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Errorf("Failed to construct request %v", err)
	}
	trace := &HTTPTrace{}
	req = trace.WithTrace(req)

	if _, err = client.Do(req); err != nil {
		t.Errorf("Failed to make request %v", err)
	}

	if trace.ConnectDuration > testDelay {
		t.Errorf("trace ConnectDuration should be (almost) 0, got %v", trace.ConnectDuration)
	}
	if trace.FirstByteDuration > testDelay {
		t.Errorf("trace FirstByteDuration should be (almost) 0, got %v", trace.FirstByteDuration)
	}
}

func TestTraceWithDelay(t *testing.T) {
	ts := newTraceTestServer(testDelay)
	defer ts.Close()

	client := ts.Client()
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Errorf("Failed to construct request %v", err)
	}
	trace := &HTTPTrace{}
	req = trace.WithTrace(req)

	if _, err = client.Do(req); err != nil {
		t.Errorf("Failed to make request %v", err)
	}

	if trace.ConnectDuration > testDelay {
		t.Errorf("trace ConnectDuration should be (almost) 0, got %v", trace.ConnectDuration)
	}
	if trace.FirstByteDuration < testDelay {
		t.Errorf("trace FirstByteDuration should be at least 200ms, got %v", trace.FirstByteDuration)
	}
}
