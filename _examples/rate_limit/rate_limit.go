package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	url := "https://httpbin.org/delay/2"

	// Instantiate default collector
	c := colly.NewCollector()

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 2,
		//Delay:      5 * time.Second,
	})

	// Before making a request print "Starting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Starting", r.URL, time.Now())
	})

	// After making a request print "Finished ..."
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL, time.Now())
	})

	// Start scraping in four threads on https://httpbin.org/delay/2
	for i := 0; i < 4; i++ {
		go c.Visit(fmt.Sprintf("%s?n=%d", url, i))
	}
	// Start scraping on https://httpbin.org/delay/2
	c.Visit(url)
	// Wait until threads are finished
	c.Wait()
}
