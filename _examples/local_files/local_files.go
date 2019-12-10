package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gocolly/colly/v2"
)

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	c := colly.NewCollector()
	c.WithTransport(t)

	pages := []string{}

	c.OnHTML("h1", func(e *colly.HTMLElement) {
		pages = append(pages, e.Text)
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		c.Visit("file://" + dir + "/html" + e.Attr("href"))
	})

	fmt.Println("file://" + dir + "/html/index.html")
	c.Visit("file://" + dir + "/html/index.html")
	c.Wait()
	for i, p := range pages {
		fmt.Printf("%d : %s\n", i, p)
	}
}
