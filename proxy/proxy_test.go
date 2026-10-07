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

package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gocolly/colly/v2"
)

// TestRoundRobinProxySwitcher_PropagatesProxyURL is the minimal smoke test:
// after a Visit through the switcher, the response must carry a non-empty
// ProxyURL on both Request and Response.
func TestRoundRobinProxySwitcher_PropagatesProxyURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	}))
	defer ts.Close()

	rp, err := RoundRobinProxySwitcher(ts.URL)
	if err != nil {
		t.Fatalf("RoundRobinProxySwitcher: %v", err)
	}

	c := colly.NewCollector(colly.IgnoreRobotsTxt())
	c.SetProxyFunc(rp)

	var called bool
	c.OnResponse(func(r *colly.Response) {
		called = true
		if r.Request.ProxyURL == "" {
			t.Errorf("Request.ProxyURL is empty — ProxyURLKey not propagated")
		}
		if r.ProxyURL == "" {
			t.Errorf("Response.ProxyURL is empty")
		}
	})

	if err := c.Visit("http://example.com/"); err != nil {
		t.Fatalf("Visit: %v", err)
	}
	if !called {
		t.Fatal("OnResponse never fired")
	}
}

// TestRoundRobinProxySwitcher_ProxyURLOnError ensures the chosen proxy URL
// is still recorded when the request fails before any response headers
// arrive (e.g. dial refused) — so OnError can report which proxy was tried.
func TestRoundRobinProxySwitcher_ProxyURLOnError(t *testing.T) {
	ln := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	dead := ln.URL
	ln.Close() // guarantees dial refused on dead

	rp, err := RoundRobinProxySwitcher(dead)
	if err != nil {
		t.Fatalf("RoundRobinProxySwitcher: %v", err)
	}
	c := colly.NewCollector(colly.IgnoreRobotsTxt())
	c.SetProxyFunc(rp)

	var called bool
	c.OnError(func(r *colly.Response, _ error) {
		called = true
		if r.Request.ProxyURL != dead {
			t.Errorf("Request.ProxyURL = %q, want %q", r.Request.ProxyURL, dead)
		}
		if r.ProxyURL != dead {
			t.Errorf("Response.ProxyURL = %q, want %q", r.ProxyURL, dead)
		}
	})

	if err := c.Visit("http://example.com/"); err == nil {
		t.Fatal("expected Visit to fail")
	}
	if !called {
		t.Fatal("OnError never fired")
	}
}

// TestSetProxyFunc_LegacyContextStringPropagates documents the interaction
// between a custom ProxyFunc that follows the legacy "WithContext+string"
// pattern and the current SetProxyFunc wrapper.
//
// The user's *pr = *pr.WithContext(...) mutation only affects the fork that
// net/http.send() created (Client.Timeout triggers forkReq), so the string
// the user writes into ProxyURLKey is discarded along with the fork. What
// actually surfaces on Request.ProxyURL is the *url.URL the ProxyFunc
// returns, written by the wrapper through the *string holder colly placed
// in the (shared) context. To make this concrete the test has the user
// write a marker string that intentionally differs from the returned URL,
// then asserts the URL — not the marker — is what propagates.
func TestSetProxyFunc_LegacyContextStringPropagates(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	}))
	defer ts.Close()

	proxyURL, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatalf("url.Parse: %v", err)
	}
	const userMarker = "user-wrote-this-but-it-should-be-ignored"

	c := colly.NewCollector(colly.IgnoreRobotsTxt())
	c.SetProxyFunc(func(pr *http.Request) (*url.URL, error) {
		ctx := context.WithValue(pr.Context(), colly.ProxyURLKey, userMarker)
		*pr = *pr.WithContext(ctx)
		return proxyURL, nil
	})

	var called bool
	c.OnResponse(func(r *colly.Response) {
		called = true
		if r.Request.ProxyURL != proxyURL.String() {
			t.Errorf("Request.ProxyURL = %q, want %q (from returned *url.URL)", r.Request.ProxyURL, proxyURL.String())
		}
		if r.ProxyURL != proxyURL.String() {
			t.Errorf("Response.ProxyURL = %q, want %q", r.ProxyURL, proxyURL.String())
		}
		if r.Request.ProxyURL == userMarker {
			t.Errorf("Request.ProxyURL leaked the user marker %q — the WithContext+string write must be isolated by forkReq", userMarker)
		}
	})

	if err := c.Visit("http://example.com/"); err != nil {
		t.Fatalf("Visit: %v", err)
	}
	if !called {
		t.Fatal("OnResponse never fired")
	}
}

// TestSetProxyFunc_LegacyContextStringOnError is the error-path counterpart:
// the same legacy WithContext+string ProxyFunc, but the proxy is a dead port
// so the request fails before any response headers. The returned *url.URL
// (not the user's discarded ctx string) must still be reflected on
// Request.ProxyURL / Response.ProxyURL — proving the *string holder write
// from SetProxyFunc's wrapper survives both forkReq and the error path.
func TestSetProxyFunc_LegacyContextStringOnError(t *testing.T) {
	ln := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	dead := ln.URL
	ln.Close() // guarantees dial refused on dead

	proxyURL, err := url.Parse(dead)
	if err != nil {
		t.Fatalf("url.Parse: %v", err)
	}
	const userMarker = "user-wrote-this-but-it-should-be-ignored"

	c := colly.NewCollector(colly.IgnoreRobotsTxt())
	c.SetProxyFunc(func(pr *http.Request) (*url.URL, error) {
		ctx := context.WithValue(pr.Context(), colly.ProxyURLKey, userMarker)
		*pr = *pr.WithContext(ctx)
		return proxyURL, nil
	})

	var called bool
	c.OnError(func(r *colly.Response, _ error) {
		called = true
		if r.Request.ProxyURL != dead {
			t.Errorf("Request.ProxyURL = %q, want %q (from returned *url.URL)", r.Request.ProxyURL, dead)
		}
		if r.ProxyURL != dead {
			t.Errorf("Response.ProxyURL = %q, want %q", r.ProxyURL, dead)
		}
		if r.Request.ProxyURL == userMarker {
			t.Errorf("Request.ProxyURL leaked the user marker %q", userMarker)
		}
	})

	if err := c.Visit("http://example.com/"); err == nil {
		t.Fatal("expected Visit to fail")
	}
	if !called {
		t.Fatal("OnError never fired")
	}
}
