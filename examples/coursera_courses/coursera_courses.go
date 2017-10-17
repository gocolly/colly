package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/asciimoo/colly"
)

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

	// Create another collector to scrape course details
	detailCollector := colly.NewCollector()

	// Visit only domains: coursera.org, www.coursera.org
	c.AllowedDomains = []string{"coursera.org", "www.coursera.org"}
	c.CacheDir = "./t"
	detailCollector.AllowedDomains = c.AllowedDomains
	detailCollector.CacheDir = c.CacheDir
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
		// Add to courses map where key is the absolute URL and the
		// values is the name of the course
		courseURL := e.Request.AbsoluteURL(e.Attr("href"))
		if strings.Index(courseURL, "coursera.org/learn") == -1 {
			return
		}
		detailCollector.Visit(courseURL)
	})

	// Extract details of the course
	detailCollector.OnHTML(`div[id=rendered-content]`, func(e *colly.HTMLElement) {
		log.Println("Course found", e.Request.URL)
		title := e.DOM.Find(".course-title").Text()
		if title == "" {
			log.Println("No title found", e.Request.URL)
		}
		course := Course{
			Title:       title,
			URL:         e.Request.URL.String(),
			Description: e.DOM.Find("div.content").Text(),
			Creator:     e.DOM.Find("div.creator-names > span").Text(),
		}
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
				log.Println("yo")
				course.Rating = s.Find("td:nth-child(2) div:nth-of-type(2)").Text()
			}
		})
		courses = append(courses, course)
	})

	// Start scraping on http://coursera.com/browse
	c.Visit("https://coursera.org/browse")

	jsonData, err := json.MarshalIndent(courses, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonData))
}
