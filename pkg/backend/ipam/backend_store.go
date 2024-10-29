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
	"sort"

	"github.com/henderiw/idxtable/pkg/iptable"
	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend/ipam"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
)

func (r *be) restore(ctx context.Context, index *ipam.IPIndex) error {
	log := log.FromContext(ctx)
	k := index.GetKey()

	cacheInstanceCtx, err := r.cache.Get(ctx, k)
	if err != nil {
		log.Error("cannot get index", "error", err.Error())
		return err
	}

	// Fetch the current entries that were stored
	curEntries, err := r.listEntries(ctx, k)
	if err != nil {
		return err
	}

	claimmap, err := r.listClaims(ctx, k)
	if err != nil {
		return nil
	}

	/*
		prefixes := make(map[string]ipamv1alpha1.Prefix)
		for _, prefix := range index.Spec.Prefixes {
			prefixes[prefix.Prefix] = prefix
		}
	*/

	/*
		if err := r.restoreIndexPrefixes(ctx, cacheInstanceCtx, curEntries, index, prefixes); err != nil {
			return err
		}
	*/
	if err := r.restoreClaims(ctx, cacheInstanceCtx, curEntries, ipam.IPIndexKind, ipam.IPClaimType_StaticPrefix, claimmap); err != nil {
		return err
	}
	if err := r.restoreClaims(ctx, cacheInstanceCtx, curEntries, ipam.IPClaimKind, ipam.IPClaimType_StaticPrefix, claimmap); err != nil {
		return err
	}
	if err := r.restoreClaims(ctx, cacheInstanceCtx, curEntries, ipam.IPClaimKind, ipam.IPClaimType_StaticRange, claimmap); err != nil {
		return err
	}
	if err := r.restoreClaims(ctx, cacheInstanceCtx, curEntries, ipam.IPClaimKind, ipam.IPClaimType_DynamicPrefix, claimmap); err != nil {
		return err
	}
	if err := r.restoreClaims(ctx, cacheInstanceCtx, curEntries, ipam.IPClaimKind, ipam.IPClaimType_StaticAddress, claimmap); err != nil {
		return err
	}
	if err := r.restoreClaims(ctx, cacheInstanceCtx, curEntries, ipam.IPClaimKind, ipam.IPClaimType_DynamicAddress, claimmap); err != nil {
		return err
	}
	log.Debug("restore prefixes entries left", "items", len(curEntries))

	return nil
}

func (r *be) saveAll(ctx context.Context, k store.Key) error {
	log := log.FromContext(ctx)
	log.Debug("SaveAll", "key", k.String())

	// entries from the memory cache
	newEntries, err := r.getEntriesFromCache(ctx, k)
	if err != nil {
		return err
	}
	// entries in the apiserver
	curEntries, err := r.listEntries(ctx, k)
	if err != nil {
		return err
	}

	news := []string{}
	for _, newEntry := range newEntries {
		news = append(news, newEntry.Name)
	}
	curs := []string{}
	for _, curEntry := range curEntries {
		curs = append(curs, curEntry.Name)
	}
	sort.Strings(news)
	sort.Strings(curs)

	for _, newEntry := range newEntries {
		log.Debug("SaveAll", "newIPEntry", newEntry.GetNamespacedName(), "apiVersion", newEntry.APIVersion)
		found := false
		var oldEntry *ipam.IPEntry
		for idx, curEntry := range curEntries {
			log.Debug("SaveAll", "curEntry", curEntry.GetNamespacedName(), "apiVersion", newEntry.APIVersion)
			idx := idx
			curEntry := curEntry
			if curEntry.GetNamespacedName() == newEntry.GetNamespacedName() {
				log.Debug("SaveAll delete entry from curEntries", "curEntry", curEntry.GetNamespacedName())
				curEntries = append(curEntries[:idx], curEntries[idx+1:]...)
				log.Debug("SaveAll delete entry from curEntries", "curEntry", curEntry.GetNamespacedName(), "length", len(curEntries))
				found = true
				oldEntry = curEntry
				break
			}
		}

		ctx = genericapirequest.WithNamespace(ctx, newEntry.GetNamespace())
		if !found {
			if err := r.bestorage.CreateEntry(ctx, newEntry); err != nil {
				log.Error("saveAll create failed", "name", newEntry.GetName(), "error", err.Error())
				return err
			}
			continue
		}
		if err := r.bestorage.UpdateEntry(ctx, newEntry, oldEntry); err != nil {
			log.Error("saveAll update failed", "name", newEntry.GetName(), "error", err.Error())
			return err
		}
	}
	//for _, curEntry := range curEntries {
	//	log.Info("saveAll delete entry", "entry", curEntry.GetNamespacedName())
	//	if err := r.bestorage.DeleteEntry(ctx, curEntry); err != nil {
	//		log.Error("saveAll update failed", "name", curEntry.GetName(), "error", err.Error())
	//		return err
	//	}
	//}
	return nil
}

// Destroy removes the store db
func (r *be) destroy(ctx context.Context, k store.Key) error {
	// TBD: what do we do when deleting the index in async mode
	if err := r.deleteClaims(ctx, k); err != nil {
		return err
	}
	return r.deleteEntries(ctx, k)
}

func (r *be) getEntriesFromCache(ctx context.Context, k store.Key) ([]*ipam.IPEntry, error) {
	//log := log.FromContext(ctx).With("key", k.String())
	cacheInstanceCtx, err := r.cache.Get(ctx, k)
	if err != nil {
		return nil, fmt.Errorf("cache index not initialized")
	}

	entries := make([]*ipam.IPEntry, 0, cacheInstanceCtx.Size())
	// add the main rib entry
	for _, route := range cacheInstanceCtx.rib.GetTable() {
		route := route
		entries = append(entries, ipam.GetIPEntry(ctx, k, "", route.Prefix(), route.Labels()))
	}
	// add all the range entries
	cacheInstanceCtx.ranges.List(func(key store.Key, t iptable.IPTable) {
		for _, route := range t.GetAll() {
			route := route
			entries = append(entries, ipam.GetIPEntry(ctx, k, key.Name, route.Prefix(), route.Labels()))
		}
	})

	return entries, nil
}

func (r *be) deleteEntries(ctx context.Context, k store.Key) error {
	log := log.FromContext(ctx)

	entries, err := r.listEntries(ctx, k)
	if err != nil {
		log.Error("cannot list entries", "error", err)
		return err
	}

	var errm error
	for _, curEntry := range entries {
		if err := r.bestorage.DeleteEntry(ctx, curEntry); err != nil {
			log.Error("delete entry failed", "name", curEntry.GetName(), "error", err.Error())
			return err
		}
	}
	return errm
}

func (r *be) deleteClaims(ctx context.Context, k store.Key) error {
	log := log.FromContext(ctx)

	log.Debug("deleteClaims list")
	claims, err := r.listClaims(ctx, k)
	if err != nil {
		log.Error("cannot list claims", "error", err)
		return err
	}

	var errm error
	for _, claim := range claims {
		log.Debug("deleteClaim from storage", "claim", claim.GetName())

		if err := r.bestorage.DeleteClaim(ctx, claim); err != nil {
			log.Error("cannot delete claim", "error", err)
			errm = errors.Join(errm, err)
			continue
		}
	}
	return errm
}

func (r *be) listEntries(ctx context.Context, k store.Key) ([]*ipam.IPEntry, error) {
	return r.bestorage.ListEntries(ctx, k)
}

func (r *be) listClaims(ctx context.Context, k store.Key) (map[string]*ipam.IPClaim, error) {
	return r.bestorage.ListClaims(ctx, k)
}

func (r *be) restoreClaims(ctx context.Context, cacheInstanceCtx *CacheInstanceContext, entries []*ipam.IPEntry, kind string, claimType ipam.IPClaimType, ipclaimmap map[string]*ipam.IPClaim) error {
	log := log.FromContext(ctx)
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		if (kind == ipam.IPIndexKind && entry.Spec.IndexEntry && claimType == entry.Spec.ClaimType) ||
			(kind != ipam.IPIndexKind && !entry.Spec.IndexEntry && claimType == entry.Spec.ClaimType) {
			claimName := ""
			if len(entry.OwnerReferences) > 0 {
				claimName = entry.OwnerReferences[0].Name
			}

			nsn := types.NamespacedName{Namespace: entry.GetNamespace(), Name: claimName}
			claim, ok := ipclaimmap[nsn.String()]
			if ok {
				log.Debug("restore claim", "kind", kind, "claimType", claimType, "claim", claim)
				if err := r.restoreClaim(ctx, cacheInstanceCtx, claim); err != nil {
					return err
				}
				// remove the entry since it is processed
				entries = append(entries[:i], entries[i+1:]...)
				delete(ipclaimmap, nsn.String()) // delete the entry to optimize
			}
		}

	}
	return nil
}

func (r *be) restoreClaim(ctx context.Context, cacheInstanceCtx *CacheInstanceContext, claim *ipam.IPClaim) error {
	ctx = initClaimContext(ctx, "restore", claim)
	log := log.FromContext(ctx)
	a, err := getApplicator(ctx, cacheInstanceCtx, claim)
	if err != nil {
		return err
	}
	// validate is needed, mainly for addresses since the parent route determines
	// e.g. the fact the address belongs to a range or not
	errList := claim.ValidateSyntax("") // needed to expand the createPrefix/prefixLength and owner
	if len(errList) != 0 {
		return fmt.Errorf("invalid syntax %v", errList)
	}
	if err := a.Validate(ctx, claim); err != nil {
		log.Error("failed to validate claim", "error", err)
		return err
	}
	if err := a.Apply(ctx, claim); err != nil {
		log.Error("failed to apply claim", "error", err)
		return err
	}
	return nil
}

func (r *be) updateIPIndexClaims(ctx context.Context, index *ipam.IPIndex) error {
	log := log.FromContext(ctx)
	log.Debug("updateIPIndexClaims", "key", index.GetKey().String())
	key := index.GetKey()

	newClaims, err := index.GetClaims()
	if err != nil {
		return err
	}

	existingClaims, err := r.listIndexClaims(ctx, key)
	if err != nil {
		return err
	}

	var errm error
	for _, newClaim := range newClaims {
		ctx = genericapirequest.WithNamespace(ctx, newClaim.GetNamespace())
		oldClaim, exists := existingClaims[newClaim.GetNamespacedName().String()]
		if !exists {
			if err := r.bestorage.CreateClaim(ctx, newClaim); err != nil {
				log.Error("updateIPIndexClaims create failed", "name", newClaim.GetName(), "error", err.Error())
				errm = errors.Join(errm, err)
				continue
			}
			continue
		}
		if err := r.bestorage.UpdateClaim(ctx, newClaim, oldClaim); err != nil {
			log.Error("updateIPIndexClaims create failed", "name", newClaim.GetName(), "error", err.Error())
			errm = errors.Join(errm, err)
			continue
		}
		delete(existingClaims, newClaim.GetNamespacedName().String())
	}
	/*
		for _, claim := range existingClaims {
			log.Info("updateIPIndexClaims: delete existing claims", "claim", claim.GetName())
			if err := r.bestorage.DeleteClaim(ctx, claim); err != nil {
				log.Error("updateIPIndexClaims delete failed", "name", claim.GetName(), "error", err.Error())
				errm = errors.Join(errm, err)
				continue
			}
		}
	*/
	if errm != nil {
		return errm
	}
	return r.saveAll(ctx, key)
}

func EntryTransformer(_ context.Context, newObj runtime.Object, oldObj runtime.Object) (runtime.Object, error) {
	// Type assertion to specific object types, assuming we are working with a type that has Spec and Status fields
	new, ok := newObj.(*ipam.IPEntry)
	if !ok {
		return nil, fmt.Errorf("newObj is not of type IPEntry, got: %s", reflect.TypeOf(newObj).Name())
	}
	old, ok := oldObj.(*ipam.IPEntry)
	if !ok {
		return nil, fmt.Errorf("oldObj is not of type IPEntry, got: %s", reflect.TypeOf(oldObj).Name())
	}

	new.SetResourceVersion(old.GetResourceVersion())
	new.SetUID(old.GetUID())

	return new, nil
}

func ClaimTransformer(_ context.Context, newObj runtime.Object, oldObj runtime.Object) (runtime.Object, error) {
	// Type assertion to specific object types, assuming we are working with a type that has Spec and Status fields
	new, ok := newObj.(*ipam.IPClaim)
	if !ok {
		return nil, fmt.Errorf("newObj is not of type IPClaim, got: %s", reflect.TypeOf(newObj).Name())
	}
	old, ok := oldObj.(*ipam.IPClaim)
	if !ok {
		return nil, fmt.Errorf("oldObj is not of type IPClaim, got: %s", reflect.TypeOf(oldObj).Name())
	}

	new.SetResourceVersion(old.GetResourceVersion())
	new.SetUID(old.GetUID())

	return new, nil
}

func (r *be) listIndexClaims(ctx context.Context, k store.Key) (map[string]*ipam.IPClaim, error) {
	return r.bestorage.ListClaims(ctx, k, &ListOptions{
		OwnerKind: ipam.IPIndexKind,
	})
}
