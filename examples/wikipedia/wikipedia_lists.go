package main

import (
	"fmt"
	"github.com/asciimoo/colly"
)

func main() {
	c := colly.NewCollector()

	//c.Cache = false
	//c.UserAgent = "myUserAgent"

	c.OnHTML("a", func(e *colly.HTMLElement) {
		fmt.Println(e.Attr("href"))
		c.Visit(e.Attr("href"))
	})

	//c.OnResponse(func(r *colly.Response) {
	//	r.Ctx.Get("x")
	//	fmt.Println(r)
	//})

	//c.OnRequest(func(r *colly.Request) {
	//	r.Ctx.Put("x", "y")
	//	if r.Path == "/" {
	//		r.UserAgent = "agent2"
	//	}
	//	fmt.Println(r)
	//})

	c.Visit("https://en.wikipedia.org/wiki/Category:Lists")
}
