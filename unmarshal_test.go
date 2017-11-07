package colly

import (
	"bytes"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var basicTestData []byte = []byte(`<ul><li class="x">list <span>item</span> 1</li><li>list item 2</li><li>3</li></ul>`)
var nestedTestData []byte = []byte(`<div><p>a</p><div><p>b</p><div><p>c</p></div></div></div>`)

func TestBasicUnmarshal(t *testing.T) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewBuffer(basicTestData))
	e := &HTMLElement{
		DOM: doc.First(),
	}
	s := struct {
		String string   `selector:"li:first-child" attr:"class"`
		Items  []string `selector:"li"`
		Struct struct {
			String string `selector:"li:last-child"`
		}
	}{}
	if err := e.Unmarshal(&s); err != nil {
		t.Error("Cannot unmarshal struct: " + err.Error())
	}
	if s.String != "x" {
		t.Errorf(`Invalid data for String: %q, expected "x"`, s.String)
	}
	if s.Struct.String != "3" {
		t.Errorf(`Invalid data for Struct.String: %q, expected "3"`, s.Struct.String)
	}
}

func TestNestedUnmarshal(t *testing.T) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewBuffer(nestedTestData))
	e := &HTMLElement{
		DOM: doc.First(),
	}
	type nested struct {
		String string  `selector:"div > p"`
		Struct *nested `selector:"div > div"`
	}
	s := nested{}
	if err := e.Unmarshal(&s); err != nil {
		t.Error("Cannot unmarshal struct: " + err.Error())
	}
	if s.String != "a" {
		t.Errorf(`Invalid data for String: %q, expected "a"`, s.String)
	}
	if s.Struct.String != "b" {
		t.Errorf(`Invalid data for Struct.String: %q, expected "b"`, s.Struct.String)
	}
	if s.Struct.Struct.String != "c" {
		t.Errorf(`Invalid data for Struct.Struct.String: %q, expected "c"`, s.Struct.Struct.String)
	}
}
