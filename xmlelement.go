package colly

import (
	"strings"

	"encoding/xml"

	"github.com/antchfx/xmlquery"
)

// XMLElement is the representation of a XML tag.
type XMLElement struct {
	// Name is the name of the tag
	Name       string
	Text       string
	attributes []xml.Attr
	// Request is the request object of the element's HTML document
	Request *Request
	// Response is the Response object of the element's HTML document
	Response *Response
	// DOM is the goquery parsed DOM object of the page. DOM is relative
	// to the current HTMLElement
	DOM *xmlquery.Node
}

// NewXMLElementFromXMLNode creates a XMLElement from a xmlquery.Node.
func NewXMLElementFromXMLNode(resp *Response, s *xmlquery.Node) *XMLElement {
	return &XMLElement{
		Name:       s.Data,
		Request:    resp.Request,
		Response:   resp,
		Text:       s.InnerText(),
		DOM:        s,
		attributes: s.Attr,
	}
}

// Attr returns the selected attribute of a HTMLElement or empty string
// if no attribute found
func (h *XMLElement) Attr(k string) string {
	for _, a := range h.attributes {
		if a.Name.Local == k {
			return a.Value
		}
	}
	return ""
}

// ChildText returns the concatenated and stripped text content of the matching
// elements.
func (h *XMLElement) ChildText(xpathQuery string) string {
	return strings.TrimSpace(xmlquery.FindOne(h.DOM, xpathQuery).InnerText())
}

// ChildAttr returns the stripped text content of the first matching
// element's attribute.
func (h *XMLElement) ChildAttr(xpathQuery, attrName string) string {
	attr := xmlquery.FindOne(h.DOM, xpathQuery).SelectAttr(attrName)
	return strings.TrimSpace(attr)
}

// ChildAttrs returns the stripped text content of all the matching
// element's attributes.
func (h *XMLElement) ChildAttrs(xpathQuery, attrName string) []string {
	res := make([]string, 0)
	xmlquery.FindEach(h.DOM, xpathQuery, func(i int, s *xmlquery.Node) {
		if attr := s.SelectAttr(attrName); attr != "" {
			res = append(res, strings.TrimSpace(attr))
		}
	})
	return res
}
