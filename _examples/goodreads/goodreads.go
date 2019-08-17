package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

func main() {

// create file 
	fileName := "quote.txt"
	file, errFile := os.Create(fileName)
	if errFile != nil {
		println("operating system create file error :%s", errFile.Error())
		panic(errFile)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			println("file close error")
		}
	}()

	c := colly.NewCollector()
  
  // optianl: if you cannot connect https://www.goodread.com, then set a proper proxy.
	errProxy := c.SetProxy("http://127.0.0.1:1080/")
	if errProxy != nil {
		println("colly set proxy error :%s", errProxy.Error())
		panic(errProxy)
	}

	c.AllowURLRevisit = true
	extensions.RandomUserAgent(c)

	c.OnHTML(".quoteText ",
		func(e *colly.HTMLElement) {
			text := strings.TrimSpace(strings.Split(e.Text, "―")[0])
			author := TrimSpaceNewlineInString(strings.TrimSpace(e.ChildText(".authorOrTitle")))

			fileWriteForMarkdown(file, text, author)
		})

	c.OnHTML(".next_page", func(e *colly.HTMLElement) {
		println("visit: ", e.Request.AbsoluteURL(e.Attr("href")))
		errHrefVisit := c.Visit(e.Request.AbsoluteURL(e.Attr("href")))
		if errHrefVisit != nil {
			panic(errHrefVisit)
		}

	})

	errVisit := c.Visit("https://www.goodreads.com/quotes/tag/philosophy")
	if errVisit != nil {
		panic(errVisit)
	}

}

// because origin response string  has newline in it, so trim these. 
func TrimSpaceNewlineInString(s string) string {
	re := regexp.MustCompile(`\n`)
	return re.ReplaceAllString(s, " ")
}

func fileWriteDirect(file *os.File,lines ...string){

	_, err := (*file).Write([]byte(lines[0]))
	if err != nil {
		println("file write error ", err.Error())
	}
	_, err = (*file).Write([]byte(lines[1]))
	if err != nil {
		println("file write error ", err.Error())
	}
}
