package colly

import (
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

// XMLElement is the representation of a XML tag.
type XMLElement struct {
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
	DOM *html.Node
}

// NewXMLElementFromHTMLNode creates a XMLElement from a xmlquery.Node.
func NewXMLElementFromHTMLNode(resp *Response, s *html.Node) *XMLElement {
	return &XMLElement{
		Name:       s.Data,
		Request:    resp.Request,
		Response:   resp,
		Text:       htmlquery.InnerText(s),
		DOM:        s,
		attributes: s.Attr,
	}
}

// Attr returns the selected attribute of a HTMLElement or empty string
// if no attribute found
func (h *XMLElement) Attr(k string) string {
	for _, a := range h.attributes {
		if a.Key == k {
			return a.Val
		}
	}
	return ""
}

// ChildText returns the concatenated and stripped text content of the matching
// elements.
func (h *XMLElement) ChildText(xpathQuery string) string {
	return strings.TrimSpace(htmlquery.InnerText(htmlquery.FindOne(h.DOM, xpathQuery)))
}

// ChildAttr returns the stripped text content of the first matching
// element's attribute.
func (h *XMLElement) ChildAttr(xpathQuery, attrName string) string {
	child := htmlquery.FindOne(h.DOM, xpathQuery)
	if child != nil {
		for _, attr := range child.Attr {
			if attr.Key == attrName {
				return strings.TrimSpace(attr.Val)
			}
		}
	}

	return ""
}

// ChildAttrs returns the stripped text content of all the matching
// element's attributes.
func (h *XMLElement) ChildAttrs(xpathQuery, attrName string) []string {
	res := make([]string, 0)
	htmlquery.FindEach(h.DOM, xpathQuery, func(i int, child *html.Node) {
		for _, attr := range child.Attr {
			if attr.Key == attrName {
				res = append(res, strings.TrimSpace(attr.Val))
			}
		}
	})
	return res
}
