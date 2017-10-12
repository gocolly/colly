# Colly

Lightning Fast and Elegant Scraping Framework for Gophers

Colly provides a clean interface to write any kind of crawler/scraper/spider.

With Colly you can easily extract structured data from websites, which can be used for a wide range of applications, like data mining, data processing or archiving.

[![GoDoc](https://godoc.org/github.com/asciimoo/colly?status.svg)](https://godoc.org/github.com/asciimoo/colly)
[![build status](https://img.shields.io/travis/asciimoo/colly/master.svg?style=flat-square)](https://travis-ci.org/asciimoo/colly)
[![report card](https://img.shields.io/badge/report%20card-a%2B-ff3333.svg?style=flat-square)](http://goreportcard.com/report/asciimoo/colly)
[![view examples](https://img.shields.io/badge/learn%20by-examples-0077b3.svg?style=flat-square)](https://github.com/asciimoo/colly/tree/master/examples)

## Features

 * Clean API
 * Fast (>1k request/sec on a single core)
 * Manages request delays and maximum concurrency per domain
 * Automatic cookie and session handling
 * Sync/async/parallel scraping
 * Caching
 * Automatic encoding of non-unicode responses


## Example

```go
func main() {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Println(link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.Visit("https://en.wikipedia.org/")
}
```

See [examples folder](https://github.com/asciimoo/colly/tree/master/examples) for more detailed examples.


## Bugs

Bugs or suggestions? Visit the [issue tracker](https://github.com/asciimoo/colly/issues) or join `#colly` on freenode
