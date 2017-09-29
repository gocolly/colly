# Colly

Scraping Framework for Gophers


Example:
```go
import (
	"fmt"
	"github.com/asciimoo/colly"
)

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


## Features

 * Clean API
 * Cookies and session handling
 * Sync/async/parallel scraping
 * Fast (>1k request/sec on a single core)


## Bugs

Bugs or suggestions? Visit the [issue tracker](https://github.com/asciimoo/colly/issues) or join `#colly` on freenode
