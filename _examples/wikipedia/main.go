package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly/v2"
	"os"
	"strings"
	"time"
)

const (
	wikipediaURL    = "https://en.wikipedia.org"
	wikiArticlePath = "/wiki/"
)

var (
	start       string
	destination string
)

func init() {
	flag.StringVar(&start, "start", "Coldplay", "title of start article")
	flag.StringVar(&destination, "destination", "Jack_Sparrow", "title of destination article")
	flag.Parse()
}

func main() {
	t0 := time.Now()

	c := colly.NewCollector(colly.AllowedDomains("en.wikipedia.org"),
		colly.MaxDepth(6),
		colly.Async(true))

	c.Limit(&colly.LimitRule{DomainGlob: "en.wikipedia.org", Parallelism: 20})

	start := wikipediaURL + wikiArticlePath + start
	destination := wikipediaURL + wikiArticlePath + destination

	// Visit links within article text where the relative path starts with /wiki/ to exclude miscellaneous links
	c.OnHTML("p > a[href]", func(e *colly.HTMLElement) {
		if strings.HasPrefix(e.Attr("href"), wikiArticlePath) {
			e.Request.Visit(wikipediaURL + e.Attr("href"))
		}
	})

	c.OnRequest(func(request *colly.Request) {
		fmt.Println("visiting", request.URL)
		if request.URL.String() == destination {
			fmt.Printf("Found a path with %d degrees of separation\n", request.Depth-1)
			fmt.Printf("Time: %s\n", time.Now().Sub(t0).String())
			os.Exit(0)
		}
	})

	c.Visit(start)

	c.Wait()
}
