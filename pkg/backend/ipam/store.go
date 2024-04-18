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

	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	ipamresv1alpha1 "github.com/kuidio/kuid/apis/resource/ipam/v1alpha1"
	"github.com/kuidio/kuid/pkg/backend"
	"github.com/kuidio/kuid/pkg/reconcilers/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
	//log := log.FromContext(ctx).With("key", k.String())
	ipIndex := ipambev1alpha1.IPIndex{}
	if err := r.client.Get(ctx, k.NamespacedName, &ipIndex); err != nil {
		if resource.IgnoreNotFound(err) != nil {
			return err
		}
		// delete entries if they exist
		return r.deleteEntries(ctx, k)
	}
	return r.restoreEntries(ctx, k)
}

// only used in configmap
func (r *bestore) SaveAll(ctx context.Context, k store.Key) error {
	log := log.FromContext(ctx)
	ipIndex := ipambev1alpha1.IPIndex{}
	if err := r.client.Get(ctx, k.NamespacedName, &ipIndex); err != nil {
		if resource.IgnoreNotFound(err) != nil {
			return err
		}
		// TBD if this is an OK scenario or we should error
		if err := r.client.Create(ctx, ipambev1alpha1.BuildIPIndex(metav1.ObjectMeta{
			Namespace: k.Namespace,
			Name:      k.Name,
		}, nil, nil)); err != nil {
			return err
		}

	}
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
		found := false
		var ipEntry *ipambev1alpha1.IPEntry
		for idx, curIPEntry := range curIPEntries.Items {
			if curIPEntry.Namespace == newIPEntry.Namespace &&
				curIPEntry.Name == newIPEntry.Name {
				curIPEntries.Items = append(curIPEntries.Items[:idx], curIPEntries.Items[idx+1:]...)
				found = true
				ipEntry = &curIPEntry
				break
			}
		}
		if !found {
			if err := r.client.Create(ctx, newIPEntry); err != nil {
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
	for _, curIPEntry := range curIPEntries.Items {
		if err := r.client.Delete(ctx, &curIPEntry); err != nil {
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

	ipEntries := make([]*ipambev1alpha1.IPEntry, 0, cacheCtx.rib.Size())
	for _, route := range cacheCtx.rib.GetTable() {
		ipEntries = append(ipEntries, ipambev1alpha1.GetIPEntry(ctx, k, route.Prefix(), route.Labels()))
	}
	return ipEntries, nil
}

func (r *bestore) restoreEntries(ctx context.Context, k store.Key) error {
	//log := log.FromContext(ctx)

	// debug start
	/*
		ipEntries, err := r.listEntries(ctx, k)
		if err != nil {
			return err
		}
		//for _, ipEntry := range ipEntries.Items {
		//	log.Info("ip entry", "nsn", ipEntry.GetNamespacedName().String())
		//}
		// debug end
	*/

	cacheCtx, err := r.cache.Get(ctx, k, true) // return a rib even when not initialized
	if err != nil {
		return err
	}

	claims, err := r.listClaims(ctx, k)
	if err != nil {
		return err
	}

	// we restore in order right now
	// 1st network instance
	// 2nd prefixes
	// 3rd claims
	niGvk := schema.GroupVersionKind{
		Group:   ipamresv1alpha1.Group,
		Version: ipamresv1alpha1.Version,
		Kind:    ipamresv1alpha1.NetworkInstanceKind,
	}
	if err := r.restorePrefixes(ctx, claims, niGvk, cacheCtx); err != nil {
		return err
	}
	pfxGvk := schema.GroupVersionKind{
		Group:   ipamresv1alpha1.Group,
		Version: ipamresv1alpha1.Version,
		Kind:    ipamresv1alpha1.IPPrefixKind,
	}
	if err := r.restorePrefixes(ctx, claims, pfxGvk, cacheCtx); err != nil {
		return err
	}
	rangeGvk := schema.GroupVersionKind{
		Group:   ipamresv1alpha1.Group,
		Version: ipamresv1alpha1.Version,
		Kind:    ipamresv1alpha1.IPRangeKind,
	}
	if err := r.restorePrefixes(ctx, claims, rangeGvk, cacheCtx); err != nil {
		return err
	}
	addrGvk := schema.GroupVersionKind{
		Group:   ipamresv1alpha1.Group,
		Version: ipamresv1alpha1.Version,
		Kind:    ipamresv1alpha1.IPAddressKind,
	}
	if err := r.restorePrefixes(ctx, claims, addrGvk, cacheCtx); err != nil {
		return err
	}
	claimGvk := schema.GroupVersionKind{
		Group:   ipambev1alpha1.Group,
		Version: ipambev1alpha1.Version,
		Kind:    ipambev1alpha1.IPClaimKind,
	}
	if err := r.restorePrefixes(ctx, claims, claimGvk, cacheCtx); err != nil {
		return err
	}

	return nil
}

func (r *bestore) deleteEntries(ctx context.Context, k store.Key) error {
	log := log.FromContext(ctx).With("key", k.String())

	ipEntries, err := r.listEntries(ctx, k)
	if err != nil {
		log.Error("cannot list entries", "error", err)
		return err
	}

	var errm error
	for _, ipEntry := range ipEntries.Items {
		if err := r.client.Delete(ctx, &ipEntry); err != nil {
			log.Error("cannot delete entry", "error", err)
			errm = errors.Join(errm, err)
			continue
		}
	}
	return errm
}

func (r *bestore) listEntries(ctx context.Context, _ store.Key) (*ipambev1alpha1.IPEntryList, error) {
	opt := []client.ListOption{
		/*
			client.MatchingFields{
				"spec.networkInstance": k.Name,
			},
		*/
	}

	ipEntries := ipambev1alpha1.IPEntryList{}
	if err := r.client.List(ctx, &ipEntries, opt...); err != nil {
		return nil, err
	}

	return &ipEntries, nil
}

func (r *bestore) listClaims(ctx context.Context, _ store.Key) (*ipambev1alpha1.IPClaimList, error) {
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

	return &claims, nil
}

func (r *bestore) restorePrefixes(ctx context.Context, claims *ipambev1alpha1.IPClaimList, gvk schema.GroupVersionKind, cacheCtx *CacheContext) error {
	log := log.FromContext(ctx)
	for _, claim := range claims.Items {
		if claim.GetCondition(conditionv1alpha1.ConditionTypeReady).Status == metav1.ConditionTrue {
			if claim.Spec.Owner.Group == gvk.Group && claim.Spec.Owner.Version == gvk.Version && claim.Spec.Owner.Kind == gvk.Kind {
				a, err := getApplicator(ctx, &claim, cacheCtx)
				if err != nil {
					return err
				}
				log.Info("restore claim", "name", claim.Name, "prefix", claim.Status.Prefix)
				if err := a.Apply(ctx, &claim); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
