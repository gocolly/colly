package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/gocolly/colly/v2"
)

var baseSearchURL = "https://factba.se/json/json-transcript.php?q=&f=&dt=&p="
var baseTranscriptURL = "https://factba.se/transcript/"

type result struct {
	Slug string `json:"slug"`
	Date string `json:"date"`
}

type results struct {
	Data []*result `json:"data"`
}

type transcript struct {
	Speaker string
	Text    string
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("factba.se"),
	)

	d := c.Clone()

	d.OnHTML("body", func(e *colly.HTMLElement) {
		t := make([]transcript, 0)
		e.ForEach(".topic-media-row", func(_ int, el *colly.HTMLElement) {
			t = append(t, transcript{
				Speaker: el.ChildText(".speaker-label"),
				Text:    el.ChildText(".transcript-text-block"),
			})
		})
		jsonData, err := json.MarshalIndent(t, "", "  ")
		if err != nil {
			return
		}
		ioutil.WriteFile(colly.SanitizeFileName(e.Request.Ctx.Get("date")+"_"+e.Request.Ctx.Get("slug"))+".json", jsonData, 0644)
	})

	stop := false
	c.OnResponse(func(r *colly.Response) {
		rs := &results{}
		err := json.Unmarshal(r.Body, rs)
		if err != nil || len(rs.Data) == 0 {
			stop = true
			return
		}
		for _, res := range rs.Data {
			u := baseTranscriptURL + res.Slug
			ctx := colly.NewContext()
			ctx.Put("date", res.Date)
			ctx.Put("slug", res.Slug)
			d.Request("GET", u, nil, ctx, nil)
		}
	})

	for i := 1; i < 1000; i++ {
		if stop {
			break
		}
		if err := c.Visit(baseSearchURL + strconv.Itoa(i)); err != nil {
			fmt.Println("Error:", err)
			break
		}
	}
}
