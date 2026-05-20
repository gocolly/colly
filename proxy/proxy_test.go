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
	"fmt"
	"net/http"
	"net/http/httptest"
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
