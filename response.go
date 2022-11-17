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
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"

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
	// Trace contains the HTTPTrace for the request. Will only be set by the
	// collector if Collector.TraceHTTP is set to true.
	Trace *HTTPTrace
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

func (r *Response) fixCharset(detectCharset bool, defaultEncoding string) error {
	if len(r.Body) == 0 {
		return nil
	}
	if defaultEncoding != "" {
		tmpBody, err := encodeBytes(r.Body, "text/plain; charset="+defaultEncoding)
		if err != nil {
			return err
		}
		r.Body = tmpBody
		return nil
	}
	contentType := strings.ToLower(r.Headers.Get("Content-Type"))

	if strings.Contains(contentType, "image/") ||
		strings.Contains(contentType, "video/") ||
		strings.Contains(contentType, "audio/") ||
		strings.Contains(contentType, "font/") {
		// These MIME types should not have textual data.

		return nil
	}

	if !strings.Contains(contentType, "charset") {
		if !detectCharset {
			return nil
		}
		bDo := false
		s1 := string(r.Body)
		k1 := "content=\"text/html; charset="
		n1 := strings.Index(s1, k1)
		if -1 == n1 {
			k1 = "<meta charset="
			n1 = strings.Index(s1, k1)
		}
		if -1 < n1 {
			k1 = s1[n1+len(k1):]
			n1 = strings.Index(k1, ">")
			if -1 < n1 {
				k1 = k1[0:n1]
				n1 = strings.Index(k1, "\"")
				if -1 == n1 {
					n1 = len(k1)
				}
				if -1 < n1 {
					k1 = k1[0:n1]
					x09 := strings.ToLower(k1)
					if x09 != "utf-8" && strings.HasPrefix(x09, "gb") {
						//r.Body = r.gbk2utf8(r.Body)
						bDo = true
						contentType = "text/plain; charset=" + k1
					}
				}
			}
		}
		if !bDo {
			d := chardet.NewTextDetector()
			r, err := d.DetectBest(r.Body)
			if err != nil {
				return err
			}
			contentType = "text/plain; charset=" + r.Charset
		}
	}
	if strings.Contains(contentType, "utf-8") || strings.Contains(contentType, "utf8") {
		return nil
	}
	tmpBody, err := encodeBytes(r.Body, contentType)
	if err != nil {
		return err
	}
	r.Body = tmpBody
	return nil
}

func (r *Response) gbk2utf8(data []byte) []byte {
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.HZGB2312.NewDecoder())
	if d, e := ioutil.ReadAll(reader); nil == e {
		return d
	}
	return data
}
func encodeBytes(b []byte, contentType string) ([]byte, error) {
	r, err := charset.NewReader(bytes.NewReader(b), contentType)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}
