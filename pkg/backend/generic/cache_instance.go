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

package generic

import (
	"github.com/henderiw/idxtable/pkg/table"
	"github.com/henderiw/idxtable/pkg/tree/gtree"
	"github.com/henderiw/store"
	"github.com/henderiw/store/memory"
)

type CacheInstanceContext struct {
	idxType string
	tree    gtree.GTree
	ranges  store.Storer[table.Table]
}

func NewCacheInstanceContext(tree gtree.GTree, idxType string) *CacheInstanceContext {
	return &CacheInstanceContext{
		idxType: idxType, // provides extra context around the
		tree:    tree,
		ranges:  memory.NewStore[table.Table](nil),
	}
}

func (r *CacheInstanceContext) Size() int {
	var size int
	size += r.tree.Size()
	r.ranges.List(func(k store.Key, t table.Table) {
		size += t.Size()
	})
	return size
}

func (r *CacheInstanceContext) Type() string {
	return r.idxType
}
