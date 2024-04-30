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
	"errors"
	"fmt"

	"github.com/henderiw/idxtable/pkg/table32"
	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	vxlanbev1alpha1 "github.com/kuidio/kuid/apis/backend/vxlan/v1alpha1"
	"github.com/kuidio/kuid/pkg/backend"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewStore(c client.Client, cache backend.Cache[*CacheContext]) backend.Store {
	return &bestore{
		client: c,
		cache:  cache,
	}
}

type bestore struct {
	client client.Client
	cache  backend.Cache[*CacheContext]
}

func (r *bestore) Restore(ctx context.Context, k store.Key) error {
	log := log.FromContext(ctx).With("key", k.String())

	cacheCtx, err := r.cache.Get(ctx, k, true)
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

	// retrieve the index to be able to restore the min/max values
	index := &vxlanbev1alpha1.VXLANIndex{}
	if err := r.client.Get(ctx, k.NamespacedName, index); err != nil {
		return err
	}

	if err := r.restoreMinMaxRanges(ctx, cacheCtx, curEntries, index); err != nil {
		return err
	}

	if err := r.restoreClaims(ctx, cacheCtx, curEntries, vxlanbev1alpha1.VXLANClaimType_Range, claimmap); err != nil {
		return err
	}
	if err := r.restoreClaims(ctx, cacheCtx, curEntries, vxlanbev1alpha1.VXLANClaimType_StaticID, claimmap); err != nil {
		return err
	}
	if err := r.restoreClaims(ctx, cacheCtx, curEntries, vxlanbev1alpha1.VXLANClaimType_DynamicID, claimmap); err != nil {
		return err
	}

	log.Info("restore entries left", "items", len(curEntries))

	return nil

}

// only used in configmap
func (r *bestore) SaveAll(ctx context.Context, k store.Key) error {
	log := log.FromContext(ctx)
	log.Info("SaveAll", "key", k.String())

	newEntries, err := r.getEntriesFromCache(ctx, k)
	if err != nil {
		return err
	}
	curEntries, err := r.listEntries(ctx, k)
	if err != nil {
		return err
	}
	// debug end
	for _, newEntry := range newEntries {
		newEntry := newEntry
		found := false
		var entry *vxlanbev1alpha1.VXLANEntry
		for idx, curEntry := range curEntries {
			idx := idx
			curEntry := curEntry
			//fmt.Println("saveAll entries", newIPEntry.Name, curIPEntry.Name)
			if curEntry.Namespace == newEntry.Namespace &&
				curEntry.Name == newEntry.Name {
				curEntries = append(curEntries[:idx], curEntries[idx+1:]...)
				found = true
				entry = curEntry
				break
			}
		}
		//fmt.Println("saveAll entries", found, newIPEntry.Name)
		if !found {
			if err := r.client.Create(ctx, newEntry); err != nil {
				log.Error("saveAll create failed", "name", newEntry.Name, "error", err.Error())
				return err
			}
			continue
		}
		entry.Spec = newEntry.Spec
		log.Debug("save all ipEntry update", "ipEntry", entry.Name)
		if err := r.client.Update(ctx, entry); err != nil {
			return err
		}
	}
	for _, curEntry := range curEntries {
		if err := r.client.Delete(ctx, curEntry); err != nil {
			return err
		}
	}
	return nil
}

// Destroy removes the store db
func (r *bestore) Destroy(ctx context.Context, k store.Key) error {
	// no need to delete the ip index as this is what this fn is supposed to do
	return r.deleteEntries(ctx, k)
}

func (r *bestore) getEntriesFromCache(ctx context.Context, k store.Key) ([]*vxlanbev1alpha1.VXLANEntry, error) {
	log := log.FromContext(ctx).With("key", k.String())
	cacheCtx, err := r.cache.Get(ctx, k, false)
	if err != nil {
		log.Error("cannot get index", "error", err.Error())
		return nil, err
	}

	entries := make([]*vxlanbev1alpha1.VXLANEntry, 0, cacheCtx.Size())
	// add the main rib entry
	for _, entry := range cacheCtx.tree.GetAll() {
		//fmt.Println("getEntriesFromCache rib entry", route.Prefix().String())
		entry := entry
		entries = append(entries, vxlanbev1alpha1.GetVXLANEntry(ctx, k, "", entry.ID().String(), entry.Labels()))
	}
	// add all the range entries
	cacheCtx.ranges.List(ctx, func(ctx context.Context, key store.Key, t *table32.Table32) {
		for _, entry := range t.GetAll() {
			//fmt.Println("getEntriesFromCache range", key.Name, route.Prefix().String())
			entry := entry
			entries = append(entries, vxlanbev1alpha1.GetVXLANEntry(ctx, k, key.Name, entry.ID().String(), entry.Labels()))
		}
	})

	return entries, nil
}

func (r *bestore) deleteEntries(ctx context.Context, k store.Key) error {
	log := log.FromContext(ctx).With("key", k.String())

	entries, err := r.listEntries(ctx, k)
	if err != nil {
		log.Error("cannot list entries", "error", err)
		return err
	}

	var errm error
	for _, entry := range entries {
		if err := r.client.Delete(ctx, entry); err != nil {
			log.Error("cannot delete entry", "error", err)
			errm = errors.Join(errm, err)
			continue
		}
	}
	return errm
}

func (r *bestore) listEntries(ctx context.Context, k store.Key) ([]*vxlanbev1alpha1.VXLANEntry, error) {
	opt := []client.ListOption{
		//client.MatchingFields{
		//	"spec.networkInstance": k.Name,
		//},
	}

	entryList := vxlanbev1alpha1.VXLANEntryList{}
	if err := r.client.List(ctx, &entryList, opt...); err != nil {
		return nil, err
	}
	entries := []*vxlanbev1alpha1.VXLANEntry{}
	for _, entry := range entryList.Items {
		entry := entry
		if entry.Spec.Index == k.Name {
			entries = append(entries, &entry)
		}
	}

	return entries, nil
}

func (r *bestore) listClaims(ctx context.Context, k store.Key) (map[string]*vxlanbev1alpha1.VXLANClaim, error) {
	opt := []client.ListOption{
		/*
			client.MatchingFields{
				"spec.networkInstance": k.Name,
			},
		*/
	}

	claims := vxlanbev1alpha1.VXLANClaimList{}
	if err := r.client.List(ctx, &claims, opt...); err != nil {
		return nil, err
	}

	claimmap := map[string]*vxlanbev1alpha1.VXLANClaim{}
	for _, claim := range claims.Items {
		claim := claim
		if claim.Spec.Index == k.Name {
			claimmap[(&claim).GetNamespacedName().String()] = &claim
		}

	}

	return claimmap, nil
}

func (r *bestore) restoreMinMaxRanges(ctx context.Context, cacheCtx *CacheContext, entries []*vxlanbev1alpha1.VXLANEntry, index *vxlanbev1alpha1.VXLANIndex) error {
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		if entry.Spec.Owner.Group == vxlanbev1alpha1.SchemeGroupVersion.Group &&
			entry.Spec.Owner.Version == vxlanbev1alpha1.SchemeGroupVersion.Version &&
			entry.Spec.Owner.Kind == vxlanbev1alpha1.VXLANIndexKind {

			entries = append(entries[:i], entries[i+1:]...)
		}
	}

	if index.Spec.MinID != nil && *index.Spec.MinID != vxlanbev1alpha1.VXLANID_Min {
		claim := index.GetMinClaim()
		if err := r.restoreClaim(ctx, cacheCtx, claim); err != nil {
			return err
		}
	}
	if index.Spec.MaxID != nil && *index.Spec.MaxID != vxlanbev1alpha1.VXLANID_Max {
		claim := index.GetMaxClaim()
		if err := r.restoreClaim(ctx, cacheCtx, claim); err != nil {
			return err
		}
	}
	return nil
}

func (r *bestore) restoreClaims(ctx context.Context, cacheCtx *CacheContext, entries []*vxlanbev1alpha1.VXLANEntry, claimType vxlanbev1alpha1.VXLANClaimType, claimmap map[string]*vxlanbev1alpha1.VXLANClaim) error {

	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		if entry.Spec.Owner.Group == vxlanbev1alpha1.SchemeGroupVersion.Group &&
			entry.Spec.Owner.Version == vxlanbev1alpha1.SchemeGroupVersion.Version &&
			entry.Spec.Owner.Kind == vxlanbev1alpha1.VXLANClaimKind {

			if claimType == entry.Spec.ClaimType {
				nsn := types.NamespacedName{Namespace: entry.Spec.Owner.Namespace, Name: entry.Spec.Owner.Name}

				claim, ok := claimmap[nsn.String()]
				if ok {
					if err := r.restoreClaim(ctx, cacheCtx, claim); err != nil {
						return err
					}
					// remove the entry since it is processed
					entries = append(entries[:i], entries[i+1:]...)
					delete(claimmap, nsn.String()) // delete the entry to optimize
				}
			}
		}
	}
	return nil
}

func (r *bestore) restoreClaim(ctx context.Context, cacheCtx *CacheContext, claim *vxlanbev1alpha1.VXLANClaim) error {
	ctx = initClaimContext(ctx, "restore", claim)
	a, err := getApplicator(ctx, cacheCtx, claim)
	if err != nil {
		return err
	}
	// validate is needed, mainly for addresses since the parent route determines
	// e.g. the fact the address belongs to a range or not
	errList := claim.ValidateSyntax() // needed to expand the createPrefix/prefixLength and owner
	if len(errList) != 0 {
		return fmt.Errorf("invalid syntax %v", errList)
	}
	if err := a.Validate(ctx, claim); err != nil {
		return err
	}
	if err := a.Apply(ctx, claim); err != nil {
		return err
	}
	return nil
}