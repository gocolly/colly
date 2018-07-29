htmlquery
====
[![Build Status](https://travis-ci.org/antchfx/htmlquery.svg?branch=master)](https://travis-ci.org/antchfx/htmlquery)
[![Coverage Status](https://coveralls.io/repos/github/antchfx/htmlquery/badge.svg?branch=master)](https://coveralls.io/github/antchfx/htmlquery?branch=master)
[![GoDoc](https://godoc.org/github.com/antchfx/htmlquery?status.svg)](https://godoc.org/github.com/antchfx/htmlquery)
[![Go Report Card](https://goreportcard.com/badge/github.com/antchfx/htmlquery)](https://goreportcard.com/report/github.com/antchfx/htmlquery)

Overview
====

htmlquery is an XPath query package for HTML, lets you extract data or evaluate from HTML documents by an XPath expression.

Installation
====

> $ go get github.com/antchfx/htmlquery

Getting Started
====

#### Load HTML document from URL.

```go
doc, err := htmlquery.LoadURL("http://example.com/")
```

#### Load HTML document from string.

```go
s := `<html>....</html>`
doc, err := htmlquery.Parse(strings.NewReader(s))
```

#### Find all A elements.

```go
list := htmlquery.Find(doc, "//a")
```

#### Find all A elements with href attribute.

```go
list := range htmlquery.Find(doc, "//a/@href")	
```

### Find the third A element.

```go
a := htmlquery.FindOne(doc, "//a[3]")
```

#### Evaluate the number of all IMG element.

```go
expr, _ := xpath.Compile("count(//img)")
v := expr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(float64)
fmt.Printf("total count is %f", v)
```

Quick Tutorial
===

```go
func main() {
	doc, err := htmlquery.LoadURL("https://www.bing.com/search?q=golang")
	if err != nil {
		panic(err)
	}
	// Find all news item.
	for i, n := range htmlquery.Find(doc, "//ol/li") {
		a := htmlquery.FindOne(n, "//a")
		fmt.Printf("%d %s(%s)\n", i, htmlquery.InnerText(a), htmlquery.SelectAttr(a, "href"))
	}
}
```

List of supported XPath query packages
===
|Name |Description |
|--------------------------|----------------|
|[htmlquery](https://github.com/antchfx/htmlquery) | XPath query package for the HTML document|
|[xmlquery](https://github.com/antchfx/xmlquery) | XPath query package for the XML document|
|[jsonquery](https://github.com/antchfx/jsonquery) | XPath query package for the JSON document|

Questions
===
Please let me know if you have any questions.
