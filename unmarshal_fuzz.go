package colly

import (
	"bytes"
	"testing"

	fuzz "github.com/AdamKorcz/go-fuzz-headers"
	"github.com/PuerkitoBio/goquery"
)

type info struct {
	Text string `selector:"span"`
}

type object struct {
	Info []*info `selector:"li.info"`
}

func FuzzHtmlElementUnmarshal(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		gfh := fuzz.NewConsumer(data)
		e := &HTMLElement{}
		err := gfh.GenerateStruct(e)
		if err != nil {
			return
		}
		d2, err := gfh.GetBytes()
		if err != nil {
			return
		}
		doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(d2))
		if err != nil {
			return
		}
		e.DOM = doc.First()
		o := object{}
		e.Unmarshal(&o)
	})
}
