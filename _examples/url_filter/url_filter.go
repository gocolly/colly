package main

import (
	"fmt"
	"regexp"

	"github.com/gocolly/colly/v2"
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only root url and urls which start with "e" or "h" on httpbin.org
		colly.URLFilters(
			regexp.MustCompile("http://httpbin\\.org/(|e.+)$"),
			regexp.MustCompile("http://httpbin\\.org/h.+"),
		),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are matched by  any of the URLFilter regexps
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on http://httpbin.org
	c.Visit("http://httpbin.org/")
}
