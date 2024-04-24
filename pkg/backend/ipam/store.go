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

	"github.com/henderiw/idxtable/pkg/iptable"
	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	ipamresv1alpha1 "github.com/kuidio/kuid/apis/resource/ipam/v1alpha1"
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
	curIPEntries, err := r.listEntries(ctx, k)
	if err != nil {
		return err
	}

	// fetch the NI, IP(s) and IPClaims
	ni, niPrefixes, err := r.getNIPrefixes(ctx, k)
	if err != nil {
		return nil
	}

	ipmap, err := r.listIPs(ctx, k)
	if err != nil {
		return nil
	}
	ipclaimmap, err := r.listClaims(ctx, k)
	if err != nil {
		return nil
	}

	if err := r.restoreNiPrefixes(ctx, cacheCtx, curIPEntries, ni, niPrefixes); err != nil {
		return err
	}
	if err := r.restoreIPs(ctx, cacheCtx, curIPEntries, ipambev1alpha1.IPClaimType_StaticPrefix, ipmap); err != nil {
		return err
	}
	if err := r.restoreIPClaims(ctx, cacheCtx, curIPEntries, ipambev1alpha1.IPClaimType_StaticPrefix, ipclaimmap); err != nil {
		return err
	}
	if err := r.restoreIPClaims(ctx, cacheCtx, curIPEntries, ipambev1alpha1.IPClaimType_StaticRange, ipclaimmap); err != nil {
		return err
	}
	if err := r.restoreIPClaims(ctx, cacheCtx, curIPEntries, ipambev1alpha1.IPClaimType_DynamicPrefix, ipclaimmap); err != nil {
		return err
	}
	if err := r.restoreIPs(ctx, cacheCtx, curIPEntries, ipambev1alpha1.IPClaimType_StaticAddress, ipmap); err != nil {
		return err
	}
	if err := r.restoreIPClaims(ctx, cacheCtx, curIPEntries, ipambev1alpha1.IPClaimType_StaticAddress, ipclaimmap); err != nil {
		return err
	}
	if err := r.restoreIPClaims(ctx, cacheCtx, curIPEntries, ipambev1alpha1.IPClaimType_DynamicAddress, ipclaimmap); err != nil {
		return err
	}

	log.Info("restore prefixes entries left", "items", len(curIPEntries))

	return nil

}

// only used in configmap
func (r *bestore) SaveAll(ctx context.Context, k store.Key) error {
	log := log.FromContext(ctx)
	log.Info("SaveAll", "key", k.String())

	newIPEntries, err := r.getEntriesFromCache(ctx, k)
	if err != nil {
		return err
	}
	curIPEntries, err := r.listEntries(ctx, k)
	if err != nil {
		return err
	}
	// debug end
	for _, newIPEntry := range newIPEntries {
		newIPEntry := newIPEntry
		found := false
		var ipEntry *ipambev1alpha1.IPEntry
		for idx, curIPEntry := range curIPEntries {
			idx := idx
			curIPEntry := curIPEntry
			//fmt.Println("saveAll entries", newIPEntry.Name, curIPEntry.Name)
			if curIPEntry.Namespace == newIPEntry.Namespace &&
				curIPEntry.Name == newIPEntry.Name {
				curIPEntries = append(curIPEntries[:idx], curIPEntries[idx+1:]...)
				found = true
				ipEntry = curIPEntry
				break
			}
		}
		//fmt.Println("saveAll entries", found, newIPEntry.Name)
		if !found {
			if err := r.client.Create(ctx, newIPEntry); err != nil {
				log.Error("saveAll create failed", "name", newIPEntry.Name, "error", err.Error())
				return err
			}
			continue
		}
		ipEntry.Spec = newIPEntry.Spec
		log.Debug("save all ipEntry update", "ipEntry", ipEntry.Name)
		if err := r.client.Update(ctx, ipEntry); err != nil {
			return err
		}
	}
	for _, curIPEntry := range curIPEntries {
		if err := r.client.Delete(ctx, curIPEntry); err != nil {
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

func (r *bestore) getEntriesFromCache(ctx context.Context, k store.Key) ([]*ipambev1alpha1.IPEntry, error) {
	log := log.FromContext(ctx).With("key", k.String())
	cacheCtx, err := r.cache.Get(ctx, k, false)
	if err != nil {
		log.Error("cannot get index", "error", err.Error())
		return nil, err
	}

	ipEntries := make([]*ipambev1alpha1.IPEntry, 0, cacheCtx.Size())
	// add the main rib entry
	for _, route := range cacheCtx.rib.GetTable() {
		//fmt.Println("getEntriesFromCache rib entry", route.Prefix().String())
		route := route
		ipEntries = append(ipEntries, ipambev1alpha1.GetIPEntry(ctx, k, route.Prefix(), route.Labels()))
	}
	// add all the range entries
	cacheCtx.ranges.List(ctx, func(ctx context.Context, key store.Key, i iptable.IPTable) {
		for _, route := range i.GetAll() {
			//fmt.Println("getEntriesFromCache range", key.Name, route.Prefix().String())
			route := route
			ipEntries = append(ipEntries, ipambev1alpha1.GetIPEntry(ctx, k, route.Prefix(), route.Labels()))
		}
	})

	return ipEntries, nil
}

func (r *bestore) deleteEntries(ctx context.Context, k store.Key) error {
	log := log.FromContext(ctx).With("key", k.String())

	ipEntries, err := r.listEntries(ctx, k)
	if err != nil {
		log.Error("cannot list entries", "error", err)
		return err
	}

	var errm error
	for _, ipEntry := range ipEntries {
		if err := r.client.Delete(ctx, ipEntry); err != nil {
			log.Error("cannot delete entry", "error", err)
			errm = errors.Join(errm, err)
			continue
		}
	}
	return errm
}

func (r *bestore) listEntries(ctx context.Context, k store.Key) ([]*ipambev1alpha1.IPEntry, error) {
	opt := []client.ListOption{
		//client.MatchingFields{
		//	"spec.networkInstance": k.Name,
		//},
	}

	ipEntries := ipambev1alpha1.IPEntryList{}
	if err := r.client.List(ctx, &ipEntries, opt...); err != nil {
		return nil, err
	}
	ipentries := []*ipambev1alpha1.IPEntry{}
	for _, ipEntry := range ipEntries.Items {
		ipEntry := ipEntry
		if ipEntry.Spec.NetworkInstance == k.Name {
			ipentries = append(ipentries, &ipEntry)
		}
	}

	return ipentries, nil
}

func (r *bestore) getNIPrefixes(ctx context.Context, k store.Key) (*ipamresv1alpha1.NetworkInstance, map[string]ipamresv1alpha1.Prefix, error) {
	ni := &ipamresv1alpha1.NetworkInstance{}
	if err := r.client.Get(ctx, k.NamespacedName, ni); err != nil {
		return nil, nil, err
	}
	niPrefixes := make(map[string]ipamresv1alpha1.Prefix)
	for _, prefix := range ni.Spec.Prefixes {
		niPrefixes[prefix.Prefix] = prefix
	}
	return ni, niPrefixes, nil
}

func (r *bestore) listClaims(ctx context.Context, k store.Key) (map[string]*ipambev1alpha1.IPClaim, error) {
	opt := []client.ListOption{
		/*
			client.MatchingFields{
				"spec.networkInstance": k.Name,
			},
		*/
	}

	claims := ipambev1alpha1.IPClaimList{}
	if err := r.client.List(ctx, &claims, opt...); err != nil {
		return nil, err
	}

	claimmap := map[string]*ipambev1alpha1.IPClaim{}
	for _, claim := range claims.Items {
		claim := claim
		if claim.Spec.NetworkInstance == k.Name {
			claimmap[(&claim).GetNamespacedName().String()] = &claim
		}

	}

	return claimmap, nil
}

func (r *bestore) listIPs(ctx context.Context, k store.Key) (map[string]*ipamresv1alpha1.IP, error) {
	opt := []client.ListOption{
		/*
			client.MatchingFields{
				"spec.networkInstance": k.Name,
			},
		*/
	}

	ipList := ipamresv1alpha1.IPList{}
	if err := r.client.List(ctx, &ipList, opt...); err != nil {
		return nil, err
	}

	ipmap := map[string]*ipamresv1alpha1.IP{}
	for _, ip := range ipList.Items {
		ip := ip
		if ip.Spec.NetworkInstance == k.Name {
			ipmap[(&ip).GetNamespacedName().String()] = &ip
		}
	}

	return ipmap, nil
}

func (r *bestore) restoreNiPrefixes(ctx context.Context, cacheCtx *CacheContext, ipEntries []*ipambev1alpha1.IPEntry, ni *ipamresv1alpha1.NetworkInstance, niPrefixes map[string]ipamresv1alpha1.Prefix) error {
	//log := log.FromContext(ctx)
	for i := len(ipEntries) - 1; i >= 0; i-- {
		ipEntry := ipEntries[i]
		if ipEntry.Spec.Owner.Group == ipamresv1alpha1.SchemeGroupVersion.Group &&
			ipEntry.Spec.Owner.Version == ipamresv1alpha1.SchemeGroupVersion.Version &&
			ipEntry.Spec.Owner.Kind == ipamresv1alpha1.NetworkInstanceKind {

			niPrefix, ok := niPrefixes[ipEntry.Spec.Prefix]
			if ok {
				claim, err := ni.GetIPClaim(niPrefix)
				if err != nil {
					return nil
				}
				if err := r.restoreClaim(ctx, cacheCtx, claim); err != nil {
					return err
				}
				// remove the entry since it is processed
				ipEntries = append(ipEntries[:i], ipEntries[i+1:]...)
				delete(niPrefixes, ipEntry.Spec.Prefix)
			}
		}
	}
	return nil
}

func (r *bestore) restoreIPs(ctx context.Context, cacheCtx *CacheContext, ipEntries []*ipambev1alpha1.IPEntry, claimType ipambev1alpha1.IPClaimType, ipmap map[string]*ipamresv1alpha1.IP) error {

	for i := len(ipEntries) - 1; i >= 0; i-- {
		ipEntry := ipEntries[i]
		if ipEntry.Spec.Owner.Group == ipamresv1alpha1.SchemeGroupVersion.Group &&
			ipEntry.Spec.Owner.Version == ipamresv1alpha1.SchemeGroupVersion.Version &&
			ipEntry.Spec.Owner.Kind == ipamresv1alpha1.IPKind {

			if claimType == ipEntry.Spec.ClaimType {
				nsn := types.NamespacedName{Namespace: ipEntry.Spec.Owner.Namespace, Name: ipEntry.Spec.Owner.Name}

				ip, ok := ipmap[nsn.String()]
				if ok {
					claim, err := ip.GetIPClaim()
					if err != nil {
						return nil
					}
					if err := r.restoreClaim(ctx, cacheCtx, claim); err != nil {
						return err
					}
					// remove the entry since it is processed
					ipEntries = append(ipEntries[:i], ipEntries[i+1:]...)
					delete(ipmap, nsn.String()) // delete the entry to optimize
				}
			}
		}
	}
	return nil
}

func (r *bestore) restoreIPClaims(ctx context.Context, cacheCtx *CacheContext, ipEntries []*ipambev1alpha1.IPEntry, claimType ipambev1alpha1.IPClaimType, ipclaimmap map[string]*ipambev1alpha1.IPClaim) error {

	for i := len(ipEntries) - 1; i >= 0; i-- {
		ipEntry := ipEntries[i]
		if ipEntry.Spec.Owner.Group == ipambev1alpha1.SchemeGroupVersion.Group &&
			ipEntry.Spec.Owner.Version == ipambev1alpha1.SchemeGroupVersion.Version &&
			ipEntry.Spec.Owner.Kind == ipambev1alpha1.IPClaimKind {

			if claimType == ipEntry.Spec.ClaimType {
				nsn := types.NamespacedName{Namespace: ipEntry.Spec.Owner.Namespace, Name: ipEntry.Spec.Owner.Name}

				claim, ok := ipclaimmap[nsn.String()]
				if ok {
					if err := r.restoreClaim(ctx, cacheCtx, claim); err != nil {
						return err
					}
					// remove the entry since it is processed
					ipEntries = append(ipEntries[:i], ipEntries[i+1:]...)
					delete(ipclaimmap, nsn.String()) // delete the entry to optimize
				}
			}
		}
	}
	return nil
}

func (r *bestore) restoreClaim(ctx context.Context, cacheCtx *CacheContext, claim *ipambev1alpha1.IPClaim) error {
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
