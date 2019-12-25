package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	// Create a collector
	c := colly.NewCollector()

	// Set HTML callback
	// Won't be called if error occurs
	c.OnHTML("*", func(e *colly.HTMLElement) {
		fmt.Println(e)
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start scraping
	c.Visit("https://definitely-not-a.website/")
}
