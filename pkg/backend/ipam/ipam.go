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

	"github.com/henderiw/idxtable/pkg/iptable"
	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	bebackend "github.com/kuidio/kuid/pkg/backend/backend"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func New(c client.Client) bebackend.Backend {
	cache := bebackend.NewCache[*CacheContext]()

	store := bebackend.NewNopStore()
	if c != nil {
		store = NewStore(c, cache)
	}
	return &be{
		cache: cache,
		store: store,
	}
}

type be struct {
	cache bebackend.Cache[*CacheContext]
	store bebackend.Store
}

// CreateIndex creates a backend index
func (r *be) CreateIndex(ctx context.Context, idx backend.IndexObject) error {
	//cr, ok := obj.(*ipamresv1alpha1.NetworkInstance)
	//if !ok {
	//	return fmt.Errorf("cannot create index expecting %s, got %s", ipamresv1alpha1.NetworkInstanceKind, reflect.TypeOf(obj).Name())
	//}
	ctx = bebackend.InitIndexContext(ctx, "create", idx)
	log := log.FromContext(ctx)
	log.Info("start")
	key := idx.GetKey()
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
func (r *be) DeleteIndex(ctx context.Context, idx backend.IndexObject) error {
	//cr, ok := obj.(*ipamresv1alpha1.NetworkInstance)
	//if !ok {
	//	return fmt.Errorf("cannot delete index expecting %s, got %s", ipamresv1alpha1.NetworkInstanceKind, reflect.TypeOf(obj).Name())
	//}
	ctx = bebackend.InitIndexContext(ctx, "delete", idx)
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
func (r *be) Claim(ctx context.Context, obj backend.ClaimObject) error {
	claim, ok := obj.(*ipambev1alpha1.IPClaim)
	if !ok {
		return fmt.Errorf("cannot claim ip expecting %s, got %s", ipambev1alpha1.IPClaimKind, reflect.TypeOf(obj).Name())
	}
	ctx = initClaimContext(ctx, "create", claim)
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
func (r *be) Release(ctx context.Context, obj backend.ClaimObject) error {
	claim, ok := obj.(*ipambev1alpha1.IPClaim)
	if !ok {
		return fmt.Errorf("cannot delete claimm expecting %s, got %s", ipambev1alpha1.IPClaimKind, reflect.TypeOf(obj).Name())
	}
	ctx = initClaimContext(ctx, "delete", claim)
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

func (r *be) GetCache(ctx context.Context, key store.Key) (*CacheContext, error) {
	return r.cache.Get(ctx, key, false)
}

func getApplicator(_ context.Context, cacheCtx *CacheContext, claim *ipambev1alpha1.IPClaim) (Applicator, error) {
	ipClaimType, err := claim.GetIPClaimType()
	if err != nil {
		return nil, err
	}
	var a Applicator
	switch ipClaimType {
	case ipambev1alpha1.IPClaimType_StaticAddress:
		a = &staticAddressApplicator{name: string(ipambev1alpha1.IPClaimType_StaticAddress), applicator: applicator{cacheCtx: cacheCtx}}
	case ipambev1alpha1.IPClaimType_StaticPrefix:
		a = &staticPrefixApplicator{name: string(ipambev1alpha1.IPClaimType_StaticPrefix), applicator: applicator{cacheCtx: cacheCtx}}
	case ipambev1alpha1.IPClaimType_StaticRange:
		a = &staticRangeApplicator{name: string(ipambev1alpha1.IPClaimType_StaticRange), applicator: applicator{cacheCtx: cacheCtx}}
	case ipambev1alpha1.IPClaimType_DynamicAddress:
		a = &dynamicAddressApplicator{name: string(ipambev1alpha1.IPClaimType_DynamicAddress), applicator: applicator{cacheCtx: cacheCtx}}
	case ipambev1alpha1.IPClaimType_DynamicPrefix:
		a = &dynamicPrefixApplicator{name: string(ipambev1alpha1.IPClaimType_DynamicPrefix), applicator: applicator{cacheCtx: cacheCtx}}
	default:
		return nil, fmt.Errorf("invalid addressing, got: %s", string(ipClaimType))
	}

	return a, nil
}

func (r *be) PrintEntries(ctx context.Context, k store.Key) error {
	cachectx, err := r.cache.Get(ctx, k, false)
	if err != nil {
		return fmt.Errorf("key not found: %s", err.Error())
	}
	fmt.Println("---------")
	for _, entry := range cachectx.rib.GetTable() {
		entry := entry
		fmt.Println("entry", entry.String())
	}
	cachectx.ranges.List(ctx, func(ctx context.Context, k store.Key, t iptable.IPTable) {
		fmt.Println("range", k.Name)
		if t != nil {
			for _, entry := range t.GetAll() {
				entry := entry
				fmt.Println("entry", entry.String())
			}
		}

	})
	return nil
}
