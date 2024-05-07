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

import (
	"context"

	"github.com/henderiw/idxtable/pkg/table"
	"github.com/henderiw/idxtable/pkg/tree/gtree"
	"github.com/henderiw/store"
	"github.com/henderiw/store/memory"
)

// newCacheContext holds the cache instance context
// with a status to indicate if it is initialized or not
// initialized false: means it is NOT initialized,
// initialized true means it is initialized
func newCacheContext[T1 any](i T1) *cacheContext[T1] {
	return &cacheContext[T1]{
		initialized: false,
		instance:    i,
	}
}

type cacheContext[T1 any] struct {
	initialized bool
	instance    T1
}

func (r *cacheContext[T1]) Initialized() {
	r.initialized = true
}

func (r *cacheContext[T1]) IsInitialized() bool {
	return r.initialized
}

type CacheContext struct {
	idxType string
	tree    gtree.GTree
	ranges  store.Storer[table.Table]
}

func NewCacheContext(tree gtree.GTree, idxType string) *CacheContext {
	return &CacheContext{
		idxType: idxType, // provides extra context around the
		tree:    tree,
		ranges:  memory.NewStore[table.Table](),
	}
}

func (r *CacheContext) Size() int {
	var size int
	size += r.tree.Size()
	r.ranges.List(context.Background(), func(ctx context.Context, k store.Key, t table.Table) {
		size += t.Size()
	})
	return size
}

func (r *CacheContext) Type() string {
	return r.idxType
}
