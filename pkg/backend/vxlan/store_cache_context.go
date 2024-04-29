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

package vxlan

import (
	"context"

	"github.com/henderiw/idxtable/pkg/table32"
	"github.com/henderiw/idxtable/pkg/tree32"
	"github.com/henderiw/store"
	"github.com/henderiw/store/memory"
)

type CacheContext struct {
	tree   *tree32.Tree32
	ranges store.Storer[*table32.Table32]
}

func NewCacheContext() *CacheContext {
	return &CacheContext{
		tree:   tree32.New(),
		ranges: memory.NewStore[*table32.Table32](),
	}

}

func (r *CacheContext) Size() int {
	var size int
	size += r.tree.Size()
	r.ranges.List(context.Background(), func(ctx context.Context, k store.Key, t *table32.Table32) {
		size += t.Size()
	})
	return size
}
