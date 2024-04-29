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

package ipam

import (
	"context"

	"github.com/hansthienpondt/nipam/pkg/table"
	"github.com/henderiw/idxtable/pkg/iptable"
	"github.com/henderiw/store"
	"github.com/henderiw/store/memory"
)

type CacheContext struct {
	rib    *table.RIB
	ranges store.Storer[iptable.IPTable]
}

func NewCacheContext() *CacheContext {
	return &CacheContext{
		rib:    table.NewRIB(),
		ranges: memory.NewStore[iptable.IPTable](),
	}

}

func (r *CacheContext) Size() int {
	var size int
	size += r.rib.Size()
	r.ranges.List(context.Background(), func(ctx context.Context, k store.Key, i iptable.IPTable) {
        size +=i.Size()
    })
    return size
}
