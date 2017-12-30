package colly

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"

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
	return SanitizeFileName(r.Request.URL.Path[1:])
}

func (r *Response) fixCharset() {
	contentType := strings.ToLower(r.Headers.Get("Content-Type"))
	if !strings.Contains(contentType, "charset") {
		return
	}
	if strings.Contains(contentType, "utf-8") || strings.Contains(contentType, "utf8") {
		return
	}
	encodedBodyReader, err := charset.NewReader(bytes.NewReader(r.Body), contentType)
	if err != nil {
		return
	}
	tmpBody, err := ioutil.ReadAll(encodedBodyReader)
	if err != nil {
		return
	}
	r.Body = tmpBody
}
