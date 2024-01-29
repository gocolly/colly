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
	"context"
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
	"time"

	"google.golang.org/appengine/urlfetch"
)

type appEngineBackend struct {
	LimitRules []*LimitRule
	Client     *http.Client
	lock       *sync.RWMutex
}

func (h *appEngineBackend) Init(jar http.CookieJar) {
	rand.Seed(time.Now().UnixNano())

	ctx := context.Background()
	h.Client = urlfetch.Client(ctx)
	h.Client.Jar = jar
	h.Client.Timeout = 10 * time.Second

	h.lock = &sync.RWMutex{}
}

func (h *appEngineBackend) GetMatchingRule(domain string) *LimitRule {
	if h.LimitRules == nil {
		return nil
	}
	h.lock.RLock()
	defer h.lock.RUnlock()
	for _, r := range h.LimitRules {
		if r.Match(domain) {
			return r
		}
	}
	return nil
}

func (h *appEngineBackend) Cache(request *http.Request, bodySize int, checkHeadersFunc CheckHeadersFunc, cacheDir string) (*Response, error) {
	if cacheDir == "" || request.Method != "GET" {
		return h.Do(request, bodySize, checkHeadersFunc)
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
	resp, err := h.Do(request, bodySize, checkHeadersFunc)
	if err != nil || resp.StatusCode >= 500 {
		return resp, err
	}
	if _, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, 0750); err != nil {
			return resp, err
		}
	}
	file, err := os.Create(filename + "~")
	if err != nil {
		return resp, err
	}
	if err := gob.NewEncoder(file).Encode(resp); err != nil {
		file.Close()
		return resp, err
	}
	file.Close()
	return resp, os.Rename(filename+"~", filename)
}

func (h *appEngineBackend) Do(request *http.Request, bodySize int, checkHeadersFunc CheckHeadersFunc) (*Response, error) {
	r := h.GetMatchingRule(request.URL.Host)
	if r != nil {
		r.WaitChan <- true
		defer func(r *LimitRule) {
			randomDelay := time.Duration(0)
			if r.RandomDelay != 0 {
				randomDelay = time.Duration(rand.Intn(int(r.RandomDelay)))
			}
			time.Sleep(r.Delay + randomDelay)
			<-r.WaitChan
		}(r)
	}

	res, err := h.Client.Do(request)
	if err != nil {
		return nil, err
	}
	*request = *res.Request

	var bodyReader io.Reader = res.Body
	if bodySize > 0 {
		bodyReader = io.LimitReader(bodyReader, int64(bodySize))
	}
	body, err := ioutil.ReadAll(bodyReader)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	if !checkHeadersFunc(request, res.StatusCode, res.Header) {
		// closing res.Body (see defer above) without reading it aborts
		// the download
		return nil, ErrAbortedAfterHeaders
	}
	return &Response{
		StatusCode: res.StatusCode,
		Body:       body,
		Headers:    &res.Header,
	}, nil
}

func (h *appEngineBackend) Limit(rule *LimitRule) error {
	h.lock.Lock()
	if h.LimitRules == nil {
		h.LimitRules = make([]*LimitRule, 0, 8)
	}
	h.LimitRules = append(h.LimitRules, rule)
	h.lock.Unlock()
	return rule.Init()
}

func (h *appEngineBackend) Limits(rules []*LimitRule) error {
	for _, r := range rules {
		if err := h.Limit(r); err != nil {
			return err
		}
	}
	return nil
}

func (h *appEngineBackend) Jar(j http.CookieJar) {
	h.Client.Jar = j
}

func (h *appEngineBackend) GetJar() http.CookieJar {
	return h.Client.Jar
}

func (h *appEngineBackend) Transport(t http.RoundTripper) {
	h.Client.Transport = t
}

func (h *appEngineBackend) Timeout(t time.Duration) {
	h.Client.Timeout = t
}

func (h *appEngineBackend) GetTimeout() time.Duration {
	return h.Client.Timeout
}

func (h *appEngineBackend) Proxy(pf ProxyFunc) {
	t, ok := h.Client.Transport.(*http.Transport)
	if h.Client.Transport != nil && ok {
		t.Proxy = pf
		t.DisableKeepAlives = true
	} else {
		h.Client.Transport = &http.Transport{
			Proxy:             pf,
			DisableKeepAlives: true,
		}
	}
}

func (h *appEngineBackend) SetCookies(url *url.URL, cookies []*http.Cookie) error {
	if h.Client.Jar == nil {
		return ErrNoCookieJar
	}
	h.Client.Jar.SetCookies(url, cookies)
	return nil
}

func (h *appEngineBackend) Cookies(url *url.URL) []*http.Cookie {
	if h.Client.Jar == nil {
		return nil
	}

	return h.Client.Jar.Cookies(url)
}

func (h *appEngineBackend) CheckRedirect(f func(req *http.Request, via []*http.Request) error) {
	h.Client.CheckRedirect = f
}

func (h *appEngineBackend) SetClient(client *http.Client) {
	h.Client = client
}
