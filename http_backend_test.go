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
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

func TestHTTPBackendDoAllowsUnexpectedEOFWithUnknownLength(t *testing.T) {
	body := []byte("decoded response body")
	var compressed bytes.Buffer

	gzipWriter := gzip.NewWriter(&compressed)
	if _, err := gzipWriter.Write(body); err != nil {
		t.Fatalf("Failed to write gzip body: %v", err)
	}
	if err := gzipWriter.Close(); err != nil {
		t.Fatalf("Failed to close gzip body: %v", err)
	}

	truncated := compressed.Bytes()[:compressed.Len()-8]
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		_, _ = w.Write(truncated)
	}))
	defer ts.Close()

	backend := &httpBackend{}
	jar, _ := cookiejar.New(nil)
	backend.Init(jar)

	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := backend.Do(
		req,
		0,
		func(req *http.Request) bool { return true },
		func(req *http.Request, statusCode int, header http.Header) bool { return true },
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !bytes.Equal(resp.Body, body) {
		t.Fatalf("Invalid response body: %q (expected %q)", resp.Body, body)
	}
}
