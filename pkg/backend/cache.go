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
	"fmt"

	"github.com/henderiw/store"
	"github.com/henderiw/store/memory"
)

type Cache[T1 any] interface {
	IsInitialized(ctx context.Context, k store.Key) bool
	SetInitialized(ctx context.Context, k store.Key) error
	Get(ctx context.Context, k store.Key, i bool) (T1, error)
	Create(ctx context.Context, k store.Key, i T1)
	Delete(ctx context.Context, k store.Key)
}

func NewCache[T1 any]() Cache[T1] {
	return &cache[T1]{
		store: memory.NewStore[*cacheContext[T1]](),
	}
}

type cache[T1 any] struct {
	store store.Storer[*cacheContext[T1]]
}

func (r *cache[T1]) Create(ctx context.Context, k store.Key, i T1) {
	cacheCtx, err := r.store.Get(ctx, k)
	if err != nil {
		_ = r.store.Create(ctx, k, newCacheContext(i))
	}
	if cacheCtx == nil {
		_ = r.store.Create(ctx, k, newCacheContext(i))
	}
	//_ = r.store.Delete(ctx, k)
	//_ = r.store.Create(ctx, k, newCacheContext(i))
}

func (r *cache[T1]) Delete(ctx context.Context, k store.Key) {
	_ = r.store.Delete(ctx, k)
}

// Get returns the cache; the initialized flag can be used to return a cahce even if not initialized
func (r *cache[T1]) Get(ctx context.Context, k store.Key, ignoreInitializing bool) (T1, error) {
	cacheCtx, err := r.store.Get(ctx, k)
	if err != nil {
		return *new(T1), fmt.Errorf("index %s not initialized", k.String())
	}

	if !ignoreInitializing && !cacheCtx.IsInitialized() {
		return *new(T1), fmt.Errorf("index %s is initializing", k.String())
	}
	return cacheCtx.instance, nil
}

// IsInitialized returns true if the cache is initialized and false if the cache is
// not initialized
func (r *cache[T1]) IsInitialized(ctx context.Context, k store.Key) bool {
	cacheCtx, err := r.store.Get(ctx, k)
	if err != nil {
		return false
	}
	return cacheCtx.IsInitialized()
}

// SetInitialized sets the status in the cacheContext to initialized
func (r *cache[T1]) SetInitialized(ctx context.Context, k store.Key) error {
	cacheCtx, err := r.store.Get(ctx, k)
	if err != nil {
		return fmt.Errorf("index %s not initialized", k.String())
	}
	cacheCtx.Initialized()
	return nil
}
