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

/*

import (
	"context"
	"fmt"

	"github.com/henderiw/idxtable/pkg/table"
	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Backend interface {
	// CreateIndex creates a backend index
	CreateIndex(ctx context.Context, obj backend.IndexObject) error
	// DeleteIndex deletes a backend index
	DeleteIndex(ctx context.Context, obj backend.IndexObject) error
	// Claim claims an entry in the backend index
	Claim(ctx context.Context, obj backend.ClaimObject) error
	// Release a claim in the backend
	Release(ctx context.Context, obj backend.ClaimObject) error
	// PrintEntries prints the entries of the cache
	PrintEntries(ctx context.Context, k store.Key) error
}

func New(
	c client.Client,
	newIndex func() backend.IndexObject,
	newEntryList func() backend.ObjectList,
	newClaimList func() backend.ObjectList,
	newEntry func(k store.Key, vrange, id string, labels map[string]string) backend.EntryObject,
	indexGVK schema.GroupVersionKind,
	claimGVK schema.GroupVersionKind) Backend {

	cache := NewCache[*CacheContext]()

	store := NewNopStore()
	if c != nil {
		store = NewStore(c, cache, newIndex, newEntryList, newClaimList, newEntry, indexGVK, claimGVK)
	}
	return &be{
		cache: cache,
		store: store,
	}
}

type be struct {
	cache Cache[*CacheContext]
	store Store
}

// CreateIndex creates a backend index
func (r *be) CreateIndex(ctx context.Context, idx backend.IndexObject) error {
	ctx = InitIndexContext(ctx, "create", idx)
	log := log.FromContext(ctx)
	log.Info("start")
	key := idx.GetKey()
	//log := log.FromContext(ctx).With("key", key)

	log.Info("start", "isInitialized", r.cache.IsInitialized(ctx, key))
	// if the Cache is not initialized -> restore the cache
	// this happens upon initialization or backend restart
	cacheCtx := NewCacheContext(idx.GetTree(), idx.GetType())
	r.cache.Create(ctx, key, cacheCtx)
	if r.cache.IsInitialized(ctx, key) {
		log.Info("already initialized")
		return nil
	}
	if err := r.store.Restore(ctx, key); err != nil {
		log.Error("cannot restore index", "error", err.Error())
		return err
	}
	log.Info("finished")
	return r.cache.SetInitialized(ctx, key)
}

// DeleteIndex deletes a backend index
func (r *be) DeleteIndex(ctx context.Context, idx backend.IndexObject) error {
	ctx = InitIndexContext(ctx, "delete", idx)
	log := log.FromContext(ctx)
	log.Debug("start")
	key := idx.GetKey()

	log.Debug("start", "isInitialized", r.cache.IsInitialized(ctx, key))
	// delete the data from the backend
	if err := r.store.Destroy(ctx, key); err != nil {
		log.Error("cannot delete Index", "error", err.Error())
		return err
	}
	r.cache.Delete(ctx, key)

	log.Debug("finished")
	return nil

}

// Claim claims an entry in the backend index
func (r *be) Claim(ctx context.Context, claim backend.ClaimObject) error {
	ctx = InitClaimContext(ctx, "create", claim)
	log := log.FromContext(ctx)
	log.Debug("start")

	cacheCtx, err := r.cache.Get(ctx, claim.GetKey(), false)
	if err != nil {
		return err
	}

	a, err := getApplicator(ctx, cacheCtx, claim)
	if err != nil {
		return err
	}
	if err := a.Validate(ctx, claim); err != nil {
		return err
	}
	if err := a.Apply(ctx, claim); err != nil {
		return err
	}

	// store the resources in the backend
	return r.store.SaveAll(ctx, claim.GetKey())
}

// Release delete a claim in the backend index
func (r *be) Release(ctx context.Context, claim backend.ClaimObject) error {
	ctx = InitClaimContext(ctx, "delete", claim)
	log := log.FromContext(ctx)
	log.Debug("start")

	cacheCtx, err := r.cache.Get(ctx, claim.GetKey(), false)
	if err != nil {
		return err
	}

	// ip claim delete and store
	a, err := getApplicator(ctx, cacheCtx, claim)
	if err != nil {
		// error gets returned when rib is not initialized -> this means we can safely return
		// and pretend nothing is wrong (hence return nil) since the cleanup already happened
		return nil
	}
	if err := a.Delete(ctx, claim); err != nil {
		return err
	}

	return r.store.SaveAll(ctx, claim.GetKey())
}

func getApplicator(_ context.Context, cacheCtx *CacheContext, claim backend.ClaimObject) (Applicator, error) {
	claimType := claim.GetClaimType()
	var a Applicator
	switch claimType {
	case backend.ClaimType_DynamicID:
		a = &dynamicApplicator{name: string(claimType), applicator: applicator{cacheCtx: cacheCtx}}
	case backend.ClaimType_StaticID:
		a = &staticApplicator{name: string(claimType), applicator: applicator{cacheCtx: cacheCtx}}
	case backend.ClaimType_Range:
		a = &rangeApplicator{name: string(claimType), applicator: applicator{cacheCtx: cacheCtx}}
	default:
		return nil, fmt.Errorf("invalid addressing, got: %s", string(claimType))
	}

	return a, nil
}

func (r *be) PrintEntries(ctx context.Context, k store.Key) error {
	cachectx, err := r.cache.Get(ctx, k, false)
	if err != nil {
		return fmt.Errorf("key not found: %s", err.Error())
	}
	fmt.Println("---------")
	for _, entry := range cachectx.tree.GetAll() {
		entry := entry
		fmt.Println("entry", entry.ID().String())
	}
	cachectx.ranges.List(ctx, func(ctx context.Context, k store.Key, t table.Table) {
		fmt.Println("range", k.Name)
		if t != nil {
			for _, entry := range t.GetAll() {
				entry := entry
				fmt.Println("entry", entry.ID().String())
			}
		}

	})
	return nil
}
*/