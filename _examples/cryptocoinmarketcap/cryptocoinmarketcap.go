package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

func main() {
	fName := "cryptocoinmarketcap.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Name", "Symbol", "Market Cap (USD)", "Price (USD)", "Circulating Supply (USD)", "Volume (24h)", "Change (1h)", "Change (24h)", "Change (7d)"})

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		writer.Write([]string{
			e.ChildText(".cmc-table__column-name"),
			e.ChildText(".cmc-table__cell--sort-by__symbol"),
			e.ChildText(".cmc-table__cell--sort-by__market-cap"),
			e.ChildText(".cmc-table__cell--sort-by__price"),
			e.ChildText(".cmc-table__cell--sort-by__circulating-supply"),
			e.ChildText(".cmc-table__cell--sort-by__volume-24-h"),
			e.ChildText(".cmc-table__cell--sort-by__percent-change-1-h"),
			e.ChildText(".cmc-table__cell--sort-by__percent-change-24-h"),
			e.ChildText(".cmc-table__cell--sort-by__percent-change-7-d"),
		})
	})

	c.Visit("https://coinmarketcap.com/all/views/all/")

	log.Printf("Scraping finished, check file %q for results\n", fName)
}
