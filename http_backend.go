package colly

import (
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path"
	"regexp"
	"sync"
	"time"

	"github.com/gobwas/glob"
)

type httpBackend struct {
	LimitRules []*LimitRule
	Client     *http.Client
	lock       *sync.RWMutex
}

// LimitRule provides connection restrictions for domains.
// Both DomainRegexp and DomainGlob can be used to specify
// the included domains patterns, but at least one is required.
// There can be two kind of limitations:
//  - Parallelism: Set limit for the number of concurrent requests to matching domains
//  - Delay: Wait specified amount of time between requests (parallelism is 1 in this case)
type LimitRule struct {
	// DomainRegexp is a regular expression to match against domains
	DomainRegexp string
	// DomainRegexp is a glob pattern to match against domains
	DomainGlob string
	// Delay is the duration to wait before creating a new request to the matching domains
	Delay time.Duration
	// RandomDelay is the extra randomized duration to wait added to Delay before creating a new request
	RandomDelay time.Duration
	// Parallelism is the number of the maximum allowed concurrent requests of the matching domains
	Parallelism    int
	waitChan       chan bool
	compiledRegexp *regexp.Regexp
	compiledGlob   glob.Glob
}

// Init initializes the private members of LimitRule
func (r *LimitRule) Init() error {
	waitChanSize := 1
	if r.Parallelism > 1 {
		waitChanSize = r.Parallelism
	}
	r.waitChan = make(chan bool, waitChanSize)
	hasPattern := false
	if r.DomainRegexp != "" {
		c, err := regexp.Compile(r.DomainRegexp)
		if err != nil {
			return err
		}
		r.compiledRegexp = c
		hasPattern = true
	}
	if r.DomainGlob != "" {
		c, err := glob.Compile(r.DomainGlob)
		if err != nil {
			return err
		}
		r.compiledGlob = c
		hasPattern = true
	}
	if !hasPattern {
		return errors.New("No pattern defined in LimitRule")
	}
	return nil
}

func (h *httpBackend) Init() {
	rand.Seed(time.Now().UnixNano())
	h.LimitRules = make([]*LimitRule, 0, 8)
	jar, _ := cookiejar.New(nil)
	h.Client = &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}
	h.lock = &sync.RWMutex{}
}

// Match checks that the domain parameter triggers the rule
func (r *LimitRule) Match(domain string) bool {
	match := false
	if r.compiledRegexp != nil && r.compiledRegexp.MatchString(domain) {
		match = true
	}
	if r.compiledGlob != nil && r.compiledGlob.Match(domain) {
		match = true
	}
	return match
}

func (h *httpBackend) GetMatchingRule(domain string) *LimitRule {
	h.lock.RLock()
	defer h.lock.RUnlock()
	for _, r := range h.LimitRules {
		if r.Match(domain) {
			return r
		}
	}
	return nil
}

func (h *httpBackend) Cache(request *http.Request, bodySize int, cacheDir string) (*Response, error) {
	if cacheDir == "" || request.Method != "GET" {
		return h.Do(request, bodySize)
	}
	sum := sha1.Sum([]byte(request.URL.String()))
	hash := hex.EncodeToString(sum[:])
	dir := path.Join(cacheDir, hash[:2])
	filename := path.Join(dir, hash)
	if file, err := os.Open(filename); err == nil {
		resp := new(Response)
		err := gob.NewDecoder(file).Decode(resp)
		file.Close()
		if resp.StatusCode < 500 {
			return resp, err
		}
	}
	resp, err := h.Do(request, bodySize)
	if err != nil || resp.StatusCode >= 500 {
		return resp, err
	}
	if _, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, 0750); err != nil {
			return resp, err
		}
	}
	file, err := os.Create(filename + "~")
	defer file.Close()
	if err != nil {
		return resp, err
	}
	if err := gob.NewEncoder(file).Encode(resp); err != nil {
		return resp, err
	}
	return resp, os.Rename(filename+"~", filename)
}

func (h *httpBackend) Do(request *http.Request, bodySize int) (*Response, error) {
	r := h.GetMatchingRule(request.URL.Host)
	if r != nil {
		r.waitChan <- true
		defer func(r *LimitRule) {
			randomDelay := time.Duration(0)
			if r.RandomDelay != 0 {
				randomDelay = time.Duration(rand.Intn(int(r.RandomDelay)))
			}
			time.Sleep(r.Delay + randomDelay)
			<-r.waitChan
		}(r)
	}

	res, err := h.Client.Do(request)
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader = res.Body
	if bodySize > 0 {
		bodyReader = io.LimitReader(bodyReader, int64(bodySize))
	}
	body, err := ioutil.ReadAll(bodyReader)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	return &Response{
		StatusCode: res.StatusCode,
		Body:       body,
		Headers:    &res.Header,
	}, nil
}

func (h *httpBackend) Limit(rule *LimitRule) error {
	h.lock.Lock()
	h.LimitRules = append(h.LimitRules, rule)
	h.lock.Unlock()
	return rule.Init()
}

func (h *httpBackend) Limits(rules []*LimitRule) error {
	for _, r := range rules {
		if err := h.Limit(r); err != nil {
			return err
		}
	}
	return nil
}
