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
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestContextIteration(t *testing.T) {
	ctx := NewContext()
	for i := 0; i < 10; i++ {
		ctx.Put(strconv.Itoa(i), i)
	}
	values := ctx.ForEach(func(k string, v interface{}) interface{} {
		return v.(int)
	})
	if len(values) != 10 {
		t.Fatal("fail to iterate context")
	}
	for _, i := range values {
		v := i.(int)
		if v != ctx.GetAny(strconv.Itoa(v)).(int) {
			t.Fatal("value not equal")
		}
	}
}

func TestContextCloneConcurrentWrite(t *testing.T) {
	ctx := NewContext()
	for i := 0; i < 50; i++ {
		ctx.Put(strconv.Itoa(i), i)
	}
	done := make(chan struct{})
	var wg sync.WaitGroup
	for w := 0; w < 4; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				default:
					ctx.Put("x", 1)
				}
			}
		}()
	}
	for c := 0; c < 8; c++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				default:
					_ = ctx.Clone()
				}
			}
		}()
	}
	time.Sleep(time.Second)
	close(done)
	finished := make(chan struct{})
	go func() { wg.Wait(); close(finished) }()
	select {
	case <-finished:
	case <-time.After(5 * time.Second):
		t.Fatal("Clone deadlocked on recursive read lock")
	}
}

func TestContextClone(t *testing.T) {
	ctxOrg := NewContext()
	for i := 0; i < 10; i++ {
		ctxOrg.Put(strconv.Itoa(i), i)
	}

	ctx := ctxOrg.Clone()
	values := ctx.ForEach(func(k string, v interface{}) interface{} {
		return v.(int)
	})
	if len(values) != 10 {
		t.Fatal("fail to iterate context")
	}
	for _, i := range values {
		v := i.(int)
		if v != ctx.GetAny(strconv.Itoa(v)).(int) {
			t.Fatal("value not equal")
		}
	}
}
