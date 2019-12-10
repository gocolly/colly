package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

// Mail is the container of a single e-mail
type Mail struct {
	Title   string
	Link    string
	Author  string
	Date    string
	Message string
}

func main() {
	var groupName string
	flag.StringVar(&groupName, "group", "hspbp", "Google Groups group name")
	flag.Parse()

	threads := make(map[string][]Mail)

	threadCollector := colly.NewCollector()
	mailCollector := colly.NewCollector()

	// Collect threads
	threadCollector.OnHTML("tr", func(e *colly.HTMLElement) {
		ch := e.DOM.Children()
		author := ch.Eq(1).Text()
		// deleted topic
		if author == "" {
			return
		}

		title := ch.Eq(0).Text()
		link, _ := ch.Eq(0).Children().Eq(0).Attr("href")
		// fix link to point to the pure HTML version of the thread
		link = strings.Replace(link, ".com/d/topic", ".com/forum/?_escaped_fragment_=topic", 1)
		date := ch.Eq(2).Text()

		log.Printf("Thread found: %s %q %s %s\n", link, title, author, date)
		mailCollector.Visit(link)
	})

	// Visit next page
	threadCollector.OnHTML("body > a[href]", func(e *colly.HTMLElement) {
		log.Println("Next page link found:", e.Attr("href"))
		e.Request.Visit(e.Attr("href"))
	})

	// Extract mails
	mailCollector.OnHTML("body", func(e *colly.HTMLElement) {
		// Find subject
		threadSubject := e.ChildText("h2")
		if _, ok := threads[threadSubject]; !ok {
			threads[threadSubject] = make([]Mail, 0, 8)
		}

		// Extract mails
		e.ForEach("table tr", func(_ int, el *colly.HTMLElement) {
			mail := Mail{
				Title:   el.ChildText("td:nth-of-type(1)"),
				Link:    el.ChildAttr("td:nth-of-type(1)", "href"),
				Author:  el.ChildText("td:nth-of-type(2)"),
				Date:    el.ChildText("td:nth-of-type(3)"),
				Message: el.ChildText("td:nth-of-type(4)"),
			}
			threads[threadSubject] = append(threads[threadSubject], mail)
		})

		// Follow next page link
		if link, found := e.DOM.Find("> a[href]").Attr("href"); found {
			e.Request.Visit(link)
		} else {
			log.Printf("Thread %q done\n", threadSubject)
		}
	})

	threadCollector.Visit("https://groups.google.com/forum/?_escaped_fragment_=forum/" + groupName)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	enc.Encode(threads)
}
