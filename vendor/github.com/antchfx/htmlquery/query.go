/*
Package htmlquery provides extract data from HTML documents using XPath expression.
*/
package htmlquery

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

var _ xpath.NodeNavigator = &NodeNavigator{}

// CreateXPathNavigator creates a new xpath.NodeNavigator for the specified html.Node.
func CreateXPathNavigator(top *html.Node) *NodeNavigator {
	return &NodeNavigator{curr: top, root: top, attr: -1}
}

// Find searches the html.Node that matches by the specified XPath expr.
func Find(top *html.Node, expr string) []*html.Node {
	var elems []*html.Node
	exp, err := xpath.Compile(expr)
	if err != nil {
		panic(err)
	}
	t := exp.Select(CreateXPathNavigator(top))
	for t.MoveNext() {
		elems = append(elems, (t.Current().(*NodeNavigator)).curr)
	}
	return elems
}

// FindOne searches the html.Node that matches by the specified XPath expr,
// and returns first element of matched html.Node.
func FindOne(top *html.Node, expr string) *html.Node {
	var elem *html.Node
	exp, err := xpath.Compile(expr)
	if err != nil {
		panic(err)
	}
	t := exp.Select(CreateXPathNavigator(top))
	if t.MoveNext() {
		elem = (t.Current().(*NodeNavigator)).curr
	}
	return elem
}

// FindEach searches the html.Node and calls functions cb.
func FindEach(top *html.Node, expr string, cb func(int, *html.Node)) {
	exp, err := xpath.Compile(expr)
	if err != nil {
		panic(err)
	}
	t := exp.Select(CreateXPathNavigator(top))
	i := 0
	for t.MoveNext() {
		cb(i, (t.Current().(*NodeNavigator)).curr)
		i++
	}
}

// LoadURL loads the HTML document from the specified URL.
func LoadURL(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	return html.Parse(r)
}

// Parse returns the parse tree for the HTML from the given Reader.
func Parse(r io.Reader) (*html.Node, error) {
	return html.Parse(r)
}

// InnerText returns the text between the start and end tags of the object.
func InnerText(n *html.Node) string {
	var output func(*bytes.Buffer, *html.Node)
	output = func(buf *bytes.Buffer, n *html.Node) {
		switch n.Type {
		case html.TextNode:
			buf.WriteString(n.Data)
			return
		case html.CommentNode:
			return
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			output(buf, child)
		}
	}

	var buf bytes.Buffer
	output(&buf, n)
	return buf.String()
}

// SelectAttr returns the attribute value with the specified name.
func SelectAttr(n *html.Node, name string) (val string) {
	if n == nil {
		return
	}
	for _, attr := range n.Attr {
		if attr.Key == name {
			val = attr.Val
			break
		}
	}
	return
}

// OutputHTML returns the text including tags name.
func OutputHTML(n *html.Node, self bool) string {
	var buf bytes.Buffer
	if self {
		html.Render(&buf, n)
	} else {
		for n := n.FirstChild; n != nil; n = n.NextSibling {
			html.Render(&buf, n)
		}
	}
	return buf.String()
}

type NodeNavigator struct {
	root, curr *html.Node
	attr       int
}

func (h *NodeNavigator) Current() *html.Node {
	return h.curr
}

func (h *NodeNavigator) NodeType() xpath.NodeType {
	switch h.curr.Type {
	case html.CommentNode:
		return xpath.CommentNode
	case html.TextNode:
		return xpath.TextNode
	case html.DocumentNode:
		return xpath.RootNode
	case html.ElementNode:
		if h.attr != -1 {
			return xpath.AttributeNode
		}
		return xpath.ElementNode
	case html.DoctypeNode:
		// ignored <!DOCTYPE HTML> declare and as Root-Node type.
		return xpath.RootNode
	}
	panic(fmt.Sprintf("unknown HTML node type: %v", h.curr.Type))
}

func (h *NodeNavigator) LocalName() string {
	if h.attr != -1 {
		return h.curr.Attr[h.attr].Key
	}
	return h.curr.Data
}

func (*NodeNavigator) Prefix() string {
	return ""
}

func (h *NodeNavigator) Value() string {
	switch h.curr.Type {
	case html.CommentNode:
		return h.curr.Data
	case html.ElementNode:
		if h.attr != -1 {
			return h.curr.Attr[h.attr].Val
		}
		return InnerText(h.curr)
	case html.TextNode:
		return h.curr.Data
	}
	return ""
}

func (h *NodeNavigator) Copy() xpath.NodeNavigator {
	n := *h
	return &n
}

func (h *NodeNavigator) MoveToRoot() {
	h.curr = h.root
}

func (h *NodeNavigator) MoveToParent() bool {
	if h.attr != -1 {
		h.attr = -1
		return true
	} else if node := h.curr.Parent; node != nil {
		h.curr = node
		return true
	}
	return false
}

func (h *NodeNavigator) MoveToNextAttribute() bool {
	if h.attr >= len(h.curr.Attr)-1 {
		return false
	}
	h.attr++
	return true
}

func (h *NodeNavigator) MoveToChild() bool {
	if h.attr != -1 {
		return false
	}
	if node := h.curr.FirstChild; node != nil {
		h.curr = node
		return true
	}
	return false
}

func (h *NodeNavigator) MoveToFirst() bool {
	if h.attr != -1 || h.curr.PrevSibling == nil {
		return false
	}
	for {
		node := h.curr.PrevSibling
		if node == nil {
			break
		}
		h.curr = node
	}
	return true
}

func (h *NodeNavigator) String() string {
	return h.Value()
}

func (h *NodeNavigator) MoveToNext() bool {
	if h.attr != -1 {
		return false
	}
	if node := h.curr.NextSibling; node != nil {
		h.curr = node
		return true
	}
	return false
}

func (h *NodeNavigator) MoveToPrevious() bool {
	if h.attr != -1 {
		return false
	}
	if node := h.curr.PrevSibling; node != nil {
		h.curr = node
		return true
	}
	return false
}

func (h *NodeNavigator) MoveTo(other xpath.NodeNavigator) bool {
	node, ok := other.(*NodeNavigator)
	if !ok || node.root != h.root {
		return false
	}

	h.curr = node.curr
	h.attr = node.attr
	return true
}
