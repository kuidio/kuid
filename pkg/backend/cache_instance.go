/*
Copyright 2024 Nokia.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package backend

import "sync"

// newCacheContext holds the cache instance context
// with a status to indicate if it is initialized or not
// initialized false: means it is NOT initialized,
// initialized true means it is initialized
func newCacheInstance[T1 any](i T1) *cacheInstance[T1] {
	return &cacheInstance[T1]{
		initialized: false,
		instance:    i,
	}
}

type cacheInstance[T1 any] struct {
	m           sync.RWMutex
	initialized bool
	instance    T1
}

func (r *cacheInstance[T1]) SetInitialized() {
	r.m.Lock()
	defer r.m.Unlock()
	r.initialized = true
}

func (r *cacheInstance[T1]) IsInitialized() bool {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.initialized
}
