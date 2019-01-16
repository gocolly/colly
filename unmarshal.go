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
	"errors"
	"reflect"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Unmarshal is a shorthand for colly.UnmarshalHTML
func (h *HTMLElement) Unmarshal(v interface{}) error {
	return UnmarshalHTML(v, h.DOM, nil)
}

// UnmarshalWithMap is a shorthand for colly.UnmarshalHTML, extended to allow maps to be passed in.
func (h *HTMLElement) UnmarshalWithMap(v interface{}, structMap map[string]string) error {
	return UnmarshalHTML(v, h.DOM, structMap)
}

// UnmarshalHTML declaratively extracts text or attributes to a struct from
// HTML response using struct tags composed of css selectors.
// Allowed struct tags:
//  - "selector" (required): CSS (goquery) selector of the desired data
//  - "attr" (optional): Selects the matching element's attribute's value.
//     Leave it blank or omit to get the text of the element.
//
// Example struct declaration:
//
//   type Nested struct {
//   	String  string   `selector:"div > p"`
//      Classes []string `selector:"li" attr:"class"`
//   	Struct  *Nested  `selector:"div > div"`
//   }
//
// Supported types: struct, *struct, string, []string
func UnmarshalHTML(v interface{}, s *goquery.Selection, structMap map[string]string) error {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("Invalid type or nil-pointer")
	}

	sv := rv.Elem()
	st := reflect.TypeOf(v).Elem()
	if structMap != nil {
		for k, v := range structMap {
			attrV := sv.FieldByName(k)
			if !attrV.CanAddr() || !attrV.CanSet() {
				continue
			}
			if err := unmarshalSelector(s, attrV, v); err != nil {
				return err
			}
		}
	} else {
		for i := 0; i < sv.NumField(); i++ {
			attrV := sv.Field(i)
			if !attrV.CanAddr() || !attrV.CanSet() {
				continue
			}
			if err := unmarshalAttr(s, attrV, st.Field(i)); err != nil {
				return err
			}

		}
	}

	return nil
}

func unmarshalSelector(s *goquery.Selection, attrV reflect.Value, selector string) error {
	//selector is "-" specify that field should ignore.
	if selector == "-" {
		return nil
	}
	htmlAttr := ""
	// TODO support more types
	switch attrV.Kind() {
	case reflect.Slice:
		if err := unmarshalSlice(s, selector, htmlAttr, attrV); err != nil {
			return err
		}
	case reflect.String:
		val := getDOMValue(s.Find(selector), htmlAttr)
		attrV.Set(reflect.Indirect(reflect.ValueOf(val)))
	case reflect.Struct:
		if err := unmarshalStruct(s, selector, attrV); err != nil {
			return err
		}
	case reflect.Ptr:
		if err := unmarshalPtr(s, selector, attrV); err != nil {
			return err
		}
	default:
		return errors.New("Invalid type: " + attrV.String())
	}
	return nil
}

func unmarshalAttr(s *goquery.Selection, attrV reflect.Value, attrT reflect.StructField) error {
	selector := attrT.Tag.Get("selector")
	//selector is "-" specify that field should ignore.
	if selector == "-" {
		return nil
	}
	htmlAttr := attrT.Tag.Get("attr")
	// TODO support more types
	switch attrV.Kind() {
	case reflect.Slice:
		if err := unmarshalSlice(s, selector, htmlAttr, attrV); err != nil {
			return err
		}
	case reflect.String:
		val := getDOMValue(s.Find(selector), htmlAttr)
		attrV.Set(reflect.Indirect(reflect.ValueOf(val)))
	case reflect.Struct:
		if err := unmarshalStruct(s, selector, attrV); err != nil {
			return err
		}
	case reflect.Ptr:
		if err := unmarshalPtr(s, selector, attrV); err != nil {
			return err
		}
	default:
		return errors.New("Invalid type: " + attrV.String())
	}
	return nil
}

func unmarshalStruct(s *goquery.Selection, selector string, attrV reflect.Value) error {
	newS := s
	if selector != "" {
		newS = newS.Find(selector)
	}
	if newS.Nodes == nil {
		return nil
	}
	v := reflect.New(attrV.Type())
	err := UnmarshalHTML(v.Interface(), newS, nil)
	if err != nil {
		return err
	}
	attrV.Set(reflect.Indirect(v))
	return nil
}

func unmarshalPtr(s *goquery.Selection, selector string, attrV reflect.Value) error {
	newS := s
	if selector != "" {
		newS = newS.Find(selector)
	}
	if newS.Nodes == nil {
		return nil
	}
	e := attrV.Type().Elem()
	if e.Kind() != reflect.Struct {
		return errors.New("Invalid slice type")
	}
	v := reflect.New(e)
	err := UnmarshalHTML(v.Interface(), newS, nil)
	if err != nil {
		return err
	}
	attrV.Set(v)
	return nil
}

func unmarshalSlice(s *goquery.Selection, selector, htmlAttr string, attrV reflect.Value) error {
	if attrV.Pointer() == 0 {
		v := reflect.MakeSlice(attrV.Type(), 0, 0)
		attrV.Set(v)
	}
	switch attrV.Type().Elem().Kind() {
	case reflect.String:
		s.Find(selector).Each(func(_ int, s *goquery.Selection) {
			val := getDOMValue(s, htmlAttr)
			attrV.Set(reflect.Append(attrV, reflect.Indirect(reflect.ValueOf(val))))
		})
	case reflect.Ptr:
		s.Find(selector).Each(func(_ int, innerSel *goquery.Selection) {
			someVal := reflect.New(attrV.Type().Elem().Elem())
			UnmarshalHTML(someVal.Interface(), innerSel, nil)
			attrV.Set(reflect.Append(attrV, someVal))
		})
	case reflect.Struct:
		s.Find(selector).Each(func(_ int, innerSel *goquery.Selection) {
			someVal := reflect.New(attrV.Type().Elem())
			UnmarshalHTML(someVal.Interface(), innerSel, nil)
			attrV.Set(reflect.Append(attrV, reflect.Indirect(someVal)))
		})
	default:
		return errors.New("Invalid slice type")
	}
	return nil
}

func getDOMValue(s *goquery.Selection, attr string) string {
	if attr == "" {
		return strings.TrimSpace(s.First().Text())
	}
	attrV, _ := s.Attr(attr)
	return attrV
}
