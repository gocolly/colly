// Package extensions implements various helper addons for Colly
package extensions

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/gocolly/colly/v2"
)

// DynamicContentOptions contains options for the DynamicContent extension
type DynamicContentOptions struct {
	// Timeout is the maximum time to wait for the page to load
	Timeout time.Duration
	// WaitForSelector is a CSS selector to wait for before considering the page loaded
	WaitForSelector string
	// UserAgent is the user agent to use for the headless browser
	UserAgent string
	// Headless determines whether to run the browser in headless mode
	Headless bool
	// CustomBrowserPath is the path to a custom browser executable
	CustomBrowserPath string
	// ExtraHeaders are additional headers to send with each request
	ExtraHeaders map[string]string
}

// DefaultDynamicContentOptions returns the default options for DynamicContent
func DefaultDynamicContentOptions() *DynamicContentOptions {
	return &DynamicContentOptions{
		Timeout:         30 * time.Second,
		WaitForSelector: "body",
		UserAgent:       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
		Headless:        true,
		ExtraHeaders:    make(map[string]string),
	}
}

// dynamicContentTransport is a custom http.RoundTripper that uses Rod to render JavaScript
type dynamicContentTransport struct {
	browser       *rod.Browser
	options       *DynamicContentOptions
	nextTransport http.RoundTripper
}

// RoundTrip implements the http.RoundTripper interface
func (t *dynamicContentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Only process GET requests
	if req.Method != "GET" {
		return t.nextTransport.RoundTrip(req)
	}

	// Create a new page
	page := t.browser.MustPage()
	defer page.Close()

	// Set user agent
	if t.options.UserAgent != "" {
		page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent: t.options.UserAgent,
		})
	}

	// Set extra headers
	if len(t.options.ExtraHeaders) > 0 {
		extraHeaders := []string{}
		for name, value := range t.options.ExtraHeaders {
			extraHeaders = append(extraHeaders, name, value)
		}
		page.MustSetExtraHeaders(extraHeaders...)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), t.options.Timeout)
	defer cancel()

	// Navigate to the URL
	err := page.Context(ctx).Navigate(req.URL.String())
	if err != nil {
		// Fall back to the original transport if navigation fails
		return t.nextTransport.RoundTrip(req)
	}

	// Wait for the page to load
	if t.options.WaitForSelector != "" {
		err = page.Context(ctx).WaitElementsMoreThan(t.options.WaitForSelector, 0)
		if err != nil {
			// Fall back to the original transport if waiting fails
			return t.nextTransport.RoundTrip(req)
		}
	}

	// Get the HTML content
	html, err := page.HTML()
	if err != nil {
		// Fall back to the original transport if getting HTML fails
		return t.nextTransport.RoundTrip(req)
	}

	// Create a response with the rendered HTML
	htmlBytes := []byte(html)
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(htmlBytes)),
		Header:     make(http.Header),
		Request:    req,
	}

	// Set content type header
	resp.Header.Set("Content-Type", "text/html; charset=utf-8")
	resp.ContentLength = int64(len(htmlBytes))

	return resp, nil
}

// DynamicContent enables JavaScript rendering for Colly using the Rod headless browser
func DynamicContent(c *colly.Collector, options *DynamicContentOptions) {
	if options == nil {
		options = DefaultDynamicContentOptions()
	}

	// Initialize browser launcher
	var browser *rod.Browser

	// Create a new browser instance
	launcherURL := launcher.New().
		Headless(options.Headless).
		Set("disable-web-security", "true").
		Set("disable-setuid-sandbox", "true").
		Set("no-sandbox", "true")

	if options.CustomBrowserPath != "" {
		launcherURL = launcherURL.Bin(options.CustomBrowserPath)
	}

	browserURL := launcherURL.MustLaunch()
	browser = rod.New().ControlURL(browserURL).MustConnect()

	// Create a transport that will intercept requests and use Rod to render them
	originalTransport := http.DefaultTransport
	c.WithTransport(&dynamicContentTransport{
		browser:       browser,
		options:       options,
		nextTransport: originalTransport,
	})
}
