package main

import (
	"context"
	"fmt"

	"github.com/go-colly/colly"
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector()

	// Before making a request put the URL with
	// the key of "url" into the context of the request
	c.OnRequest(func(_ context.Context, r *colly.Request) {
		dctx := colly.ContextDataContext(r.Ctx)
		dctx.Put("url", r.URL.String())
	})

	// After making a request get "url" from
	// the context of the request
	c.OnResponse(func(_ context.Context, r *colly.Response) {
		dctx := colly.ContextDataContext(r.Ctx)
		fmt.Println(dctx.Get("url"))
	})

	// Start scraping on https://en.wikipedia.org
	c.Visit("https://en.wikipedia.org/")
}
