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
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"github.com/gocolly/colly/v2/debug"
)

var serverIndexResponse = []byte("hello world\n")
var robotsFile = `
User-agent: *
Allow: /allowed
Disallow: /disallowed
Disallow: /allowed*q=
`

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write(serverIndexResponse)
	})

	mux.HandleFunc("/html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<title>Test Page</title>
</head>
<body>
<h1>Hello World</h1>
<p class="description">This is a test page</p>
<p class="description">This is a test paragraph</p>
</body>
</html>
		`))
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(r.FormValue("name")))
		}
	})

	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(robotsFile))
	})

	mux.HandleFunc("/allowed", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("allowed"))
	})

	mux.HandleFunc("/disallowed", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("disallowed"))
	})

	mux.Handle("/redirect", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/redirected/", http.StatusSeeOther)

	}))

	mux.Handle("/redirected/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<a href="test">test</a>`)
	}))

	mux.HandleFunc("/set_cookie", func(w http.ResponseWriter, r *http.Request) {
		c := &http.Cookie{Name: "test", Value: "testv", HttpOnly: false}
		http.SetCookie(w, c)
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/check_cookie", func(w http.ResponseWriter, r *http.Request) {
		cs := r.Cookies()
		if len(cs) != 1 || r.Cookies()[0].Value != "testv" {
			w.WriteHeader(500)
			w.Write([]byte("nok"))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(500)
		w.Write([]byte("<p>error</p>"))
	})

	mux.HandleFunc("/user_agent", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(r.Header.Get("User-Agent")))
	})

	mux.HandleFunc("/base", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<title>Test Page</title>
<base href="http://xy.com/" />
</head>
<body>
<a href="z">link</a>
</body>
</html>
		`))
	})

	mux.HandleFunc("/large_binary", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		ww := bufio.NewWriter(w)
		defer ww.Flush()
		for {
			// have to check error to detect client aborting download
			if _, err := ww.Write([]byte{0x41}); err != nil {
				return
			}
		}
	})

	return httptest.NewServer(mux)
}

var newCollectorTests = map[string]func(*testing.T){
	"UserAgent": func(t *testing.T) {
		for _, ua := range []string{
			"foo",
			"bar",
		} {
			c := NewCollector(UserAgent(ua))

			if got, want := c.UserAgent, ua; got != want {
				t.Fatalf("c.UserAgent = %q, want %q", got, want)
			}
		}
	},
	"MaxDepth": func(t *testing.T) {
		for _, depth := range []int{
			12,
			34,
			0,
		} {
			c := NewCollector(MaxDepth(depth))

			if got, want := c.MaxDepth, depth; got != want {
				t.Fatalf("c.MaxDepth = %d, want %d", got, want)
			}
		}
	},
	"AllowedDomains": func(t *testing.T) {
		for _, domains := range [][]string{
			{"example.com", "example.net"},
			{"example.net"},
			{},
			nil,
		} {
			c := NewCollector(AllowedDomains(domains...))

			if got, want := c.AllowedDomains, domains; !reflect.DeepEqual(got, want) {
				t.Fatalf("c.AllowedDomains = %q, want %q", got, want)
			}
		}
	},
	"DisallowedDomains": func(t *testing.T) {
		for _, domains := range [][]string{
			{"example.com", "example.net"},
			{"example.net"},
			{},
			nil,
		} {
			c := NewCollector(DisallowedDomains(domains...))

			if got, want := c.DisallowedDomains, domains; !reflect.DeepEqual(got, want) {
				t.Fatalf("c.DisallowedDomains = %q, want %q", got, want)
			}
		}
	},
	"DisallowedURLFilters": func(t *testing.T) {
		for _, filters := range [][]*regexp.Regexp{
			{regexp.MustCompile(`.*not_allowed.*`)},
		} {
			c := NewCollector(DisallowedURLFilters(filters...))

			if got, want := c.DisallowedURLFilters, filters; !reflect.DeepEqual(got, want) {
				t.Fatalf("c.DisallowedURLFilters = %v, want %v", got, want)
			}
		}
	},
	"URLFilters": func(t *testing.T) {
		for _, filters := range [][]*regexp.Regexp{
			{regexp.MustCompile(`\w+`)},
			{regexp.MustCompile(`\d+`)},
			{},
			nil,
		} {
			c := NewCollector(URLFilters(filters...))

			if got, want := c.URLFilters, filters; !reflect.DeepEqual(got, want) {
				t.Fatalf("c.URLFilters = %v, want %v", got, want)
			}
		}
	},
	"AllowURLRevisit": func(t *testing.T) {
		c := NewCollector(AllowURLRevisit())

		if !c.AllowURLRevisit {
			t.Fatal("c.AllowURLRevisit = false, want true")
		}
	},
	"MaxBodySize": func(t *testing.T) {
		for _, sizeInBytes := range []int{
			1024 * 1024,
			1024,
			0,
		} {
			c := NewCollector(MaxBodySize(sizeInBytes))

			if got, want := c.MaxBodySize, sizeInBytes; got != want {
				t.Fatalf("c.MaxBodySize = %d, want %d", got, want)
			}
		}
	},
	"CacheDir": func(t *testing.T) {
		for _, path := range []string{
			"/tmp/",
			"/var/cache/",
		} {
			c := NewCollector(CacheDir(path))

			if got, want := c.CacheDir, path; got != want {
				t.Fatalf("c.CacheDir = %q, want %q", got, want)
			}
		}
	},
	"IgnoreRobotsTxt": func(t *testing.T) {
		c := NewCollector(IgnoreRobotsTxt())

		if !c.IgnoreRobotsTxt {
			t.Fatal("c.IgnoreRobotsTxt = false, want true")
		}
	},
	"ID": func(t *testing.T) {
		for _, id := range []uint32{
			0,
			1,
			2,
		} {
			c := NewCollector(ID(id))

			if got, want := c.ID, id; got != want {
				t.Fatalf("c.ID = %d, want %d", got, want)
			}
		}
	},
	"DetectCharset": func(t *testing.T) {
		c := NewCollector(DetectCharset())

		if !c.DetectCharset {
			t.Fatal("c.DetectCharset = false, want true")
		}
	},
	"Debugger": func(t *testing.T) {
		d := &debug.LogDebugger{}
		c := NewCollector(Debugger(d))

		if got, want := c.debugger, d; got != want {
			t.Fatalf("c.debugger = %v, want %v", got, want)
		}
	},
	"CheckHead": func(t *testing.T) {
		c := NewCollector(CheckHead())

		if !c.CheckHead {
			t.Fatal("c.CheckHead = false, want true")
		}
	},
}

func TestNewCollector(t *testing.T) {
	t.Run("Functional Options", func(t *testing.T) {
		for name, test := range newCollectorTests {
			t.Run(name, test)
		}
	})
}

func TestCollectorVisit(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector()

	onRequestCalled := false
	onResponseCalled := false
	onScrapedCalled := false

	c.OnRequest(func(r *Request) {
		onRequestCalled = true
		r.Ctx.Put("x", "y")
	})

	c.OnResponse(func(r *Response) {
		onResponseCalled = true

		if r.Ctx.Get("x") != "y" {
			t.Error("Failed to retrieve context value for key 'x'")
		}

		if !bytes.Equal(r.Body, serverIndexResponse) {
			t.Error("Response body does not match with the original content")
		}
	})

	c.OnScraped(func(r *Response) {
		if !onResponseCalled {
			t.Error("OnScraped called before OnResponse")
		}

		if !onRequestCalled {
			t.Error("OnScraped called before OnRequest")
		}

		onScrapedCalled = true
	})

	c.Visit(ts.URL)

	if !onRequestCalled {
		t.Error("Failed to call OnRequest callback")
	}

	if !onResponseCalled {
		t.Error("Failed to call OnResponse callback")
	}

	if !onScrapedCalled {
		t.Error("Failed to call OnScraped callback")
	}
}

func TestCollectorVisitWithAllowedDomains(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector(AllowedDomains("localhost", "127.0.0.1", "::1"))
	err := c.Visit(ts.URL)
	if err != nil {
		t.Errorf("Failed to visit url %s", ts.URL)
	}

	err = c.Visit("http://example.com")
	if err != ErrForbiddenDomain {
		t.Errorf("c.Visit should return ErrForbiddenDomain, but got %v", err)
	}
}

func TestCollectorVisitWithDisallowedDomains(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector(DisallowedDomains("localhost", "127.0.0.1", "::1"))
	err := c.Visit(ts.URL)
	if err != ErrForbiddenDomain {
		t.Errorf("c.Visit should return ErrForbiddenDomain, but got %v", err)
	}

	c2 := NewCollector(DisallowedDomains("example.com"))
	err = c2.Visit("http://example.com:8080")
	if err != ErrForbiddenDomain {
		t.Errorf("c.Visit should return ErrForbiddenDomain, but got %v", err)
	}
	err = c2.Visit(ts.URL)
	if err != nil {
		t.Errorf("Failed to visit url %s", ts.URL)
	}
}

func TestCollectorVisitResponseHeaders(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	var onResponseHeadersCalled bool

	c := NewCollector()
	c.OnResponseHeaders(func(r *Response) {
		onResponseHeadersCalled = true
		if r.Headers.Get("Content-Type") == "application/octet-stream" {
			r.Request.Abort()
		}
	})
	c.OnResponse(func(r *Response) {
		t.Error("OnResponse was called")
	})
	c.Visit(ts.URL + "/large_binary")
	if !onResponseHeadersCalled {
		t.Error("OnResponseHeaders was not called")
	}
}

func TestCollectorOnHTML(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector()

	titleCallbackCalled := false
	paragraphCallbackCount := 0

	c.OnHTML("title", func(e *HTMLElement) {
		titleCallbackCalled = true
		if e.Text != "Test Page" {
			t.Error("Title element text does not match, got", e.Text)
		}
	})

	c.OnHTML("p", func(e *HTMLElement) {
		paragraphCallbackCount++
		if e.Attr("class") != "description" {
			t.Error("Failed to get paragraph's class attribute")
		}
	})

	c.OnHTML("body", func(e *HTMLElement) {
		if e.ChildAttr("p", "class") != "description" {
			t.Error("Invalid class value")
		}
		classes := e.ChildAttrs("p", "class")
		if len(classes) != 2 {
			t.Error("Invalid class values")
		}
	})

	c.Visit(ts.URL + "/html")

	if !titleCallbackCalled {
		t.Error("Failed to call OnHTML callback for <title> tag")
	}

	if paragraphCallbackCount != 2 {
		t.Error("Failed to find all <p> tags")
	}
}

func TestCollectorURLRevisit(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector()

	visitCount := 0

	c.OnRequest(func(r *Request) {
		visitCount++
	})

	c.Visit(ts.URL)
	c.Visit(ts.URL)

	if visitCount != 1 {
		t.Error("URL revisited")
	}

	c.AllowURLRevisit = true

	c.Visit(ts.URL)
	c.Visit(ts.URL)

	if visitCount != 3 {
		t.Error("URL not revisited")
	}
}

func TestCollectorURLRevisitCheck(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector()

	visited, err := c.HasVisited(ts.URL)

	if err != nil {
		t.Error(err.Error())
	}

	if visited != false {
		t.Error("Expected URL to NOT have been visited")
	}

	c.Visit(ts.URL)

	visited, err = c.HasVisited(ts.URL)

	if err != nil {
		t.Error(err.Error())
	}

	if visited != true {
		t.Error("Expected URL to have been visited")
	}
}

// TestCollectorURLRevisitDisallowed ensures that disallowed URL is not considered visited.
func TestCollectorURLRevisitDomainDisallowed(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	parsedURL, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	c := NewCollector(DisallowedDomains(parsedURL.Hostname()))
	err = c.Visit(ts.URL)
	if got, want := err, ErrForbiddenDomain; got != want {
		t.Fatalf("wrong error on first visit: got=%v want=%v", got, want)
	}
	err = c.Visit(ts.URL)
	if got, want := err, ErrForbiddenDomain; got != want {
		t.Fatalf("wrong error on second visit: got=%v want=%v", got, want)
	}

}

func TestCollectorPost(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	postValue := "hello"
	c := NewCollector()

	c.OnResponse(func(r *Response) {
		if postValue != string(r.Body) {
			t.Error("Failed to send data with POST")
		}
	})

	c.Post(ts.URL+"/login", map[string]string{
		"name": postValue,
	})
}

func TestRedirect(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector()
	c.OnHTML("a[href]", func(e *HTMLElement) {
		u := e.Request.AbsoluteURL(e.Attr("href"))
		if !strings.HasSuffix(u, "/redirected/test") {
			t.Error("Invalid URL after redirect: " + u)
		}
	})
	c.OnResponse(func(r *Response) {
		if !strings.HasSuffix(r.Request.URL.String(), "/redirected/") {
			t.Error("Invalid URL in Request after redirect: " + r.Request.URL.String())
		}
	})
	c.Visit(ts.URL + "/redirect")
}

func TestBaseTag(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector()
	c.OnHTML("a[href]", func(e *HTMLElement) {
		u := e.Request.AbsoluteURL(e.Attr("href"))
		if u != "http://xy.com/z" {
			t.Error("Invalid <base /> tag handling in OnHTML: expected https://xy.com/z, got " + u)
		}
	})
	c.Visit(ts.URL + "/base")

	c2 := NewCollector()
	c2.OnXML("//a", func(e *XMLElement) {
		u := e.Request.AbsoluteURL(e.Attr("href"))
		if u != "http://xy.com/z" {
			t.Error("Invalid <base /> tag handling in OnXML: expected https://xy.com/z, got " + u)
		}
	})
	c2.Visit(ts.URL + "/base")
}

func TestCollectorCookies(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector()

	if err := c.Visit(ts.URL + "/set_cookie"); err != nil {
		t.Fatal(err)
	}

	if err := c.Visit(ts.URL + "/check_cookie"); err != nil {
		t.Fatalf("Failed to use previously set cookies: %s", err)
	}
}

func TestRobotsWhenAllowed(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector()
	c.IgnoreRobotsTxt = false

	c.OnResponse(func(resp *Response) {
		if resp.StatusCode != 200 {
			t.Fatalf("Wrong response code: %d", resp.StatusCode)
		}
	})

	err := c.Visit(ts.URL + "/allowed")

	if err != nil {
		t.Fatal(err)
	}
}

func TestRobotsWhenDisallowed(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector()
	c.IgnoreRobotsTxt = false

	c.OnResponse(func(resp *Response) {
		t.Fatalf("Received response: %d", resp.StatusCode)
	})

	err := c.Visit(ts.URL + "/disallowed")
	if err.Error() != "URL blocked by robots.txt" {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestRobotsWhenDisallowedWithQueryParameter(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector()
	c.IgnoreRobotsTxt = false

	c.OnResponse(func(resp *Response) {
		t.Fatalf("Received response: %d", resp.StatusCode)
	})

	err := c.Visit(ts.URL + "/allowed?q=1")
	if err.Error() != "URL blocked by robots.txt" {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestIgnoreRobotsWhenDisallowed(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector()
	c.IgnoreRobotsTxt = true

	c.OnResponse(func(resp *Response) {
		if resp.StatusCode != 200 {
			t.Fatalf("Wrong response code: %d", resp.StatusCode)
		}
	})

	err := c.Visit(ts.URL + "/disallowed")

	if err != nil {
		t.Fatal(err)
	}

}

func TestConnectionErrorOnRobotsTxtResultsInError(t *testing.T) {
	ts := newTestServer()
	ts.Close() // immediately close the server to force a connection error

	c := NewCollector()
	c.IgnoreRobotsTxt = false
	err := c.Visit(ts.URL)

	if err == nil {
		t.Fatal("Error expected")
	}
}

func TestEnvSettings(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	os.Setenv("COLLY_USER_AGENT", "test")
	defer os.Unsetenv("COLLY_USER_AGENT")

	c := NewCollector()

	valid := false

	c.OnResponse(func(resp *Response) {
		if string(resp.Body) == "test" {
			valid = true
		}
	})

	c.Visit(ts.URL + "/user_agent")

	if !valid {
		t.Fatalf("Wrong user-agent from environment")
	}
}

func TestParseHTTPErrorResponse(t *testing.T) {
	contentCount := 0
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector(
		AllowURLRevisit(),
	)

	c.OnHTML("p", func(e *HTMLElement) {
		if e.Text == "error" {
			contentCount++
		}
	})

	c.Visit(ts.URL + "/500")

	if contentCount != 0 {
		t.Fatal("Content is parsed without ParseHTTPErrorResponse enabled")
	}

	c.ParseHTTPErrorResponse = true

	c.Visit(ts.URL + "/500")

	if contentCount != 1 {
		t.Fatal("Content isn't parsed with ParseHTTPErrorResponse enabled")
	}

}

func TestHTMLElement(t *testing.T) {
	ctx := &Context{}
	resp := &Response{
		Request: &Request{
			Ctx: ctx,
		},
		Ctx: ctx,
	}

	in := `<a href="http://go-colly.org">Colly</a>`
	sel := "a[href]"
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer([]byte(in)))
	if err != nil {
		t.Fatal(err)
	}
	elements := []*HTMLElement{}
	i := 0
	doc.Find(sel).Each(func(_ int, s *goquery.Selection) {
		for _, n := range s.Nodes {
			elements = append(elements, NewHTMLElementFromSelectionNode(resp, s, n, i))
			i++
		}
	})
	elementsLen := len(elements)
	if elementsLen != 1 {
		t.Errorf("element length mismatch. got %d, expected %d.\n", elementsLen, 1)
	}
	v := elements[0]
	if v.Name != "a" {
		t.Errorf("element tag mismatch. got %s, expected %s.\n", v.Name, "a")
	}
	if v.Text != "Colly" {
		t.Errorf("element content mismatch. got %s, expected %s.\n", v.Text, "Colly")
	}
	if v.Attr("href") != "http://go-colly.org" {
		t.Errorf("element href mismatch. got %s, expected %s.\n", v.Attr("href"), "http://go-colly.org")
	}
}

func TestCollectorOnXML(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector()

	titleCallbackCalled := false
	paragraphCallbackCount := 0

	c.OnXML("/html/head/title", func(e *XMLElement) {
		titleCallbackCalled = true
		if e.Text != "Test Page" {
			t.Error("Title element text does not match, got", e.Text)
		}
	})

	c.OnXML("/html/body/p", func(e *XMLElement) {
		paragraphCallbackCount++
		if e.Attr("class") != "description" {
			t.Error("Failed to get paragraph's class attribute")
		}
	})

	c.OnXML("/html/body", func(e *XMLElement) {
		if e.ChildAttr("p", "class") != "description" {
			t.Error("Invalid class value")
		}
		classes := e.ChildAttrs("p", "class")
		if len(classes) != 2 {
			t.Error("Invalid class values")
		}
	})

	c.Visit(ts.URL + "/html")

	if !titleCallbackCalled {
		t.Error("Failed to call OnXML callback for <title> tag")
	}

	if paragraphCallbackCount != 2 {
		t.Error("Failed to find all <p> tags")
	}
}

func TestCollectorVisitWithTrace(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector(AllowedDomains("localhost", "127.0.0.1", "::1"), TraceHTTP())
	c.OnResponse(func(resp *Response) {
		if resp.Trace == nil {
			t.Error("Failed to initialize trace")
		}
	})

	err := c.Visit(ts.URL)
	if err != nil {
		t.Errorf("Failed to visit url %s", ts.URL)
	}
}

func TestCollectorVisitWithCheckHead(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector(CheckHead())
	var requestMethodChain []string
	c.OnResponse(func(resp *Response) {
		requestMethodChain = append(requestMethodChain, resp.Request.Method)
	})

	err := c.Visit(ts.URL)
	if err != nil {
		t.Errorf("Failed to visit url %s", ts.URL)
	}
	if requestMethodChain[0] != "HEAD" && requestMethodChain[1] != "GET" {
		t.Errorf("Failed to perform a HEAD request before GET")
	}
}

func BenchmarkOnHTML(b *testing.B) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector()
	c.OnHTML("p", func(_ *HTMLElement) {})

	for n := 0; n < b.N; n++ {
		c.Visit(fmt.Sprintf("%s/html?q=%d", ts.URL, n))
	}
}

func BenchmarkOnXML(b *testing.B) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector()
	c.OnXML("//p", func(_ *XMLElement) {})

	for n := 0; n < b.N; n++ {
		c.Visit(fmt.Sprintf("%s/html?q=%d", ts.URL, n))
	}
}

func BenchmarkOnResponse(b *testing.B) {
	ts := newTestServer()
	defer ts.Close()

	c := NewCollector()
	c.AllowURLRevisit = true
	c.OnResponse(func(_ *Response) {})

	for n := 0; n < b.N; n++ {
		c.Visit(ts.URL)
	}
}
