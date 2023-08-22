package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// DATE_FORMAT default format date used in openedx
const DATE_FORMAT = "02 Jan, 2006"

// Course store openedx course data
type Course struct {
	CourseID  string
	Run       string
	Name      string
	Number    string
	StartDate *time.Time
	EndDate   *time.Time
	URL       string
}

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Using IndonesiaX as sample
		colly.AllowedDomains("indonesiax.co.id", "www.indonesiax.co.id"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./cache"),
	)

	courses := make([]Course, 0, 200)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", "_", func(_ string, e *colly.HTMLElement) {
		link := e.Attr("href")
		if !strings.HasPrefix(link, "/courses/") {
			return
		}
		// start scraping the page under the link found
		e.Request.Visit(link)
	})

	c.OnHTML("div[class=main-container]", "_", func(_ string, e *colly.HTMLElement) {
		if e.DOM.Find("section#course-info").Length() == 0 {
			return
		}
		title := strings.Split(e.ChildText(".course-info__title"), "\n")[0]
		course_id := e.ChildAttr("input[name=course_id]", "value")
		texts := e.ChildTexts("span[data-datetime]")
		start_date, _ := time.Parse(DATE_FORMAT, texts[0])
		end_date, _ := time.Parse(DATE_FORMAT, texts[1])
		var run string
		if len(strings.Split(course_id, "_")) > 1 {
			run = strings.Split(course_id, "_")[1]
		}
		course := Course{
			CourseID:  course_id,
			Run:       run,
			Name:      title,
			Number:    e.ChildText("span.course-number"),
			StartDate: &start_date,
			EndDate:   &end_date,
			URL:       fmt.Sprintf("/courses/%s/about", course_id),
		}
		courses = append(courses, course)
	})

	// Start scraping on https://openedxdomain/courses
	c.Visit("https://www.indonesiax.co.id/courses")

	// Convert results to JSON data if the scraping job has finished
	jsonData, err := json.MarshalIndent(courses, "", "  ")
	if err != nil {
		panic(err)
	}

	// Dump json to the standard output (can be redirected to a file)
	fmt.Println(string(jsonData))
}
