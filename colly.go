package colly

import (
	//"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html"

	"github.com/PuerkitoBio/goquery"
)

type Collector struct {
	Cache          bool
	UserAgent      string
	AllowedDomains []string
	visitedURLs    []string
	htmlCallbacks  map[string]HTMLCallback
}

type HTMLElement struct {
	Name       string
	attributes []html.Attribute
}

type HTMLCallback func(*HTMLElement)

func NewCollector() *Collector {
	return &Collector{
		Cache:          true,
		UserAgent:      "colly",
		AllowedDomains: nil,
		visitedURLs:    make([]string, 0),
		htmlCallbacks:  make(map[string]HTMLCallback, 0),
	}
}

func (c *Collector) OnHTML(xpath string, f HTMLCallback) {
	c.htmlCallbacks[xpath] = f
}

func (c *Collector) Visit(url string) error {
	// TODO create request
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromResponse(res)
	//body, err := ioutil.ReadAll(res.Body)
	//res.Body.Close()
	if err != nil {
		return err
	}
	if strings.Index(strings.ToLower(res.Header.Get("Content-Type")), "html") > -1 {
		for expr, f := range c.htmlCallbacks {
			doc.Find(expr).Each(func(i int, s *goquery.Selection) {
				for _, n := range s.Nodes {
					f(&HTMLElement{attributes: n.Attr, Name: n.Data})
				}
			})
		}
	}
	return nil
}

func (h *HTMLElement) Attr(k string) string {
	for _, a := range h.attributes {
		if a.Key == k {
			return a.Val
		}
	}
	return ""
}
