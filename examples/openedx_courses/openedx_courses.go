package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/asciimoo/colly"
)

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
	c := colly.NewCollector()
	// Using IndonesiaX as sample
	c.AllowedDomains = []string{"indonesiax.co.id", "www.indonesiax.co.id"}

	// Cache responses to prevent multiple download of pages
	// even if the collector is restarted
	c.CacheDir = "./cache"

	courses := make([]Course, 0, 200)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if !strings.HasPrefix(link, "/courses/") {
			return
		}
		// start scaping the page under the link found
		e.Request.Visit(link)
	})

	c.OnHTML("div[class=content-wrapper]", func(e *colly.HTMLElement) {
		if e.DOM.Find("section.course-info").Length() == 0 {
			return
		}
		title := strings.Split(e.ChildText(".course-title"), "\n")[0]
		//provider := e.DOM.Find(".course-title > a").First()
		course_id := e.DOM.Find("input[name=course_id]").First().AttrOr("value", "")
		var run string
		if len(strings.Split(course_id, "_")) > 1 {
			run = strings.Split(course_id, "_")[1]
		}
		course := Course{
			CourseID: course_id,
			Run:      run,
			Name:     title,
			URL:      fmt.Sprintf("/courses/%s/about", course_id),
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
