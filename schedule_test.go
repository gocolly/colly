package colly

import (
	"testing"
	"time"
)

func TestStartSchedules(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector()
	visited := make(chan bool, 1)
	c.OnRequest(func(r *Request) {
		visited <- true
	})

	// start schedule on-the-minute
	c.Schedule("@always", ts.URL)
	c.StartSchedules()

	sec := 65 - time.Now().Second()
	tick := time.NewTimer(time.Duration(sec) * time.Second)
L:
	for {
		select {
		case <-visited:
			break L
		case <-tick.C:
			t.Errorf("Schedule failed to start after %d sec", sec)
		default:
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func TestStartSchedulesWaitAndStop(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector()
	visited := make(chan bool, 1)
	c.OnRequest(func(r *Request) {
		visited <- true
	})

	go func() {
		sec := 65 - time.Now().Second()
		tick := time.NewTimer(time.Duration(sec) * time.Second)
	L:
		for {
			select {
			case <-visited:
				c.StopSchedules()
				break L
			case <-tick.C:
				t.Errorf("Schedule failed to start/stop after %d sec", sec)
			default:
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	// start schedule on-the-minute
	c.Schedule("@always", ts.URL)
	c.StartSchedulesWait()
}
