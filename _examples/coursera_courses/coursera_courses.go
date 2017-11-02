package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/asciimoo/colly"
)

// Course stores information about a coursera course
type Course struct {
	Title       string
	Description string
	Creator     string
	Level       string
	URL         string
	Language    string
	Commitment  string
	HowToPass   string
	Rating      string
}

func main() {
	// Instantiate default collector
	c := colly.NewCollector()

	// Visit only domains: coursera.org, www.coursera.org
	c.AllowedDomains = []string{"coursera.org", "www.coursera.org"}

	// Cache responses to prevent multiple download of pages
	// even if the collector is restarted
	c.CacheDir = "./coursera_cache"

	// Create another collector to scrape course details
	detailCollector := c.Clone()

	courses := make([]Course, 0, 200)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// If attribute class of a is this long string return from callback
		// As this a is irrelevant
		if e.Attr("class") == "Button_1qxkboh-o_O-primary_cv02ee-o_O-md_28awn8-o_O-primaryLink_109aggg" {
			return
		}
		link := e.Attr("href")
		// If link start with browse or includes either signup or login return from callback
		if !strings.HasPrefix(link, "/browse") || strings.Index(link, "=signup") > -1 || strings.Index(link, "=login") > -1 {
			return
		}
		// start scaping the page under the link found
		e.Request.Visit(link)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	// On every a HTML element which has name attribute call callback
	c.OnHTML(`a[name]`, func(e *colly.HTMLElement) {
		// Activate detailCollector if the link contains "coursera.org/learn"
		courseURL := e.Request.AbsoluteURL(e.Attr("href"))
		if strings.Index(courseURL, "coursera.org/learn") != -1 {
			detailCollector.Visit(courseURL)
		}
	})

	// Extract details of the course
	detailCollector.OnHTML(`div[id=rendered-content]`, func(e *colly.HTMLElement) {
		log.Println("Course found", e.Request.URL)
		title := e.ChildText(".course-title")
		if title == "" {
			log.Println("No title found", e.Request.URL)
		}
		course := Course{
			Title:       title,
			URL:         e.Request.URL.String(),
			Description: e.ChildText("div.content"),
			Creator:     e.ChildText("div.creator-names > span"),
		}
		// Iterate over rows of the table which contains different information
		// about the course
		e.DOM.Find("table.basic-info-table tr").Each(func(_ int, s *goquery.Selection) {
			switch s.Find("td:first-child").Text() {
			case "Language":
				course.Language = s.Find("td:nth-child(2)").Text()
			case "Level":
				course.Level = s.Find("td:nth-child(2)").Text()
			case "Commitment":
				course.Commitment = s.Find("td:nth-child(2)").Text()
			case "How To Pass":
				course.HowToPass = s.Find("td:nth-child(2)").Text()
			case "User Ratings":
				course.Rating = s.Find("td:nth-child(2) div:nth-of-type(2)").Text()
			}
		})
		courses = append(courses, course)
	})

	// Start scraping on http://coursera.com/browse
	c.Visit("https://coursera.org/browse")

	// Convert results to JSON data if the scraping job has finished
	jsonData, err := json.MarshalIndent(courses, "", "  ")
	if err != nil {
		panic(err)
	}

	// Dump json to the standard output (can be redirected to a file)
	fmt.Println(string(jsonData))
}
