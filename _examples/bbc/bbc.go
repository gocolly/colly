package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36"),
	)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})
	c.OnHTML("h3", func(e *colly.HTMLElement) {
		title := e.Text
		fmt.Println("title", title)
	})


	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("error:", err)
		fmt.Println("status:", r.StatusCode)
		fmt.Println("body:", string(r.Body))
	})
	c.Visit("https://www.bbc.com/")

}
