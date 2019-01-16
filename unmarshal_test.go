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
	"bytes"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var basicTestData = []byte(`<ul><li class="x">list <span>item</span> 1</li><li>list item 2</li><li>3</li></ul>`)
var nestedTestData = []byte(`<div><p>a</p><div><p>b</p><div><p>c</p></div></div></div>`)
var pointerSliceTestData = []byte(`<ul class="object"><li class="info">Information: <span>Info 1</span></li><li class="info">Information: <span>Info 2</span></li></ul>`)

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

func TestNestedUnmarshalMap(t *testing.T) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewBuffer(nestedTestData))
	e := &HTMLElement{
		DOM: doc.First(),
	}
	doc2, _ := goquery.NewDocumentFromReader(bytes.NewBuffer(basicTestData))
	e2 := &HTMLElement{
		DOM: doc2.First(),
	}
	type nested struct {
		String string
	}
	mapSelector := make(map[string]string)
	mapSelector["String"] = "div > p"

	mapSelector2 := make(map[string]string)
	mapSelector2["String"] = "span"

	s := nested{}
	s2 := nested{}
	if err := e.UnmarshalWithMap(&s, mapSelector); err != nil {
		t.Error("Cannot unmarshal struct: " + err.Error())
	}
	if err := e2.UnmarshalWithMap(&s2, mapSelector2); err != nil {
		t.Error("Cannot unmarshal struct: " + err.Error())
	}
	if s.String != "a" {
		t.Errorf(`Invalid data for String: %q, expected "a"`, s.String)
	}
	if s2.String != "item" {
		t.Errorf(`Invalid data for String: %q, expected "a"`, s.String)
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

func TestPointerSliceUnmarshall(t *testing.T) {
	type info struct {
		Text string `selector:"span"`
	}
	type object struct {
		Info []*info `selector:"li.info"`
	}

	doc, _ := goquery.NewDocumentFromReader(bytes.NewBuffer(pointerSliceTestData))
	e := HTMLElement{DOM: doc.First()}
	o := object{}
	err := e.Unmarshal(&o)
	if err != nil {
		t.Fatalf("Failed to unmarshal page: %s\n", err.Error())
	}

	if len(o.Info) != 2 {
		t.Errorf("Invalid length for Info: %d, expected 2", len(o.Info))
	}
	if o.Info[0].Text != "Info 1" {
		t.Errorf("Invalid data for Info.[0].Text: %s, expected Info 1", o.Info[0].Text)
	}
	if o.Info[1].Text != "Info 2" {
		t.Errorf("Invalid data for Info.[1].Text: %s, expected Info 2", o.Info[1].Text)
	}

}

func TestStructSliceUnmarshall(t *testing.T) {
	type info struct {
		Text string `selector:"span"`
	}
	type object struct {
		Info []info `selector:"li.info"`
	}

	doc, _ := goquery.NewDocumentFromReader(bytes.NewBuffer(pointerSliceTestData))
	e := HTMLElement{DOM: doc.First()}
	o := object{}
	err := e.Unmarshal(&o)
	if err != nil {
		t.Fatalf("Failed to unmarshal page: %s\n", err.Error())
	}

	if len(o.Info) != 2 {
		t.Errorf("Invalid length for Info: %d, expected 2", len(o.Info))
	}
	if o.Info[0].Text != "Info 1" {
		t.Errorf("Invalid data for Info.[0].Text: %s, expected Info 1", o.Info[0].Text)
	}
	if o.Info[1].Text != "Info 2" {
		t.Errorf("Invalid data for Info.[1].Text: %s, expected Info 2", o.Info[1].Text)
	}

}
