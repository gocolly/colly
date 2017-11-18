// Package colly implements a HTTP scraping framework
package colly

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gocolly/colly/debug"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"

	"github.com/PuerkitoBio/goquery"
	"github.com/kennygrant/sanitize"
	"github.com/temoto/robotstxt"
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
	// DisallowedDomains is a domain blacklist.
	DisallowedDomains []string
	// URLFilters is a list of regular expressions which restricts
	// visiting URLs. If any of the rules matches to a URL the
	// request won't be stopped.
	// Leave it blank to allow any URLs to be visited
	URLFilters []*regexp.Regexp
	// AllowURLRevisit allows multiple downloads of the same URL
	AllowURLRevisit bool
	// MaxBodySize is the limit of the retrieved response body in bytes.
	// 0 means unlimited.
	// The default value for MaxBodySize is 10MB (10 * 1024 * 1024 bytes).
	MaxBodySize int
	// CacheDir specifies a location where GET requests are cached as files.
	// When it's not defined, caching is disabled.
	CacheDir string
	// IgnoreRobotsTxt allows the Collector to ignore any restrictions set by
	// the target host's robots.txt file.  See http://www.robotstxt.org/ for more
	// information.
	IgnoreRobotsTxt bool
	// Id is the unique identifier of a collector
	Id                int32
	debugger          debug.Debugger
	visitedURLs       map[uint64]bool
	robotsMap         map[string]*robotstxt.RobotsData
	htmlCallbacks     []*htmlCallbackContainer
	requestCallbacks  []RequestCallback
	responseCallbacks []ResponseCallback
	errorCallbacks    []ErrorCallback
	requestCount      int32
	responseCount     int32
	backend           *httpBackend
	wg                *sync.WaitGroup
	lock              *sync.RWMutex
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
	Depth int
	// Unique identifier of the request
	Id        int32
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
	contextMap map[string]interface{}
	lock       *sync.RWMutex
}

// RequestCallback is a type alias for OnRequest callback functions
type RequestCallback func(*Request)

// ResponseCallback is a type alias for OnResponse callback functions
type ResponseCallback func(*Response)

// HTMLCallback is a type alias for OnHTML callback functions
type HTMLCallback func(*HTMLElement)

// ErrorCallback is a type alias for OnError callback functions
type ErrorCallback func(*Response, error)

type htmlCallbackContainer struct {
	Selector string
	Function HTMLCallback
}

var collectorCounter int32 = 0

// NewCollector creates a new Collector instance with default configuration
func NewCollector() *Collector {
	c := &Collector{}
	c.Init()
	return c
}

// NewContext initializes a new Context instance
func NewContext() *Context {
	return &Context{
		contextMap: make(map[string]interface{}),
		lock:       &sync.RWMutex{},
	}
}

// Init initializes the Collector's private variables and sets default
// configuration for the Collector
func (c *Collector) Init() {
	c.UserAgent = "colly - https://github.com/gocolly/colly"
	c.MaxDepth = 0
	c.visitedURLs = make(map[uint64]bool, 0)
	c.htmlCallbacks = make([]*htmlCallbackContainer, 0, 8)
	c.requestCallbacks = make([]RequestCallback, 0, 8)
	c.responseCallbacks = make([]ResponseCallback, 0, 8)
	c.errorCallbacks = make([]ErrorCallback, 0, 8)
	c.MaxBodySize = 10 * 1024 * 1024
	c.backend = &httpBackend{}
	c.backend.Init()
	c.backend.Client.CheckRedirect = c.checkRedirectFunc()
	c.wg = &sync.WaitGroup{}
	c.lock = &sync.RWMutex{}
	c.robotsMap = make(map[string]*robotstxt.RobotsData, 0)
	c.IgnoreRobotsTxt = true
	c.Id = atomic.AddInt32(&collectorCounter, 1)
}

// Appengine will replace the Collector's backend http.Client
// With an Http.Client that is provided by appengine/urlfetch
// This function should be used when the scraper is initiated
// by a http.Request to Google App Engine
func (c *Collector) Appengine(req *http.Request) {
	ctx := appengine.NewContext(req)
	client := urlfetch.Client(ctx)
	client.Jar = c.backend.Client.Jar
	client.CheckRedirect = c.backend.Client.CheckRedirect
	client.Timeout = c.backend.Client.Timeout

	c.backend.Client = client
}

// Visit starts Collector's collecting job by creating a
// request to the URL specified in parameter.
// Visit also calls the previously provided callbacks
func (c *Collector) Visit(URL string) error {
	return c.scrape(URL, "GET", 1, nil, nil, nil)
}

// Post starts a collector job by creating a POST request.
// Post also calls the previously provided callbacks
func (c *Collector) Post(URL string, requestData map[string]string) error {
	return c.scrape(URL, "POST", 1, createFormReader(requestData), nil, nil)
}

// PostRaw starts a collector job by creating a POST request with raw binary data.
// Post also calls the previously provided callbacks
func (c *Collector) PostRaw(URL string, requestData []byte) error {
	return c.scrape(URL, "POST", 1, bytes.NewReader(requestData), nil, nil)
}

// PostMultipart starts a collector job by creating a Multipart POST request
// with raw binary data.  PostMultipart also calls the previously provided callbacks
func (c *Collector) PostMultipart(URL string, requestData map[string][]byte) error {
	boundary := randomBoundary()
	hdr := http.Header{}
	hdr.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	hdr.Set("User-Agent", c.UserAgent)
	return c.scrape(URL, "POST", 1, createMultipartReader(boundary, requestData), nil, hdr)
}

// Request starts a collector job by creating a custom HTTP request
// where method, context, headers and request data can be specified.
// Set requestData, ctx, hdr parameters to nil if you don't want to use them.
// Valid methods:
//   - "GET"
//   - "POST"
//   - "PUT"
//   - "DELETE"
//   - "PATCH"
//   - "OPTIONS"
func (c *Collector) Request(method, URL string, requestData io.Reader, ctx *Context, hdr http.Header) error {
	return c.scrape(URL, method, 1, requestData, ctx, hdr)
}

// SetDebugger attaches a debugger to the collector
func (c *Collector) SetDebugger(d debug.Debugger) {
	d.Init()
	c.debugger = d
}

func (c *Collector) scrape(u, method string, depth int, requestData io.Reader, ctx *Context, hdr http.Header) error {
	c.wg.Add(1)
	defer c.wg.Done()
	if err := c.requestCheck(u, method, depth); err != nil {
		return err
	}
	parsedURL, err := url.Parse(u)
	if err != nil {
		return err
	}
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "http"
	}
	if !c.isDomainAllowed(parsedURL.Host) {
		return errors.New("Forbidden domain")
	}
	if !c.IgnoreRobotsTxt {
		if err = c.checkRobots(parsedURL); err != nil {
			return err
		}
	}
	req, err := http.NewRequest(method, parsedURL.String(), requestData)
	if err != nil {
		return err
	}
	if hdr == nil {
		req.Header.Set("User-Agent", c.UserAgent)
		if method == "POST" {
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		}
	} else {
		req.Header = hdr
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
		Id:        atomic.AddInt32(&c.requestCount, 1),
	}

	c.handleOnRequest(request)

	if method == "POST" && req.Header.Get("Content-Type") == "" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	response, err := c.backend.Cache(req, c.MaxBodySize, c.CacheDir)
	if err := c.handleOnError(response, err, request, ctx); err != nil {
		return err
	}
	atomic.AddInt32(&c.responseCount, 1)
	response.Ctx = ctx
	response.Request = request
	response.fixCharset()

	c.handleOnResponse(response)

	c.handleOnHTML(response)

	return nil
}

func (c *Collector) requestCheck(u, method string, depth int) error {
	if u == "" {
		return errors.New("Missing URL")
	}
	if c.MaxDepth > 0 && c.MaxDepth < depth {
		return errors.New("Max depth limit reached")
	}
	if len(c.URLFilters) > 0 {
		matched := false
		for _, r := range c.URLFilters {
			if r.Match([]byte(u)) {
				matched = true
				break
			}
		}
		if !matched {
			return errors.New("No URLFilters match")
		}
	}
	if !c.AllowURLRevisit && method == "GET" {
		h := fnv.New64a()
		h.Write([]byte(u))
		uHash := h.Sum64()
		c.lock.RLock()
		visited := c.visitedURLs[uHash]
		c.lock.RUnlock()
		if visited {
			return errors.New("URL already visited")
		}
		c.lock.Lock()
		c.visitedURLs[uHash] = true
		c.lock.Unlock()
	}
	return nil
}

func (c *Collector) isDomainAllowed(domain string) bool {
	for _, d2 := range c.DisallowedDomains {
		if d2 == domain {
			return false
		}
	}
	if c.AllowedDomains == nil || len(c.AllowedDomains) == 0 {
		return true
	}
	for _, d2 := range c.AllowedDomains {
		if d2 == domain {
			return true
		}
	}
	return false
}

func (c *Collector) checkRobots(u *url.URL) error {
	// var robot *robotstxt.RobotsData
	// var ok bool
	var err error

	c.lock.RLock()
	robot, ok := c.robotsMap[u.Host]
	c.lock.RUnlock()

	if !ok {
		// no robots file cached
		resp, _ := c.backend.Client.Get(u.Scheme + "://" + u.Host + "/robots.txt")
		robot, err = robotstxt.FromResponse(resp)
		if err != nil {
			return err
		}
		c.lock.Lock()
		c.robotsMap[u.Host] = robot
		c.lock.Unlock()
	}

	uaGroup := robot.FindGroup(c.UserAgent)
	if uaGroup == nil {
		return nil
	}

	if !uaGroup.Test(u.EscapedPath()) {
		return errors.New("URL blocked by robots.txt")
	}
	return nil
}

// String is the text representation of the collector.
// It contains useful debug information about the collector's internals
func (c *Collector) String() string {
	return fmt.Sprintf(
		"Requests made: %d (%d responses) | Callbacks: OnRequest: %d, OnHTML: %d, OnResponse: %d, OnError: %d",
		c.requestCount,
		c.responseCount,
		len(c.requestCallbacks),
		len(c.htmlCallbacks),
		len(c.responseCallbacks),
		len(c.errorCallbacks),
	)
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
// element matched by the GoQuery Selector parameter.
// GoQuery Selector is a selector used by https://github.com/PuerkitoBio/goquery
func (c *Collector) OnHTML(goquerySelector string, f HTMLCallback) {
	c.lock.Lock()
	c.htmlCallbacks = append(c.htmlCallbacks, &htmlCallbackContainer{
		Selector: goquerySelector,
		Function: f,
	})
	c.lock.Unlock()
}

// OnHTMLDetach deregister a function. Function will not be execute after detached
func (c *Collector) OnHTMLDetach(goquerySelector string) {
	c.lock.Lock()
	deleteIdx := -1
	for i, cc := range c.htmlCallbacks {
		if cc.Selector == goquerySelector {
			deleteIdx = i
			break
		}
	}
	if deleteIdx != -1 {
		c.htmlCallbacks = append(c.htmlCallbacks[:deleteIdx], c.htmlCallbacks[deleteIdx+1:]...)
	}
	c.lock.Unlock()
}

// OnError registers a function. Function will be executed if an error
// occurs during the HTTP request.
func (c *Collector) OnError(f ErrorCallback) {
	c.lock.Lock()
	c.errorCallbacks = append(c.errorCallbacks, f)
	c.lock.Unlock()
}

// WithTransport allows you to set a custom http.RoundTripper (transport)
func (c *Collector) WithTransport(transport http.RoundTripper) {
	c.backend.Client.Transport = transport
}

// DisableCookies turns off cookie handling
func (c *Collector) DisableCookies() {
	c.backend.Client.Jar = nil
}

// SetCookieJar overrides the previously set cookie jar
func (c *Collector) SetCookieJar(j *cookiejar.Jar) {
	c.backend.Client.Jar = j
}

// SetRequestTimeout overrides the default timeout (10 seconds) for this collector
func (c *Collector) SetRequestTimeout(timeout time.Duration) {
	c.backend.Client.Timeout = timeout
}

// SetProxy sets a proxy for the collector. This overrides the previously
// used http.Transport if the type of the transport is not http.RoundTripper
func (c *Collector) SetProxy(proxyURL string) error {
	proxyParsed, err := url.Parse(proxyURL)
	if err != nil {
		return err
	}

	t, ok := c.backend.Client.Transport.(*http.Transport)
	if c.backend.Client.Transport != nil && ok {
		t.Proxy = http.ProxyURL(proxyParsed)
	} else {
		c.backend.Client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyParsed),
		}
	}

	return nil
}

func createEvent(eventType string, requestId, collectorId int32, kvargs map[string]string) *debug.Event {
	return &debug.Event{
		CollectorId: collectorId,
		RequestId:   requestId,
		Type:        eventType,
		Values:      kvargs,
	}
}

func (c *Collector) handleOnRequest(r *Request) {
	if c.debugger != nil {
		c.debugger.Event(createEvent("request", r.Id, c.Id, map[string]string{
			"url": r.URL.String(),
		}))
	}
	for _, f := range c.requestCallbacks {
		f(r)
	}
}

func (c *Collector) handleOnResponse(r *Response) {
	if c.debugger != nil {
		c.debugger.Event(createEvent("response", r.Request.Id, c.Id, map[string]string{
			"url":    r.Request.URL.String(),
			"status": http.StatusText(r.StatusCode),
		}))
	}
	for _, f := range c.responseCallbacks {
		f(r)
	}
}

func (c *Collector) handleOnHTML(resp *Response) {
	if strings.Index(strings.ToLower(resp.Headers.Get("Content-Type")), "html") == -1 {
		return
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(resp.Body))
	if err != nil {
		return
	}
	for _, cc := range c.htmlCallbacks {
		doc.Find(cc.Selector).Each(func(i int, s *goquery.Selection) {
			for _, n := range s.Nodes {
				e := &HTMLElement{
					Name:       n.Data,
					Request:    resp.Request,
					Response:   resp,
					Text:       goquery.NewDocumentFromNode(n).Text(),
					DOM:        s,
					attributes: n.Attr,
				}
				if c.debugger != nil {
					c.debugger.Event(createEvent("html", resp.Request.Id, c.Id, map[string]string{
						"selector": cc.Selector,
						"url":      resp.Request.URL.String(),
					}))
				}
				cc.Function(e)
			}
		})
	}
}

func (c *Collector) handleOnError(response *Response, err error, request *Request, ctx *Context) error {
	if err == nil && response.StatusCode < 203 {
		return nil
	}
	if err == nil {
		err = errors.New(http.StatusText(response.StatusCode))
	}
	if response == nil {
		response = &Response{
			Request: request,
			Ctx:     ctx,
		}
	}
	if c.debugger != nil {
		c.debugger.Event(createEvent("error", request.Id, c.Id, map[string]string{
			"url":    request.URL.String(),
			"status": http.StatusText(response.StatusCode),
		}))
	}
	if response.Request == nil {
		response.Request = request
	}
	for _, f := range c.errorCallbacks {
		f(response, err)
	}
	return err
}

// Limit adds a new LimitRule to the collector
func (c *Collector) Limit(rule *LimitRule) error {
	return c.backend.Limit(rule)
}

// Limits adds new LimitRules to the collector
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

// Clone creates an exact copy of a Collector without callbacks.
// HTTP backend, robots.txt cache and cookie jar are shared
// between collectors.
func (c *Collector) Clone() *Collector {
	return &Collector{
		UserAgent:         c.UserAgent,
		MaxDepth:          c.MaxDepth,
		visitedURLs:       make(map[uint64]bool, 0),
		htmlCallbacks:     make([]*htmlCallbackContainer, 0, 8),
		requestCallbacks:  make([]RequestCallback, 0, 8),
		responseCallbacks: make([]ResponseCallback, 0, 8),
		errorCallbacks:    make([]ErrorCallback, 0, 8),
		CacheDir:          c.CacheDir,
		MaxBodySize:       c.MaxBodySize,
		backend:           c.backend,
		wg:                c.wg,
		lock:              c.lock,
		robotsMap:         c.robotsMap,
		IgnoreRobotsTxt:   c.IgnoreRobotsTxt,
		Id:                atomic.AddInt32(&collectorCounter, 1),
		debugger:          c.debugger,
	}
}

func (c *Collector) checkRedirectFunc() func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if !c.isDomainAllowed(req.URL.Host) {
			return fmt.Errorf("Not following redirect to %s because its not in AllowedDomains", req.URL.Host)
		}

		// Honor golangs default of maximum of 10 redirects
		if len(via) >= 10 {
			return http.ErrUseLastResponse
		}

		lastRequest := via[len(via)-1]

		// Copy the headers from last request
		for hName, hValues := range lastRequest.Header {
			for _, hValue := range hValues {
				req.Header.Set(hName, hValue)
			}
		}

		// If domain has changed, remove the Authorization-header if it exists
		if req.URL.Host != lastRequest.URL.Host {
			req.Header.Del("Authorization")
		}

		return nil
	}
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

// ChildText returns the concatenated and stripped text content of the matching
// elements.
func (h *HTMLElement) ChildText(goquerySelector string) string {
	return strings.TrimSpace(h.DOM.Find(goquerySelector).Text())
}

// ChildAttr returns the stripped text content of the first matching
// element's attribute.
func (h *HTMLElement) ChildAttr(goquerySelector, attrName string) string {
	if attr, ok := h.DOM.Find(goquerySelector).Attr(attrName); ok {
		return strings.TrimSpace(attr)
	}
	return ""
}

// ChildAttrs returns the stripped text content of all the matching
// element's attributes.
func (h *HTMLElement) ChildAttrs(goquerySelector, attrName string) []string {
	res := make([]string, 0)
	h.DOM.Find(goquerySelector).Each(func(_ int, s *goquery.Selection) {
		if attr, ok := s.Attr(attrName); ok {
			res = append(res, strings.TrimSpace(attr))
		}
	})
	return res
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
// Visit also calls the previously provided callbacks
func (r *Request) Visit(URL string) error {
	return r.collector.scrape(r.AbsoluteURL(URL), "GET", r.Depth+1, nil, r.Ctx, nil)
}

// Post continues a collector job by creating a POST request and preserves the Context
// of the previous request.
// Post also calls the previously provided callbacks
func (r *Request) Post(URL string, requestData map[string]string) error {
	return r.collector.scrape(r.AbsoluteURL(URL), "POST", r.Depth+1, createFormReader(requestData), r.Ctx, nil)
}

// PostRaw starts a collector job by creating a POST request with raw binary data.
// PostRaw preserves the Context of the previous request
// and calls the previously provided callbacks
func (r *Request) PostRaw(URL string, requestData []byte) error {
	return r.collector.scrape(r.AbsoluteURL(URL), "POST", r.Depth+1, bytes.NewReader(requestData), r.Ctx, nil)
}

// PostMultipart starts a collector job by creating a Multipart POST request
// with raw binary data.  PostMultipart also calls the previously provided.
// callbacks
func (r *Request) PostMultipart(URL string, requestData map[string][]byte) error {
	boundary := randomBoundary()
	hdr := http.Header{}
	hdr.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	hdr.Set("User-Agent", r.collector.UserAgent)
	return r.collector.scrape(r.AbsoluteURL(URL), "POST", r.Depth+1, createMultipartReader(boundary, requestData), r.Ctx, hdr)
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

// Put stores a value of any type in Context
func (c *Context) Put(key string, value interface{}) {
	c.lock.Lock()
	c.contextMap[key] = value
	c.lock.Unlock()
}

// Get retrieves a string value from Context.
// Get returns an empty string if key not found
func (c *Context) Get(key string) string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if v, ok := c.contextMap[key]; ok {
		return v.(string)
	}
	return ""
}

// GetAny retrieves a value from Context.
// GetAny returns nil if key not found
func (c *Context) GetAny(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if v, ok := c.contextMap[key]; ok {
		return v
	}
	return nil
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

// SanitizeFileName replaces dangerous characters in a string
// so the return value can be used as a safe file name.
func SanitizeFileName(fileName string) string {
	ext := filepath.Ext(fileName)
	cleanExt := sanitize.BaseName(ext)
	if cleanExt == "" {
		cleanExt = ".unknown"
	}
	return strings.Replace(fmt.Sprintf(
		"%s.%s",
		sanitize.BaseName(fileName[:len(fileName)-len(ext)]),
		cleanExt[1:],
	), "-", "_", -1)
}

func createFormReader(data map[string]string) io.Reader {
	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}
	return strings.NewReader(form.Encode())
}

func createMultipartReader(boundary string, data map[string][]byte) io.Reader {
	dashBoundary := "--" + boundary

	body := []byte{}
	buffer := bytes.NewBuffer(body)

	buffer.WriteString("Content-type: multipart/form-data; boundary=" + boundary + "\n\n")
	for contentType, content := range data {
		buffer.WriteString(dashBoundary + "\n")
		buffer.WriteString("Content-Disposition: form-data; name=" + contentType + "\n")
		buffer.WriteString(fmt.Sprintf("Content-Length: %d \n\n", len(content)))
		buffer.Write(content)
		buffer.WriteString("\n")
	}
	buffer.WriteString(dashBoundary + "--\n\n")
	return buffer
}

// randomBoundary was borrowed from
// github.com/golang/go/mime/multipart/writer.go#randomBoundary
func randomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
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
