package extensions

import (
	"context"
	"github.com/go-colly/colly"
)

// Referrer sets valid Referrer HTTP header to requests.
// Warning: this extension works only if you use Request.Visit
// from callbacks instead of Collector.Visit.
func Referrer(c *colly.Collector) {
	c.OnResponse(func(_ context.Context, r *colly.Response) {
		dctx := colly.ContextDataContext(r.Ctx)
		dctx.Put("_referrer", r.Request.URL.String())
	})
	c.OnRequest(func(_ context.Context, r *colly.Request) {
		dctx := colly.ContextDataContext(r.Ctx)
		if ref := dctx.Get("_referrer"); ref != "" {
			r.Headers.Set("Referrer", ref)
		}
	})
}
