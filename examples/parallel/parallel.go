package main

import (
	"fmt"
	"github.com/asciimoo/colly"
)

func main() {
	c := colly.NewCollector()

	c.MaxDepth = 2

	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Println(link)
		go e.Request.Visit(link)
	})

	c.Visit("https://en.wikipedia.org/")
	c.Wait()
}
