package colly

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"testing"
)

var testServerPort int = 31337
var testServerAddr string = fmt.Sprintf("127.0.0.1:%d", testServerPort)
var testServerRootURL string = fmt.Sprintf("http://%s/", testServerAddr)
var serverIndexResponse []byte = []byte("hello world\n")
var robotsFile string = `
User-agent: *
Allow: /allowed
Disallow: /disallowed
`

func init() {
	srv := &http.Server{}
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: testServerPort})
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(serverIndexResponse)
	})

	http.HandleFunc("/html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Conent-Type", "text/html")
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

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			w.Header().Set("Conent-Type", "text/html")
			w.Write([]byte(r.FormValue("name")))
		}
	})

	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(robotsFile))
	})

	http.HandleFunc("/allowed", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("allowed"))
	})

	http.HandleFunc("/disallowed", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("disallowed"))
	})

	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()
}

func TestCollectorVisit(t *testing.T) {
	c := NewCollector()

	onRequestCalled := false
	onResponseCalled := false

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

	c.Visit(testServerRootURL)

	if !onRequestCalled {
		t.Error("Failed to call OnRequest callback")
	}

	if !onResponseCalled {
		t.Error("Failed to call OnResponse callback")
	}
}

func TestCollectorOnHTML(t *testing.T) {
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
		paragraphCallbackCount += 1
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

	c.Visit(testServerRootURL + "/html")

	if !titleCallbackCalled {
		t.Error("Failed to call OnHTML callback for <title> tag")
	}

	if paragraphCallbackCount != 2 {
		t.Error("Failed to find all <p> tags")
	}
}

func TestCollectorPost(t *testing.T) {
	postValue := "hello"
	c := NewCollector()

	c.OnResponse(func(r *Response) {
		if postValue != string(r.Body) {
			t.Error("Failed to send data with POST")
		}
	})

	c.Post(testServerRootURL+"login", map[string]string{
		"name": postValue,
	})
}

func BenchmarkVisit(b *testing.B) {
	c := NewCollector()
	c.OnHTML("p", func(_ *HTMLElement) {})

	for n := 0; n < b.N; n++ {
		c.Visit(fmt.Sprintf("%shtml?q=%d", testServerRootURL, n))
	}
}

func TestRobotsWhenAllowed(t *testing.T) {
	c := NewCollector()
	c.IgnoreRobotsTxt = false

	c.OnResponse(func(resp *Response) {
		if resp.StatusCode != 200 {
			t.Fatalf("Wrong response code: %d", resp.StatusCode)
		}
	})

	err := c.Visit(testServerRootURL + "allowed")

	if err != nil {
		t.Fatal(err)
	}
}

func TestRobotsWhenDisallowed(t *testing.T) {
	c := NewCollector()
	c.IgnoreRobotsTxt = false

	c.OnResponse(func(resp *Response) {
		t.Fatalf("Received response: %d", resp.StatusCode)
	})

	err := c.Visit(testServerRootURL + "disallowed")
	if err.Error() != "URL blocked by robots.txt" {
		t.Fatalf("wrong error message: %v", err)
	}
}

func TestIgnoreRobotsWhenDisallowed(t *testing.T) {
	c := NewCollector()
	c.IgnoreRobotsTxt = true

	c.OnResponse(func(resp *Response) {
		if resp.StatusCode != 200 {
			t.Fatalf("Wrong response code: %d", resp.StatusCode)
		}
	})

	err := c.Visit(testServerRootURL + "disallowed")

	if err != nil {
		t.Fatal(err)
	}

}
