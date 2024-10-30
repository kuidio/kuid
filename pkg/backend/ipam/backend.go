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
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	"github.com/kform-dev/choreo/apis/condition"
	"github.com/kuidio/kuid/apis/backend/ipam"
	bebackend "github.com/kuidio/kuid/pkg/backend"
	"k8s.io/apimachinery/pkg/runtime"
)

func New() bebackend.Backend {
	cache := bebackend.NewCache[*CacheInstanceContext]()
	return &be{
		cache: cache,
	}
}

type be struct {
	cache bebackend.Cache[*CacheInstanceContext]
	m     sync.RWMutex
	// added later
	//entryStorage *registry.Store
	//claimStorage *registry.Store
	bestorage BackendStorage
}

func (r *be) PrintEntries(ctx context.Context, index string) {
	entries, _ := r.listEntries(ctx, store.ToKey(index))
	for _, entry := range entries {
		uobj, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(entry)
		fmt.Println("entry", uobj)
	}
}

func (r *be) AddStorageInterfaces(bes any) error {
	bestorage, ok := bes.(BackendStorage)
	if !ok {
		return fmt.Errorf("AddStorageInterfaces did not supply a ipam BackendStorage interface, got: %s", reflect.TypeOf(bes).Name())
	}
	r.bestorage = bestorage
	return nil
}

// CreateIndex creates a backend index
func (r *be) CreateIndex(ctx context.Context, obj runtime.Object) error {
	r.m.Lock()
	defer r.m.Unlock()
	index, ok := obj.(*ipam.IPIndex)
	if !ok {
		return errors.New("runtime object is not IPIndex")
	}

	ctx = bebackend.InitIndexContext(ctx, "create", index)
	log := log.FromContext(ctx)
	key := index.GetKey()

	log.Debug("create index", "isInitialized", r.cache.IsInitialized(ctx, key))
	// if the Cache is not initialized -> restore the cache
	// this happens upon initialization or backend restart
	if _, err := r.cache.Get(ctx, key); err != nil {
		// if it does not exist create the cache
		cacheInstanceCtx := NewCacheInstanceContext()
		r.cache.Create(ctx, key, cacheInstanceCtx)
	}
	if !r.cache.IsInitialized(ctx, key) {
		if err := r.restore(ctx, index); err != nil {
			log.Error("cannot restore index", "error", err.Error())
			index.SetConditions(condition.Failed(err.Error()))
			return err
		}
		log.Debug("restored")
		index.SetConditions(condition.Ready())
		obj = index

		if err := r.cache.SetInitialized(ctx, key); err != nil {
			return err
		}
	}
	log.Debug("update IPIndex claims", "object", obj)
	return r.updateIPIndexClaims(ctx, index)
}

// DeleteIndex deletes a backend index
func (r *be) DeleteIndex(ctx context.Context, obj runtime.Object) error {
	r.m.Lock()
	defer r.m.Unlock()
	index, ok := obj.(*ipam.IPIndex)
	if !ok {
		return errors.New("runtime object is not IPIndex")
	}

	ctx = bebackend.InitIndexContext(ctx, "delete", index)
	log := log.FromContext(ctx)
	log.Debug("start")
	key := index.GetKey()

	log.Debug("start", "isInitialized", r.cache.IsInitialized(ctx, key))
	// delete the data from the backend
	if err := r.destroy(ctx, key); err != nil {
		log.Error("cannot delete Index", "error", err.Error())
		return err
	}
	log.Debug("destroyed")
	r.cache.Delete(ctx, key)

	log.Debug("finished")
	return nil
}

func (r *be) Claim(ctx context.Context, obj runtime.Object, recursion bool) error {
	// index delete/create can call the claim create/delete -> this avoid double locking
	if !recursion {
		r.m.Lock()
		defer r.m.Unlock()
	}
	claim, ok := obj.(*ipam.IPClaim)
	if !ok {
		return errors.New("runtime object is not IPClaim")
	}

	ctx = initClaimContext(ctx, "create", claim)
	log := log.FromContext(ctx)
	log.Debug("start")

	cacheCtx, err := r.cache.Get(ctx, claim.GetKey())
	if err != nil {
		return err
	}
	if !r.cache.IsInitialized(ctx, claim.GetKey()) {
		return fmt.Errorf("cache not initialized")
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
	if err := r.saveAll(ctx, claim.GetKey()); err != nil {
		return err
	}
	obj = claim
	return nil

}

func (r *be) Release(ctx context.Context, obj runtime.Object, recursion bool) error {
	// index delete/create can call the claim create/delete -> this avoid double locking
	if !recursion {
		r.m.Lock()
		defer r.m.Unlock()
	}
	claim, ok := obj.(*ipam.IPClaim)
	if !ok {
		return errors.New("runtime object is not IPClaim")
	}

	ctx = initClaimContext(ctx, "delete", claim)
	log := log.FromContext(ctx)
	log.Debug("start")

	cacheCtx, err := r.cache.Get(ctx, claim.GetKey())
	if err != nil {
		return err
	}
	if !r.cache.IsInitialized(ctx, claim.GetKey()) {
		return fmt.Errorf("cache not initialized")
	}

	a, err := getApplicator(ctx, cacheCtx, claim)
	if err != nil {
		return err
	}
	if err := a.Delete(ctx, claim); err != nil {
		return err
	}

	return r.saveAll(ctx, claim.GetKey())
}

func getApplicator(_ context.Context, cacheInstanceCtx *CacheInstanceContext, claim *ipam.IPClaim) (Applicator, error) {
	ipClaimType, err := claim.GetIPClaimType()
	if err != nil {
		return nil, err
	}
	var a Applicator
	switch ipClaimType {
	case ipam.IPClaimType_StaticAddress:
		a = &staticAddressApplicator{name: string(ipam.IPClaimType_StaticAddress), applicator: applicator{cacheInstanceCtx: cacheInstanceCtx}}
	case ipam.IPClaimType_StaticPrefix:
		a = &staticPrefixApplicator{name: string(ipam.IPClaimType_StaticPrefix), applicator: applicator{cacheInstanceCtx: cacheInstanceCtx}}
	case ipam.IPClaimType_StaticRange:
		a = &staticRangeApplicator{name: string(ipam.IPClaimType_StaticRange), applicator: applicator{cacheInstanceCtx: cacheInstanceCtx}}
	case ipam.IPClaimType_DynamicAddress:
		a = &dynamicAddressApplicator{name: string(ipam.IPClaimType_DynamicAddress), applicator: applicator{cacheInstanceCtx: cacheInstanceCtx}}
	case ipam.IPClaimType_DynamicPrefix:
		a = &dynamicPrefixApplicator{name: string(ipam.IPClaimType_DynamicPrefix), applicator: applicator{cacheInstanceCtx: cacheInstanceCtx}}
	default:
		return nil, fmt.Errorf("invalid addressing, got: %s", string(ipClaimType))
	}

	return a, nil
}
