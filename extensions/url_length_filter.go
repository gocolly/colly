package extensions

import (
	"context"
	"github.com/go-colly/colly"
)

// URLLengthFilter filters out requests with URLs longer than URLLengthLimit
func URLLengthFilter(c *colly.Collector, URLLengthLimit int) {
	c.OnRequest(func(_ context.Context, r *colly.Request) {
		if len(r.URL.String()) > URLLengthLimit {
			r.Abort()
		}
	})
}
