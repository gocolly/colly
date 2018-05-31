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
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
	"regexp"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

// Response is the representation of a HTTP response made by a Collector
type Response struct {
	// StatusCode is the status code of the Response
	StatusCode int
	// Body is the content of the Response
	Body []byte
	// Ctx is a context between a Request and a Response
	Ctx *Context
	// Request is the Request object of the response
	Request *Request
	// Headers contains the Response's HTTP headers
	Headers *http.Header
}

// Save writes response body to disk
func (r *Response) Save(fileName string) error {
	return ioutil.WriteFile(fileName, r.Body, 0644)
}

// FileName returns the sanitized file name parsed from "Content-Disposition"
// header or from URL
func (r *Response) FileName() string {
	_, params, err := mime.ParseMediaType(r.Headers.Get("Content-Disposition"))
	if fName, ok := params["filename"]; ok && err == nil {
		return SanitizeFileName(fName)
	}
	if r.Request.URL.RawQuery != "" {
		return SanitizeFileName(fmt.Sprintf("%s_%s", r.Request.URL.Path, r.Request.URL.RawQuery))
	}
	return SanitizeFileName(strings.TrimPrefix(r.Request.URL.Path, "/"))
}

// regex from https://github.com/apache/tika/blob/0e8f44459fbbed991171e9eafb3395df6060fb7a/tika-parsers/src/main/java/org/apache/tika/parser/html/HtmlEncodingDetector.java#L99
var charsetRe = regexp.MustCompile(`(?i)charset\s*=\s*(?:['\\"]\\s*)?([-_:\\.a-z0-9]+)`)
var httpMetaPattern = regexp.MustCompile("(?is)<\\s*meta\\s+([^<>]+)")

func headEncoding(response *Response) string { // get encoding from head
	contentType := response.Headers.Get("content-type")
	rtnValue := ""
	if len(contentType) > 0 {
		if strings.Contains(contentType, "charset") {
			re := regexp.MustCompile(`(?i)charset=(?P<charset>.*)`)
			a := re.FindSubmatch([]byte(contentType))
			if len(a) > 0 {
				rtnValue = string(a[1])
			}
		}
	}
	return rtnValue
}

func bodyEncoding(response *Response) string {
	maxSize := 1024 * 10
	temp := make([]byte, maxSize)
	if len(response.Body) > maxSize {
		temp = response.Body[1:maxSize]
	} else {
		temp = response.Body
	}
	metaTags := httpMetaPattern.FindAll(temp, -1)
	for _, i := range metaTags {
		cs := charsetRe.FindSubmatch(i)
		return string(cs[1])
	}
	return ""
}
func getEncoding(response *Response) string {
	head_encoding := headEncoding(response)
	if (len(head_encoding) > 0) {
		return head_encoding
	} else {
		return bodyEncoding(response)
	}
}

func (r *Response) fixCharset(detectCharset bool, defaultEncoding string) error {
	if defaultEncoding != "" {
		tmpBody, err := encodeBytes(r.Body, "text/plain; charset="+defaultEncoding)
		if err != nil {
			return err
		}
		r.Body = tmpBody
		return nil
	}
	contentType := strings.ToLower(getEncoding(r))

	if len(contentType) == 0 { // no charset found
		if !detectCharset {
			return nil
		}
		d := chardet.NewTextDetector()
		r, err := d.DetectBest(r.Body)
		if err != nil {
			return err
		}
		contentType = "text/plain; charset=" + r.Charset
	}
	if strings.Compare(contentType, "utf-8") == 0 || strings.Compare(contentType, "utf8") == 0 {
		return nil
	}
	tmpBody, err := encodeBytes(r.Body, contentType)
	if err != nil {
		return err
	}
	r.Body = tmpBody
	return nil
}

func encodeBytes(b []byte, contentType string) ([]byte, error) {
	r, err := charset.NewReader(bytes.NewReader(b), contentType)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}
