package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jawher/mow.cli"
)

var scraper_head_template string = `package main

import (
	"log"

	"github.com/asciimoo/colly"
)

func main() {
	c := colly.NewCollector()
`

var scraper_end_template string = `
	c.Visit("https://yourdomain.com/")
}
`

var html_callback_template string = `
	c.OnHTML("element-selector", func(e *colly.HTMLElement) {
		log.Println(e.Text)
	})
`

var request_callback_template string = `
	c.OnRequest("element-selector", func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})
`

var response_callback_template string = `
	c.OnResponse("element-selector", func(r *colly.Response) {
		log.Println("Visited", r.Request.URL, r.StatusCode)
	})
`

var error_callback_template string = `
	c.OnError("element-selector", func(r *colly.Response, err error) {
		log.Printf("Error on %s: %s", r.Request.URL, err)
	})
`

func main() {
	app := cli.App("colly", "Scraping Framework for Gophers")

	app.Command("new", "Create new scraper", func(cmd *cli.Cmd) {
		var (
			callbacks = cmd.StringOpt("callbacks", "", "Add callbacks to the template. (E.g. '--callbacks=html,response,error')")
			hosts     = cmd.StringOpt("hosts", "", "Specify scraper's allowed hosts. (e.g. '--hosts=xy.com,abcd.com')")
			path      = cmd.StringArg("PATH", "", "Path of the new scraper")
		)

		cmd.Spec = "[--callbacks] [--hosts] [PATH]"

		cmd.Action = func() {
			scraper := bytes.NewBufferString(scraper_head_template)
			outfile := os.Stdout
			if *path != "" {
				var err error
				outfile, err = os.Create(*path)
				if err != nil {
					log.Fatal(err)
				}
				defer outfile.Close()
			}
			if *hosts != "" {
				scraper.WriteString("\n	c.AllowedDomains = []string{")
				for i, h := range strings.Split(*hosts, ",") {
					if i > 0 {
						scraper.WriteString(", ")
					}
					scraper.WriteString(fmt.Sprintf("%q", h))
				}
				scraper.WriteString("}\n")
			}
			if len(*callbacks) > 0 {
				for _, c := range strings.Split(*callbacks, ",") {
					switch c {
					case "html":
						scraper.WriteString(html_callback_template)
					case "request":
						scraper.WriteString(request_callback_template)
					case "response":
						scraper.WriteString(response_callback_template)
					case "error":
						scraper.WriteString(error_callback_template)
					}
				}
			}
			scraper.WriteString(scraper_end_template)
			outfile.Write(scraper.Bytes())
		}
	})

	app.Run(os.Args)
}
