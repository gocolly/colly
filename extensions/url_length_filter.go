package extensions

import (
	"context"
	"github.com/gocolly/colly"
)

// URLLengthFilter filters out requests with URLs longer than URLLengthLimit
func URLLengthFilter(c *colly.Collector, URLLengthLimit int) {
	c.OnRequest(func(ctx context.Context, r *colly.Request) {
		if len(r.URL.String()) > URLLengthLimit {
			r.Abort()
		}
	})
}
