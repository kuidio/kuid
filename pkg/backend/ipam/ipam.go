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
	"fmt"
	"reflect"

	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"

	"github.com/kuidio/kuid/pkg/backend"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func New(c client.Client) backend.Backend[*CacheContext] {
	cache := backend.NewCache[*CacheContext]()

	store := backend.NewNopStore()
	if c != nil {
		store = NewStore(c, cache)
	}
	return &be{
		cache: cache,
		store: store,
	}
}

type be struct {
	cache backend.Cache[*CacheContext]
	store backend.Store
}

// CreateIndex creates a backend index
func (r *be) CreateIndex(ctx context.Context, obj runtime.Object) error {
	cr, ok := obj.(*ipambev1alpha1.IPIndex)
	if !ok {
		return fmt.Errorf("cannot create index expecting %s, got %s", ipambev1alpha1.IPIndexKind, reflect.TypeOf(obj).Name())
	}
	ctx = initIndexContext(ctx, "create", cr)
	log := log.FromContext(ctx)
	log.Info("start")
	key := cr.GetKey()
	//log := log.FromContext(ctx).With("key", key)

	log.Info("start", "isInitialized", r.cache.IsInitialized(ctx, key))
	// if the Cache is not initialized -> restore the cache
	// this happens upon initialization or backend restart
	r.cache.Create(ctx, key, NewCacheContext())
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
func (r *be) DeleteIndex(ctx context.Context, obj runtime.Object) error {
	cr, ok := obj.(*ipambev1alpha1.IPIndex)
	if !ok {
		return fmt.Errorf("cannot delete index expecting %s, got %s", ipambev1alpha1.IPIndexKind, reflect.TypeOf(obj).Name())
	}
	ctx = initIndexContext(ctx, "delete", cr)
	log := log.FromContext(ctx)
	log.Debug("start")
	key := cr.GetKey()

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
func (r *be) Claim(ctx context.Context, obj runtime.Object) error {
	claim, ok := obj.(*ipambev1alpha1.IPClaim)
	if !ok {
		return fmt.Errorf("cannot claim ip expecting %s, got %s", ipambev1alpha1.IPIndexKind, reflect.TypeOf(obj).Name())
	}
	ctx = initClaimContext(ctx, "create", claim)
	log := log.FromContext(ctx)
	log.Debug("start")

	cacheCtx, err := r.cache.Get(ctx, claim.GetKey(), false)
	if err != nil {
		return err
	}

	a, err := getApplicator(ctx, claim, cacheCtx)
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
func (r *be) Release(ctx context.Context, obj runtime.Object) error {
	claim, ok := obj.(*ipambev1alpha1.IPClaim)
	if !ok {
		return fmt.Errorf("cannot delete ip cliam expecting %s, got %s", ipambev1alpha1.IPIndexKind, reflect.TypeOf(obj).Name())
	}
	ctx = initClaimContext(ctx, "delete", claim)
	log := log.FromContext(ctx)
	log.Debug("start")

	rib, err := r.cache.Get(ctx, claim.GetKey(), false)
	if err != nil {
		return err
	}

	// ip claim delete and store
	a, err := getApplicator(ctx, claim, rib)
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

func (r *be) GetCache(ctx context.Context, key store.Key) (*CacheContext, error) {
	return r.cache.Get(ctx, key, false)
}

func getApplicator(_ context.Context, claim *ipambev1alpha1.IPClaim, cacheCtx *CacheContext) (Applicator, error) {

	addressing, err := claim.GetAddressing()
	if err != nil {
		return nil, err
	}
	var a Applicator
	switch addressing {
	case ipambev1alpha1.IPClaimAddressing_StaticAddress:
		a = &staticAddressApplicator{name: "staticIPAddress", applicator: applicator{cacheCtx: cacheCtx}}
	case ipambev1alpha1.IPClaimAddressing_StaticPrefix:
		a = &staticPrefixApplicator{name: "staticIPprefix", applicator: applicator{cacheCtx: cacheCtx}}
	case ipambev1alpha1.IPClaimAddressing_StaticRange:
		a = &staticRangeApplicator{name: "staticIPRange", applicator: applicator{cacheCtx: cacheCtx}}
	case ipambev1alpha1.IPClaimAddressing_DynamicAddress:
		a = &dynamicAddressApplicator{name: "dynamicIPRange", applicator: applicator{cacheCtx: cacheCtx}}
	case ipambev1alpha1.IPClaimAddressing_DynamicPrefix:
		a = &dynamicPrefixApplicator{name: "dynamicIPprefix", applicator: applicator{cacheCtx: cacheCtx}}
	default:
		return nil, fmt.Errorf("invalid addressing, got: %s", string(addressing))
	}

	return a, nil
}
