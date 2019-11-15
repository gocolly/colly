// Copyright 2018 Adam Tauber
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	// Index stores the position of the current element within all the elements matched by an OnHTML callback
	Index int
}

// NewHTMLElementFromSelectionNode creates a HTMLElement from a goquery.Selection Node.
func NewHTMLElementFromSelectionNode(resp *Response, s *goquery.Selection, n *html.Node, idx int) *HTMLElement {
	return &HTMLElement{
		Name:       n.Data,
		Request:    resp.Request,
		Response:   resp,
		Text:       goquery.NewDocumentFromNode(n).Text(),
		DOM:        s,
		Index:      idx,
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

// ChildTexts returns the stripped text content of all the matching
// elements.
func (h *HTMLElement) ChildTexts(goquerySelector string) []string {
	var res []string
	h.DOM.Find(goquerySelector).Each(func(_ int, s *goquery.Selection) {

		res = append(res, strings.TrimSpace(s.Text()))
	})
	return res
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
	var res []string
	h.DOM.Find(goquerySelector).Each(func(_ int, s *goquery.Selection) {
		if attr, ok := s.Attr(attrName); ok {
			res = append(res, strings.TrimSpace(attr))
		}
	})
	return res
}

// ForEach iterates over the elements matched by the first argument
// and calls the callback function on every HTMLElement match.
func (h *HTMLElement) ForEach(goquerySelector string, callback func(int, *HTMLElement)) {
	i := 0
	h.DOM.Find(goquerySelector).Each(func(_ int, s *goquery.Selection) {
		for _, n := range s.Nodes {
			callback(i, NewHTMLElementFromSelectionNode(h.Response, s, n, i))
			i++
		}
	})
}

// ForEachWithBreak iterates over the elements matched by the first argument
// and calls the callback function on every HTMLElement match.
// It is identical to ForEach except that it is possible to break
// out of the loop by returning false in the callback function. It returns the
// current Selection object.
func (h *HTMLElement) ForEachWithBreak(goquerySelector string, callback func(int, *HTMLElement) bool) {
	i := 0
	h.DOM.Find(goquerySelector).EachWithBreak(func(_ int, s *goquery.Selection) bool {
		for _, n := range s.Nodes {
			if callback(i, NewHTMLElementFromSelectionNode(h.Response, s, n, i)) {
				i++
				return true
			}
		}
		return false
	})
}
