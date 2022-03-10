package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gocolly/colly"
)

func main() {

	file_name := "data.csv"
	file, err := os.Create(file_name)

	if err != nil {
		log.Fatalf("Error Occured :%q while creatig file %s", err, file_name)
	}

	defer file.Close()
	writer := csv.NewWriter(file)

	c := colly.NewCollector(
		colly.AllowedDomains("internshala.com"),
	)

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL)
	// })

	var total_pages = 1

	c.OnHTML("#pagination", func(h *colly.HTMLElement) {
		total_pages, err = strconv.Atoi(h.ChildText("#total_pages"))
		if err != nil {
			fmt.Printf("Error Occured here is %s", err)
		}
	})

	c.OnHTML(".internship_meta", func(p *colly.HTMLElement) {

		writer.Write([]string{
			p.ChildText("a"),
			p.ChildText("span"),
			p.ChildText("div.item_body"),
		})
	})

	for i := 0; i < total_pages; i++ {
		fmt.Printf("Scrapping Page %d\n", i)
		c.Visit("https://internshala.com/internships/page-" + strconv.Itoa(i))
	}

	writer.Flush()

	log.Println("Scrapping Completed ****** 100% ******")
	readcsv()
}
