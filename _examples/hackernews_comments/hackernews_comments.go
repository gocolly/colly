package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/asciimoo/colly"
)

type comment struct {
	Author  string
	URL     string
	Comment string
	Replies []*comment
	depth   int
}

func main() {
	var itemID string
	flag.StringVar(&itemID, "id", "", "hackernews post id")
	flag.Parse()

	if itemID == "" {
		log.Println("Hackernews post id required")
		os.Exit(1)
	}

	comments := make([]*comment, 0)

	// Instantiate default collector
	c := colly.NewCollector()

	// Extract comment
	c.OnHTML(".comment-tree tr.athing", func(e *colly.HTMLElement) {
		width, err := strconv.Atoi(e.ChildAttr("td.ind img", "width"))
		if err != nil {
			return
		}
		// hackernews uses 40px spacers to indent comment replies,
		// so we have to divide the width with it to get the depth
		// of the comment
		depth := width / 40
		c := &comment{
			Author:  e.ChildText("a.hnuser"),
			URL:     e.ChildAttr(".age a[href]", "href"),
			Comment: e.ChildText(".comment"),
			Replies: make([]*comment, 0),
			depth:   depth,
		}
		if depth == 0 {
			comments = append(comments, c)
			return
		}
		parent := comments[len(comments)-1]
		// append comment to its parent
		for i := 0; i < depth-1; i++ {
			parent = parent.Replies[len(parent.Replies)-1]
		}
		parent.Replies = append(parent.Replies, c)
	})

	c.Visit("https://news.ycombinator.com/item?id=" + itemID)

	// Convert results to JSON data if the scraping job has finished
	jsonData, err := json.MarshalIndent(comments, "", "  ")
	if err != nil {
		panic(err)
	}

	// Dump json to the standard output (can be redirected to a file)
	fmt.Println(string(jsonData))
}
