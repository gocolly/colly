# Colly

Lightning Fast and Elegant Scraping Framework for Gophers

Colly provides a clean interface to write any kind of crawler/scraper/spider.

With Colly you can easily extract structured data from websites, which can be used for a wide range of applications, like data mining, data processing or archiving.

[![GoDoc](https://godoc.org/github.com/gocolly/colly?status.svg)](https://godoc.org/github.com/gocolly/colly)
[![build status](https://img.shields.io/travis/gocolly/colly/master.svg?style=flat-square)](https://travis-ci.org/gocolly/colly)
[![report card](https://img.shields.io/badge/report%20card-a%2B-ff3333.svg?style=flat-square)](http://goreportcard.com/report/gocolly/colly)
[![view examples](https://img.shields.io/badge/learn%20by-examples-0077b3.svg?style=flat-square)](https://github.com/gocolly/colly/tree/master/_examples)
[![test coverage](https://cover.run/go/github.com/gocolly/colly.svg)](https://cover.run/go/github.com/gocolly/colly)

## Features

 * Clean API
 * Fast (>1k request/sec on a single core)
 * Manages request delays and maximum concurrency per domain
 * Automatic cookie and session handling
 * Sync/async/parallel scraping
 * Caching
 * Automatic encoding of non-unicode responses
 * Robots.txt support


## Example

```go
func main() {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://en.wikipedia.org/")
}
```

See [examples folder](https://github.com/gocolly/colly/tree/master/_examples) for more detailed examples.


## Installation

```
go get -u github.com/gocolly/colly/...
```


## Bugs

Bugs or suggestions? Visit the [issue tracker](https://github.com/gocolly/colly/issues) or join `#colly` on freenode


## Other Projects Using Colly

Below is a list of public, open source projects that use Colly:

 * [greenpeace/check-my-pages](https://github.com/greenpeace/check-my-pages) Scraping script to test the Spanish Greenpeace web archive

If you are using Colly in a project please send a pull request to add it to the list.
