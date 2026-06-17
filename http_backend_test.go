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
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"net/http"
	"os"
	"path"
	"testing"
)

// cacheFilename mirrors the cache path logic used by httpBackend.Cache.
// colly normalizes URLs (a bare host gets a trailing slash) before the
// request is issued, so the cache key is sha1 of the normalized URL string.
func cacheFilename(cacheDir, rawURL string) string {
	sum := sha1.Sum([]byte(rawURL))
	hash := hex.EncodeToString(sum[:])
	return path.Join(cacheDir, hash[:2], hash)
}

// writeTruncatedCacheFile encodes a valid Response gob and writes only the
// first half of the bytes, simulating an interrupted/partial cache write.
func writeTruncatedCacheFile(t *testing.T, filename string) {
	t.Helper()

	headers := http.Header{}
	resp := &Response{StatusCode: 200, Body: []byte("cached body"), Headers: &headers}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(resp); err != nil {
		t.Fatalf("encoding response: %v", err)
	}
	data := buf.Bytes()
	if len(data) < 4 {
		t.Fatalf("encoded response too small to truncate: %d bytes", len(data))
	}

	if err := os.MkdirAll(path.Dir(filename), 0750); err != nil {
		t.Fatalf("creating cache dir: %v", err)
	}
	if err := os.WriteFile(filename, data[:len(data)/2], 0600); err != nil {
		t.Fatalf("writing truncated cache file: %v", err)
	}
}

// TestCorruptCacheFileFallsBackToFetch ensures a corrupt/truncated cache file
// does not crash the process with a nil-pointer panic and instead triggers a
// live fetch. See http_backend.go Cache.
func TestCorruptCacheFileFallsBackToFetch(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	cacheDir := t.TempDir()

	// colly appends a trailing slash to a bare host when normalizing the URL,
	// so the cached request URL string is ts.URL + "/".
	filename := cacheFilename(cacheDir, ts.URL+"/")
	writeTruncatedCacheFile(t, filename)

	c := NewCollector(CacheDir(cacheDir))

	var hits int
	var body []byte
	c.OnResponse(func(r *Response) {
		hits++
		body = r.Body
	})

	if err := c.Visit(ts.URL); err != nil {
		t.Fatalf("Visit returned error: %v", err)
	}

	if hits != 1 {
		t.Fatalf("expected 1 live fetch, got %d", hits)
	}
	if !bytes.Equal(body, serverIndexResponse) {
		t.Fatalf("body = %q, want %q", body, serverIndexResponse)
	}

	// The corrupt file should have been replaced with a valid cache entry.
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("expected cache file to be rewritten: %v", err)
	}
	defer file.Close()
	cached := new(Response)
	if err := gob.NewDecoder(file).Decode(cached); err != nil {
		t.Fatalf("rewritten cache file is not a valid gob: %v", err)
	}
	if cached.Headers == nil {
		t.Fatal("rewritten cache file has nil Headers")
	}
}
