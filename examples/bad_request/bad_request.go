package main

import (
	"fmt"

	"github.com/asciimoo/colly"
)

func main() {
	c := colly.NewCollector()

	c.AllowedDomains = []string{"definitely-note-a-website.not"}

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(req *colly.Request, resp *colly.Response, err error) {
		fmt.Println("Request:", req, "\nfailed with response:", resp, "\nand error:", err)
	})

	c.Visit("https://definitely-note-a-website.not/")
}
