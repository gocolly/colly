package extensions

import (
	"github.com/gocolly/colly/v2"
)

// Referer sets valid Referer HTTP header to requests.
// Warning: this extension works only if you use Request.Visit
// from callbacks instead of Collector.Visit.
func Referer(c *colly.Collector) {
	c.OnResponse(func(r *colly.Response) error {
		r.Ctx.Put("_referer", r.Request.URL.String())
		return nil
	})
	c.OnRequest(func(r *colly.Request) {
		if ref := r.Ctx.Get("_referer"); ref != "" {
			r.Headers.Set("Referer", ref)
		}
	})
}
