package main

import (
	"context"
	"log"

	"github.com/gocolly/colly"
)

func main() {
	// create a new collector
	c := colly.NewCollector()

	// authenticate
	err := c.Post(nil, "http://example.com/login", map[string]string{"username": "admin", "password": "admin"})
	if err != nil {
		log.Fatal(err)
	}

	// attach callbacks after login
	c.OnResponse(func(_ context.Context, r *colly.Response) {
		log.Println("response received", r.StatusCode)
	})

	// start scraping
	c.Visit("https://example.com/")
}
