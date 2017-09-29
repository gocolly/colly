# Colly

Scraping Framework for Gophers


[Documentation](https://godoc.org/github.com/asciimoo/colly)

## Features

 * Clean API
 * Cookies and session handling
 * Sync/async/parallel scraping
 * Fast (>1k request/sec on a single core)


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
