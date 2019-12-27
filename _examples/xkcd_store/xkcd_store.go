package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

func main() {
	fName := "xkcd_store_items.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	// Write CSV header
	writer.Write([]string{"Name", "Price", "URL", "Image URL"})

	// Instantiate default collector
	c := colly.NewCollector(
		// Allow requests only to store.xkcd.com
		colly.AllowedDomains("store.xkcd.com"),
	)

	// Extract product details
	c.OnHTML(".product-grid-item", func(e *colly.HTMLElement) {
		writer.Write([]string{
			e.ChildAttr("a", "title"),
			e.ChildText("span"),
			e.Request.AbsoluteURL(e.ChildAttr("a", "href")),
			"https" + e.ChildAttr("img", "src"),
		})
	})

	// Find and visit next page links
	c.OnHTML(`.next a[href]`, func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.Visit("https://store.xkcd.com/collections/everything")

	log.Printf("Scraping finished, check file %q for results\n", fName)

	// Display collector's statistics
	log.Println(c)
}
