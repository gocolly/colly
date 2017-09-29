package main

import (
	"fmt"
	"github.com/asciimoo/colly"
)

func main() {
	c := colly.NewCollector()

	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Println(link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.Visit("https://en.wikipedia.org/")
}
