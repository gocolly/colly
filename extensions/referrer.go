package extensions

import (
	"github.com/gocolly/colly"
)

// Referrer sets valid Referrer HTTP header to requests.
// Warning: this extension works only if you use Request.Visit
// from callbacks instead of Collector.Visit.
func Referrer(c *colly.Collector) {
	c.OnResponse(func(r *colly.Response) {
		r.Ctx.Put("_referrer", r.Request.URL.String())
	})
	c.OnRequest(func(r *colly.Request) {
		if ref := r.Ctx.Get("_referrer"); ref != "" {
			r.Headers.Set("Referrer", ref)
		}
	})
}
