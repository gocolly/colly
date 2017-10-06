package main

import (
	"fmt"
	"time"

	"github.com/asciimoo/colly"
)

func main() {
	url := "https://httpbin.org/delay/2"

	c := colly.NewCollector()

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 2,
		//Delay:      5 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Starting", r.URL, time.Now())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL, time.Now())
	})

	for i := 0; i < 4; i++ {
		go c.Visit(fmt.Sprintf("%s?n=%d", url, i))
	}
	c.Visit(url)
	c.Wait()
}
