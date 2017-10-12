// Package colly implements a HTTP scraping framework
package colly

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"

	"github.com/PuerkitoBio/goquery"
)

// Collector provides the scraper instance for a scraping job
type Collector struct {
	// UserAgent is the User-Agent string used by HTTP requests
	UserAgent string
	// MaxDepth limits the recursion depth of visited URLs.
	// Set it to 0 for infinite recursion (default).
	MaxDepth int
	// AllowedDomains is a domain whitelist.
	// Leave it blank to allow any domains to be visited
	AllowedDomains []string
	// AllowURLRevisit allows multiple downloads of the same URL
	AllowURLRevisit bool
	// MaxBodySize is the limit of the retrieved response body in bytes.
	// `0` means unlimited.
	// The default value for MaxBodySize is 10MB (10 * 1024 * 1024 bytes).
	MaxBodySize int
	// CacheDir specifies a location where GET requests are cached as files.
	// When it's not defined, caching is disabled.
	CacheDir          string
	visitedURLs       []string
	htmlCallbacks     map[string]HTMLCallback
	requestCallbacks  []RequestCallback
	responseCallbacks []ResponseCallback
	backend           *httpBackend
	wg                *sync.WaitGroup
	lock              *sync.Mutex
}

// Request is the representation of a HTTP request made by a Collector
type Request struct {
	// URL is the parsed URL of the HTTP request
	URL *url.URL
	// Headers contains the Request's HTTP headers
	Headers *http.Header
	// Ctx is a context between a Request and a Response
	Ctx *Context
	// Depth is the number of the parents of this request
	Depth     int
	collector *Collector
}

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

// HTMLElement is the representation of a HTML tag.
type HTMLElement struct {
	// Name is the name of the tag
	Name       string
	Text       string
	attributes []html.Attribute
	// Request is the request object of the element's HTML document
	Request *Request
	// Response is the Response object of the element's HTML document
	Response *Response
	// DOM is the goquery parsed DOM object of the page. DOM is relative
	// to the current HTMLElement
	DOM *goquery.Selection
}

// Context provides a tiny layer for passing data between callbacks
type Context struct {
	contextMap map[string]string
	lock       *sync.RWMutex
}

// RequestCallback is a type alias for OnRequest callback functions
type RequestCallback func(*Request)

// ResponseCallback is a type alias for OnResponse callback functions
type ResponseCallback func(*Response)

// HTMLCallback is a type alias for OnHTML callback functions
type HTMLCallback func(*HTMLElement)

// NewCollector creates a new Collector instance with default configuration
func NewCollector() *Collector {
	c := &Collector{}
	c.Init()
	return c
}

// NewContext initializes a new Context instance
func NewContext() *Context {
	return &Context{
		contextMap: make(map[string]string),
		lock:       &sync.RWMutex{},
	}
}

// Init initializes the Collector's private variables and sets default
// configuration for the Collector
func (c *Collector) Init() {
	c.UserAgent = "colly - https://github.com/asciimoo/colly"
	c.MaxDepth = 0
	c.visitedURLs = make([]string, 0, 8)
	c.htmlCallbacks = make(map[string]HTMLCallback, 0)
	c.requestCallbacks = make([]RequestCallback, 0, 8)
	c.responseCallbacks = make([]ResponseCallback, 0, 8)
	c.MaxBodySize = 10 * 1024 * 1024
	c.backend = &httpBackend{}
	c.backend.Init()
	c.wg = &sync.WaitGroup{}
	c.lock = &sync.Mutex{}
}

// Visit starts Collector's collecting job by creating a
// request to the URL specified in parameter.
// Visit also calls the previously provided OnRequest,
// OnResponse, OnHTML callbacks
func (c *Collector) Visit(URL string) error {
	return c.scrape(URL, "GET", 1, nil, nil)
}

// Post starts a collector job by creating a POST request.
// Post also calls the previously provided OnRequest,
// OnResponse, OnHTML callbacks
func (c *Collector) Post(URL string, requestData map[string]string) error {
	return c.scrape(URL, "POST", 1, createFormReader(requestData), nil)
}

// PostRaw starts a collector job by creating a POST request with raw binary data.
// Post also calls the previously provided callbacks
func (c *Collector) PostRaw(URL string, requestData []byte) error {
	return c.scrape(URL, "POST", 1, bytes.NewReader(requestData), nil)
}

func (c *Collector) scrape(u, method string, depth int, requestData io.Reader, ctx *Context) error {
	c.wg.Add(1)
	defer c.wg.Done()
	if u == "" {
		return errors.New("Missing URL")
	}
	if c.MaxDepth > 0 && c.MaxDepth < depth {
		return errors.New("Max depth limit reached")
	}
	if !c.AllowURLRevisit {
		visited := false
		for _, u2 := range c.visitedURLs {
			if u2 == u {
				visited = true
				break
			}
		}
		if visited {
			return errors.New("URL already visited")
		}
	}
	parsedURL, err := url.Parse(u)
	if err != nil {
		return err
	}
	allowed := false
	if c.AllowedDomains == nil || len(c.AllowedDomains) == 0 {
		allowed = true
	} else {
		for _, d := range c.AllowedDomains {
			if d == parsedURL.Host {
				allowed = true
				break
			}
		}
	}
	if !allowed {
		return errors.New("Forbidden domain")
	}
	if !c.AllowURLRevisit {
		c.lock.Lock()
		c.visitedURLs = append(c.visitedURLs, u)
		c.lock.Unlock()
	}
	req, err := http.NewRequest(method, u, requestData)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", c.UserAgent)
	if method == "POST" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	if ctx == nil {
		ctx = NewContext()
	}
	request := &Request{
		URL:       parsedURL,
		Headers:   &req.Header,
		Ctx:       ctx,
		Depth:     depth,
		collector: c,
	}
	if len(c.requestCallbacks) > 0 {
		c.handleOnRequest(request)
	}
	response, err := c.backend.Cache(req, c.MaxBodySize, c.CacheDir)
	// TODO add OnError callback to handle these cases
	if err != nil {
		return err
	}
	response.Ctx = ctx
	response.Request = request
	response.fixCharset()
	if len(c.responseCallbacks) > 0 {
		c.handleOnResponse(response)
	}
	if strings.Index(strings.ToLower(response.Headers.Get("Content-Type")), "html") > -1 {
		c.handleOnHTML(request, response)
	}
	return nil
}

// Wait returns when the collector jobs are finished
func (c *Collector) Wait() {
	c.wg.Wait()
}

// OnRequest registers a function. Function will be executed on every
// request made by the Collector
func (c *Collector) OnRequest(f RequestCallback) {
	c.lock.Lock()
	c.requestCallbacks = append(c.requestCallbacks, f)
	c.lock.Unlock()
}

// OnResponse registers a function. Function will be executed on every response
func (c *Collector) OnResponse(f ResponseCallback) {
	c.lock.Lock()
	c.responseCallbacks = append(c.responseCallbacks, f)
	c.lock.Unlock()
}

// OnHTML registers a function. Function will be executed on every HTML
// element matched by the `goquerySelector` parameter.
// `goquerySelector` is a selector used by https://github.com/PuerkitoBio/goquery
func (c *Collector) OnHTML(goquerySelector string, f HTMLCallback) {
	c.lock.Lock()
	c.htmlCallbacks[goquerySelector] = f
	c.lock.Unlock()
}

// WithTransport allows you to set a custom http.Transport for this collector.
func (c *Collector) WithTransport(transport *http.Transport) {
	c.backend.Client.Transport = transport
}

// DisableCookies turns off cookie handling for this collector
func (c *Collector) DisableCookies() {
	c.backend.Client.Jar = nil
}

// SetRequestTimeout overrides the default timeout (10 seconds) for this collector
func (c *Collector) SetRequestTimeout(timeout time.Duration) {
	c.backend.Client.Timeout = timeout
}

func (c *Collector) handleOnRequest(r *Request) {
	for _, f := range c.requestCallbacks {
		f(r)
	}
}

func (c *Collector) handleOnResponse(r *Response) {
	for _, f := range c.responseCallbacks {
		f(r)
	}
}

func (c *Collector) handleOnHTML(req *Request, resp *Response) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(resp.Body))
	if err != nil {
		return
	}
	for expr, f := range c.htmlCallbacks {
		doc.Find(expr).Each(func(i int, s *goquery.Selection) {
			for _, n := range s.Nodes {
				e := &HTMLElement{
					Name:       n.Data,
					Request:    req,
					Response:   resp,
					Text:       goquery.NewDocumentFromNode(n).Text(),
					DOM:        s,
					attributes: n.Attr,
				}
				f(e)
			}
		})
	}
}

// Limit adds a new `LimitRule` to the collector
func (c *Collector) Limit(rule *LimitRule) error {
	return c.backend.Limit(rule)
}

// Limits adds new `LimitRule`s to the collector
func (c *Collector) Limits(rules []*LimitRule) error {
	return c.backend.Limits(rules)
}

// SetCookies handles the receipt of the cookies in a reply for the given URL
func (c *Collector) SetCookies(URL string, cookies []*http.Cookie) error {
	if c.backend.Client.Jar == nil {
		return errors.New("Cookie jar is not available")
	}
	u, err := url.Parse(URL)
	if err != nil {
		return err
	}
	c.backend.Client.Jar.SetCookies(u, cookies)
	return nil
}

// Cookies returns the cookies to send in a request for the given URL.
func (c *Collector) Cookies(URL string) []*http.Cookie {
	if c.backend.Client.Jar == nil {
		return nil
	}
	u, err := url.Parse(URL)
	if err != nil {
		return nil
	}
	return c.backend.Client.Jar.Cookies(u)
}

// Attr returns the selected attribute of a HTMLElement or empty string
// if no attribute found
func (h *HTMLElement) Attr(k string) string {
	for _, a := range h.attributes {
		if a.Key == k {
			return a.Val
		}
	}
	return ""
}

// AbsoluteURL returns with the resolved absolute URL of an URL chunk.
// AbsoluteURL returns empty string if the URL chunk is a fragment or
// could not be parsed
func (r *Request) AbsoluteURL(u string) string {
	if strings.HasPrefix(u, "#") {
		return ""
	}
	absURL, err := r.URL.Parse(u)
	if err != nil {
		return ""
	}
	absURL.Fragment = ""
	if absURL.Scheme == "//" {
		absURL.Scheme = r.URL.Scheme
	}
	return absURL.String()
}

// Visit continues Collector's collecting job by creating a
// request and preserves the Context of the previous request.
// Visit also calls the previously provided OnRequest,
// OnResponse, OnHTML callbacks
func (r *Request) Visit(URL string) error {
	return r.collector.scrape(r.AbsoluteURL(URL), "GET", r.Depth+1, nil, r.Ctx)
}

// Post continues a collector job by creating a POST request and preserves the Context
// of the previous request.
// Post also calls the previously provided OnRequest, OnResponse, OnHTML callbacks
func (r *Request) Post(URL string, requestData map[string]string) error {
	return r.collector.scrape(r.AbsoluteURL(URL), "POST", r.Depth+1, createFormReader(requestData), r.Ctx)
}

// PostRaw starts a collector job by creating a POST request with raw binary data.
// PostRaw preserves the Context of the previous request
// and calls the previously provided callbacks
func (r *Request) PostRaw(URL string, requestData []byte) error {
	return r.collector.scrape(r.AbsoluteURL(URL), "POST", r.Depth+1, bytes.NewReader(requestData), r.Ctx)
}

// UnmarshalBinary decodes Context value to nil
// This function is used by request caching
func (c *Context) UnmarshalBinary(_ []byte) error {
	return nil
}

// MarshalBinary encodes Context value
// This function is used by request caching
func (c *Context) MarshalBinary() (_ []byte, _ error) {
	return nil, nil
}

// Put stores a value in Context
func (c *Context) Put(key, value string) {
	c.lock.Lock()
	c.contextMap[key] = value
	c.lock.Unlock()
}

// Get retrieves a value from Context. If no value found for `k`
// Get returns an empty string if key not found
func (c *Context) Get(key string) string {
	c.lock.RLock()
	if v, ok := c.contextMap[key]; ok {
		c.lock.RUnlock()
		return v
	}
	c.lock.RUnlock()
	return ""
}

func createFormReader(data map[string]string) io.Reader {
	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}
	return strings.NewReader(form.Encode())
}

func (r *Response) fixCharset() {
	contentType := strings.ToLower(r.Headers.Get("Content-Type"))
	if strings.Index(contentType, "charset") == -1 {
		return
	}
	if strings.Index(contentType, "utf-8") != -1 || strings.Index(contentType, "utf8") != -1 {
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
