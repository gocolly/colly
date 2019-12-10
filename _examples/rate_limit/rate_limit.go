package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
)

func main() {
	url := "https://httpbin.org/delay/2"

	// Instantiate default collector
	c := colly.NewCollector(
		// Turn on asynchronous requests
		colly.Async(true),
		// Attach a debugger to the collector
		colly.Debugger(&debug.LogDebugger{}),
	)

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 2,
		//Delay:      5 * time.Second,
	})

	// Start scraping in five threads on https://httpbin.org/delay/2
	for i := 0; i < 5; i++ {
		c.Visit(fmt.Sprintf("%s?n=%d", url, i))
	}
	// Wait until threads are finished
	c.Wait()
}
