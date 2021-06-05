package colly

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestHTTPBackendDoCancelation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprint(rw, "OK")
	}))
	defer ts.Close()

	checkHeadersFunc := func(req *http.Request, statusCode int, header http.Header) bool { return true }

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	backend := &httpBackend{}
	jar, _ := cookiejar.New(nil)
	backend.Init(jar)
	limit := &LimitRule{
		DomainRegexp: ".*",
		Parallelism:  1,
		Delay:        5 * time.Millisecond,
	}
	backend.Limit(limit)

	rand.Seed(time.Now().Unix())

	// rand up to 10 to not extend the test duration too much
	n := rand.Intn(10)
	c := n + rand.Intn(10)

	var wg sync.WaitGroup
	wg.Add(c)

	errs := make(chan error)

	begin := time.Now()
	for i := 0; i < c; i++ {
		go func(i int) {
			defer wg.Done()

			req, _ := http.NewRequestWithContext(ctx, "GET", ts.URL+"/"+strconv.Itoa(i), nil)
			_, err := backend.Do(req, 0, checkHeadersFunc)
			errs <- err
		}(i)
	}

	var d time.Duration
	go func() {
		wg.Wait()
		d = time.Since(begin) // captures the duration of all calls
		close(errs)
	}()

	i := 0
	for err := range errs {
		i++
		if i == n {
			cancel()
		}

		if i <= n {
			if err != nil {
				t.Errorf("no error was expected for the first %d responses; error: %q", n, err)
			}
		} else {
			if !strings.Contains(err.Error(), "context canceled") {
				t.Error("call to Do should return with error from terminated context")
			}
		}
	}

	// the expectation is n+1 because:
	//     n: that time should have already passed, cancel is done after the second response
	//     1: the third request is cancelled just after starting forces delay
	if d < time.Duration(n+1)*limit.Delay {
		t.Error("duration is bellow the expected")
	}

	// n+1+1 is the max limit because after the minimum delay the other calls should finish
	// immediately since the function return before defer was set
	if d > time.Duration(n+1+1)*limit.Delay {
		t.Error("duration is above the expected")
	}
}
