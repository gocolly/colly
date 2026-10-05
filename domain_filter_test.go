// Copyright 2026 Adam Tauber
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

package colly_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gocolly/colly/v2"
)

type domainFilterTransport func(*http.Request) (*http.Response, error)

func (f domainFilterTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestCollectorDomainFilters(t *testing.T) {
	for _, tc := range []struct {
		name         string
		host         string
		allowed      []string
		disallowed   []string
		wantBlocked  bool
		redirectOnly bool
	}{
		{name: "unrestricted", host: "MiXeD.example"},
		{name: "deny_mixed_host", host: "ExAmPlE.COM", disallowed: []string{"example.com"}, wantBlocked: true},
		{name: "deny_uppercase_rule", host: "example.com", disallowed: []string{"EXAMPLE.COM"}, wantBlocked: true},
		{name: "deny_uppercase_host_and_rule", host: "EXAMPLE.COM", disallowed: []string{"EXAMPLE.COM"}, wantBlocked: true},
		{name: "allow_mixed_host", host: "ExAmPlE.COM", allowed: []string{"example.com"}},
		{name: "allow_uppercase_rule", host: "example.com", allowed: []string{"EXAMPLE.COM"}},
		{name: "allow_uppercase_host_and_rule", host: "EXAMPLE.COM", allowed: []string{"EXAMPLE.COM"}},
		{name: "unlisted_domain", host: "OTHER.example", allowed: []string{"example.com"}, wantBlocked: true},
		{name: "deny_precedence", host: "ExAmPlE.COM", allowed: []string{"example.com"}, disallowed: []string{"EXAMPLE.COM"}, wantBlocked: true},
		{name: "trailing_dot_deny_distinct", host: "example.com.", disallowed: []string{"example.com"}},
		{name: "trailing_dot_allow_distinct", host: "example.com.", allowed: []string{"example.com"}, wantBlocked: true},
		{name: "trailing_dot_deny_match", host: "ExAmPlE.COM.", disallowed: []string{"example.com."}, wantBlocked: true},
		{name: "trailing_dot_allow_match", host: "ExAmPlE.COM.", allowed: []string{"example.com."}},
		{name: "deny_with_port", host: "EXAMPLE.COM:8080", disallowed: []string{"example.com"}, wantBlocked: true},
		{name: "allow_with_port", host: "EXAMPLE.COM:8080", allowed: []string{"example.com"}},
		{name: "ipv6_deny", host: "[::1]", disallowed: []string{"::1"}, wantBlocked: true},
		{name: "ipv6_allow", host: "[::1]", allowed: []string{"::1"}},
		{name: "ipv6_zone_deny_match", host: "[fe80::1%25ETH0]", disallowed: []string{"fe80::1%ETH0"}, wantBlocked: true, redirectOnly: true},
		{name: "ipv6_zone_allow_match", host: "[fe80::1%25ETH0]", allowed: []string{"fe80::1%ETH0"}, redirectOnly: true},
		{name: "ipv6_zone_deny_distinct", host: "[fe80::1%25ETH0]", disallowed: []string{"fe80::1%eth0"}, redirectOnly: true},
		{name: "ipv6_zone_allow_distinct", host: "[fe80::1%25ETH0]", allowed: []string{"fe80::1%eth0"}, wantBlocked: true, redirectOnly: true},
	} {
		for _, route := range []string{"direct", "redirect", "redirect_chain"} {
			if tc.redirectOnly && route == "direct" {
				// The initial WHATWG parser rejects IPv6 zone identifiers.
				continue
			}
			t.Run(tc.name+"/"+route, func(t *testing.T) {
				allowed := tc.allowed
				if len(allowed) > 0 {
					allowed = append([]string{"start.example"}, allowed...)
				}
				c := colly.NewCollector(colly.AllowedDomains(allowed...), colly.DisallowedDomains(tc.disallowed...))
				destination := "http://" + tc.host + "/Path?Key=Value"
				targetRequests := 0
				c.WithTransport(domainFilterTransport(func(r *http.Request) (*http.Response, error) {
					response := &http.Response{
						StatusCode: http.StatusOK,
						Header:     make(http.Header),
						Body:       io.NopCloser(strings.NewReader("OK")),
						Request:    r,
					}
					if r.URL.Hostname() == "start.example" {
						location := destination
						if route == "redirect_chain" && r.URL.Path == "/start" {
							location = "/hop"
						}
						response.StatusCode = http.StatusFound
						response.Header.Set("Location", location)
					} else {
						targetRequests++
						if r.URL.Path != "/Path" || r.URL.RawQuery != "Key=Value" {
							t.Errorf("request path or query changed: %s", r.URL)
						}
						if route != "direct" && r.URL.String() != destination {
							t.Errorf("redirect URL changed: got %q, want %q", r.URL, destination)
						}
					}
					return response, nil
				}))
				visit := destination
				if route != "direct" {
					visit = "http://start.example/start"
				}
				err := c.Visit(visit)
				if tc.wantBlocked {
					if !errors.Is(err, colly.ErrForbiddenDomain) {
						t.Errorf("want ErrForbiddenDomain, got %v", err)
					}
					if targetRequests != 0 {
						t.Errorf("blocked destination received %d requests", targetRequests)
					}
				} else {
					if err != nil {
						t.Errorf("allowed destination rejected: %v", err)
					}
					if targetRequests != 1 {
						t.Errorf("want one destination request, got %d", targetRequests)
					}
				}
			})
		}
	}
}
