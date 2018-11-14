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
	"sync"
	"time"
)

type ctxKey uint8

const (
	dataCtxKey ctxKey = iota + 1
	nolimitCtxKey
	timingsCtxKey
)

var nolimitCtx = true

func WithNolimitRequestContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, nolimitCtxKey, &nolimitCtx)
}
func ContextNolimitRequest(ctx context.Context) (ok bool) {
	_, ok = ctx.Value(nolimitCtxKey).(*bool)
	return
}

func WithDataContext(ctx context.Context) context.Context {
	dataCtx := &Context{
		contextMap: make(map[string]interface{}),
		lock:       &sync.RWMutex{},
	}
	return context.WithValue(ctx, dataCtxKey, dataCtx)
}

func ContextDataContext(ctx context.Context) *Context {
	u, _ := ctx.Value(dataCtxKey).(*Context)
	return u
}

// Context provides a tiny layer for passing data between callbacks
type Context struct {
	contextMap map[string]interface{}
	lock       *sync.RWMutex
}

// NewContext initializes a new Context instance
func NewContext() *Context {
	return &Context{
		contextMap: make(map[string]interface{}),
		lock:       &sync.RWMutex{},
	}
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

// ForEach iterate context
func (c *Context) ForEach(fn func(k string, v interface{}) interface{}) []interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	ret := make([]interface{}, 0, len(c.contextMap))
	for k, v := range c.contextMap {
		ret = append(ret, fn(k, v))
	}

	return ret
}

func WithTimingsContext(ctx context.Context) context.Context {
	t := &Timings{}
	return context.WithValue(ctx, timingsCtxKey, t)
}

func ContextTimings(ctx context.Context) *Timings {
	return ctx.Value(timingsCtxKey).(*Timings)
}

type Timings struct {
	RequestStart, DownloadStart, DownloadEnd, ProcessStart, ProcessEnd, CharsetFixStart, CharsetFixEnd time.Time
}
