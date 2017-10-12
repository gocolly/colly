package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/asciimoo/colly"
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector()

	// Visit only domains: coursera.org, www.coursera.org
	c.AllowedDomains = []string{"coursera.org", "www.coursera.org"}
	courses := make(map[string]string)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// If attribute class of a is this long string return from callback
		// As this a is irrelevant
		if e.Attr("class") == "Button_1qxkboh-o_O-primary_cv02ee-o_O-md_28awn8-o_O-primaryLink_109aggg" {
			return
		}
		link := e.Attr("href")
		// If link start with browse or includes either signup or login return from callback
		if !strings.HasPrefix(link, "/browse") || strings.Index(link, "=signup") > -1 || strings.Index(link, "=login") > -1 {
			return
		}
		// start scaping the page under the link found
		e.Request.Visit(link)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	// On every a HTML element which has name attribute call callback
	c.OnHTML(`a[name]`, func(e *colly.HTMLElement) {
		// Add to courses map where key is the absolute URL and the
		// values is the name of the course
		courses[e.Request.AbsoluteURL(e.Attr("href"))] = e.Text
	})

	// Start scraping on http://coursera.com/browse
	c.Visit("https://coursera.org/browse")

	// After scraping is finished print the results
	for url, title := range courses {
		fmt.Println(url, "-", title)
	}
}
