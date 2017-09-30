package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/asciimoo/colly"
)

func main() {
	c := colly.NewCollector()

	c.AllowedDomains = []string{"coursera.org", "www.coursera.org"}
	courses := make(map[string]string)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if e.Attr("class") == "Button_1qxkboh-o_O-primary_cv02ee-o_O-md_28awn8-o_O-primaryLink_109aggg" {
			return
		}
		link := e.Attr("href")
		if !strings.HasPrefix(link, "/browse") || strings.Index(link, "=signup") > -1 || strings.Index(link, "=login") > -1 {
			return
		}
		e.Request.Visit(link)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	c.OnHTML(`a[name]`, func(e *colly.HTMLElement) {
		courses[e.Request.AbsoluteURL(e.Attr("href"))] = e.Text
	})

	c.Visit("https://coursera.org/browse")

	for url, title := range courses {
		fmt.Println(url, "-", title)
	}
}
