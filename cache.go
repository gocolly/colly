// Copyright 2023 Adam Tauber, Andrzej Lichnerowicz
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package colly

import (
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"net/http"
	"os"
	"path"
)

// Cache is an interface which handles caching Collector's responses
// The default Cache of the Collector is the NullCache.
// FileSystemCache keeps compatibility with non-pluggable caching in legacy
// Collector. For this reason, one can set cache folder via `CACHE_DIR`
// environment variable, or by passing CacheDir to NewCollector.
// Collector's caching backend can be changed too by calling new method
// Collector.SetCache
type Cache interface {
	// Init initializes the caching backend
	Init() error
	// Get retrieves a previously cached response for the given request
	Get(request *http.Request) (*Response, error)
	// Put stores a response given request in cache
	Put(request *http.Request, response *Response) error
	// Close finalizes caching backend
	Close() error
}

const (
	// DefaultCacheFolderPermissions is set to rwx(user), rx(group), nothing for others
	DefaultCacheFolderPermissions = 0750
)

var (
	ErrCacheFolderNotConfigured = errors.New("Cache's base folder cannot be empty")
	ErrCacheNotConfigured       = errors.New("Caching backend is not configured")
	ErrRequestNoCache           = errors.New("Request cannot be cached")
	ErrCachedNotFound           = errors.New("Cached response not found")
)

type NullCache struct {
}

func (c *NullCache) Init() error {
	return nil
}

// Get always retrieves an error to force re-download
func (c *NullCache) Get(request *http.Request) (*Response, error) {
	return nil, ErrCachedNotFound
}

func (c *NullCache) Put(request *http.Request, response *Response) error {
	return nil
}

func (c *NullCache) Close() error {
	return nil
}

// FileSystemCache is the default cache backend of colly.
// FileSystemCache keeps responses persisted on the disk.
type FileSystemCache struct {
	BaseDir string
}

// Init ensures that specified base folder exists
func (c *FileSystemCache) Init() error {
	if c.BaseDir == "" {
		return ErrCacheFolderNotConfigured
	}

	return os.MkdirAll(c.BaseDir, DefaultCacheFolderPermissions)
}

func (c *FileSystemCache) getFilenameFromRequest(request *http.Request) (string, string) {
	sum := sha1.Sum([]byte(request.URL.String()))
	hash := hex.EncodeToString(sum[:])
	dir := path.Join(c.BaseDir, hash[:2])
	return dir, path.Join(dir, hash)
}

// Get returns an error for HTTP verbs other than GET and if request headers
// specify `Cache-Control: no-cache`.
func (c *FileSystemCache) Get(request *http.Request) (*Response, error) {
	if request.Method != "GET" || request.Header.Get("Cache-Control") == "no-cache" {
		return nil, ErrRequestNoCache
	}

	_, filename := c.getFilenameFromRequest(request)

	if file, err := os.Open(filename); err == nil {
		resp := new(Response)
		err = gob.NewDecoder(file).Decode(resp)
		file.Close()
		return resp, err
	} else {
		return nil, err
	}
}

// Put persists response on disk. For compatibility with legacy non-pluggable version,
// it keeps only one level of folder hierarchy.
func (c *FileSystemCache) Put(request *http.Request, response *Response) error {
	dir, filename := c.getFilenameFromRequest(request)

	if _, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, DefaultCacheFolderPermissions); err != nil {
			return err
		}
	}
	file, err := os.Create(filename + "~")
	if err != nil {
		return err
	}
	if err := gob.NewEncoder(file).Encode(response); err != nil {
		file.Close()
		return err
	}
	file.Close()
	return os.Rename(filename+"~", filename)
}

func (c *FileSystemCache) Close() error {
	return nil
}
