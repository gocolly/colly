package colly

import (
	"context"

	"github.com/adhocore/gronx/pkg/tasker"
)

type schedule struct {
	cron string
	url  string
	err  chan error
}

var (
	scheduleCtx    context.Context
	scheduleCancel context.CancelFunc
)

// Schedule adds a new item to the schedule list
// eg. Visit example.com/home every 5 hours
// c.Schedule("* */5 * * *", "https://example.com/home")
//
// Cron expressions are parsed via package
//	 https://github.com/adhocore/gronx#cron-expression
func (c *Collector) Schedule(expr, u string) chan (error) {
	ch := make(chan error)
	c.schedules = append(c.schedules, schedule{cron: expr, url: u, err: ch})
	return ch
}

// StartSchedules will begin each listed schedules.
// Collector.Context can be used to cancel the entire
// list of schedules.
func (c *Collector) StartSchedules() {
	if len(c.schedules) == 0 {
		return
	}

	var ctx context.Context
	if c.Context != nil {
		ctx = c.Context
	} else {
		ctx = context.Background()
	}

	scheduleCtx, scheduleCancel = context.WithCancel(ctx)
	taskr := tasker.New(tasker.Option{}).WithContext(scheduleCtx)
	go func() {
		for _, s := range c.schedules {
			taskr.Task(s.cron, func(ctx context.Context) (int, error) {
				s.err <- c.Visit(s.url)
				return 0, nil
			})
		}
		taskr.Run()
	}()
}

// StartSchedulesWait will begin each listed schedules
// and will block until context cancelled
//
// Does not block for all currently running schedules to finish
func (c *Collector) StartSchedulesWait() {
	c.StartSchedules()
	<-scheduleCtx.Done()
}

// StopSchedules will cancel the underlying context and returns
// blocking operations like StartSchedulesWait
//
// An alternative would be cancelling Collector.Context
func (c *Collector) StopSchedules() {
	if scheduleCancel != nil {
		scheduleCancel()
	}
}
