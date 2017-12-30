package colly

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// HTMLElement is the representation of a HTML tag.
type HTMLElement struct {
	// Name is the name of the tag
	Name       string
	Text       string
	attributes []html.Attribute
	// Request is the request object of the element's HTML document
	Request *Request
	// Response is the Response object of the element's HTML document
	Response *Response
	// DOM is the goquery parsed DOM object of the page. DOM is relative
	// to the current HTMLElement
	DOM *goquery.Selection
}

// NewHTMLElementFromSelectionNode creates a HTMLElement from a goquery.Selection Node.
func NewHTMLElementFromSelectionNode(resp *Response, s *goquery.Selection, n *html.Node) *HTMLElement {
	return &HTMLElement{
		Name:       n.Data,
		Request:    resp.Request,
		Response:   resp,
		Text:       goquery.NewDocumentFromNode(n).Text(),
		DOM:        s,
		attributes: n.Attr,
	}
}

// Attr returns the selected attribute of a HTMLElement or empty string
// if no attribute found
func (h *HTMLElement) Attr(k string) string {
	for _, a := range h.attributes {
		if a.Key == k {
			return a.Val
		}
	}
	return ""
}

// ChildText returns the concatenated and stripped text content of the matching
// elements.
func (h *HTMLElement) ChildText(goquerySelector string) string {
	return strings.TrimSpace(h.DOM.Find(goquerySelector).Text())
}

// ChildAttr returns the stripped text content of the first matching
// element's attribute.
func (h *HTMLElement) ChildAttr(goquerySelector, attrName string) string {
	if attr, ok := h.DOM.Find(goquerySelector).Attr(attrName); ok {
		return strings.TrimSpace(attr)
	}
	return ""
}

// ChildAttrs returns the stripped text content of all the matching
// element's attributes.
func (h *HTMLElement) ChildAttrs(goquerySelector, attrName string) []string {
	res := make([]string, 0)
	h.DOM.Find(goquerySelector).Each(func(_ int, s *goquery.Selection) {
		if attr, ok := s.Attr(attrName); ok {
			res = append(res, strings.TrimSpace(attr))
		}
	})
	return res
}
