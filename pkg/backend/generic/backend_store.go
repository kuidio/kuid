package generic

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/henderiw/idxtable/pkg/table"
	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	bebackend "github.com/kuidio/kuid/pkg/backend"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
)

func (r *be) restore(ctx context.Context, index backend.IndexObject) error {
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

	if err := r.restoreClaims(ctx, cacheInstanceCtx, curEntries, r.indexKind, backend.ClaimType_Range, claimmap); err != nil {
		return err
	}
	if err := r.restoreClaims(ctx, cacheInstanceCtx, curEntries, r.indexKind, backend.ClaimType_StaticID, claimmap); err != nil {
		return err
	}
	if err := r.restoreClaims(ctx, cacheInstanceCtx, curEntries, r.claimKind, backend.ClaimType_Range, claimmap); err != nil {
		return err
	}
	if err := r.restoreClaims(ctx, cacheInstanceCtx, curEntries, r.claimKind, backend.ClaimType_StaticID, claimmap); err != nil {
		return err
	}
	if err := r.restoreClaims(ctx, cacheInstanceCtx, curEntries, r.claimKind, backend.ClaimType_DynamicID, claimmap); err != nil {
		return err
	}

	log.Debug("restore entries left", "items", len(curEntries))

	return nil
}

func (r *be) saveAll(ctx context.Context, k store.Key) error {
	log := log.FromContext(ctx)
	log.Debug("SaveAll")

	cacheEntries, err := r.getEntriesFromCache(ctx, k)
	if err != nil {
		return err
	}


	apiEntries, err := r.listEntries(ctx, k)
	if err != nil {
		return err
	}

	for _, cacheEntry := range cacheEntries {
		found := false
		var oldEntry backend.EntryObject
		for idx, apiEntry := range apiEntries {
			if apiEntry.GetNamespacedName() == cacheEntry.GetNamespacedName() {
				// delete the current entry
				apiEntries = append(apiEntries[:idx], apiEntries[idx+1:]...)
				found = true
				oldEntry = apiEntry
				break
			}
		}

		if !found {
			if err := r.bestorage.CreateEntry(ctx, cacheEntry); err != nil {
				log.Error("saveAll create failed", "name", cacheEntry.GetName(), "error", err.Error())
				return err
			}
			continue
		}
		if err := r.bestorage.UpdateEntry(ctx, cacheEntry, oldEntry); err != nil {
			log.Error("saveAll update failed", "name", cacheEntry.GetName(), "error", err.Error())
			return err
		}
	}

	for _, apiEntry := range apiEntries {
		if err := r.bestorage.DeleteEntry(ctx, apiEntry); err != nil {
			log.Error("saveAll delete failed", "name", apiEntry.GetName(), "error", err.Error())
			return err
		}
	}
	return nil
}

// Destroy removes the store db
func (r *be) destroy(ctx context.Context, k store.Key) error {
	// no need to delete the index as this is what this fn is supposed to do
	if err := r.deleteClaims(ctx, k); err != nil {
		return err
	}
	return r.deleteEntries(ctx, k)
}

func (r *be) getEntriesFromCache(ctx context.Context, k store.Key) ([]backend.EntryObject, error) {
	//log := log.FromContext(ctx).With("key", k.String())

	cacheInstanceCtx, err := r.cache.Get(ctx, k)
	if err != nil {
		return nil, fmt.Errorf("cache index not initialized")
	}

	entries := make([]backend.EntryObject, 0, cacheInstanceCtx.Size())
	// add the main rib entry
	for _, entry := range cacheInstanceCtx.tree.GetAll() {
		entries = append(entries, r.entryFromCacheFn(k, "", entry.ID().String(), entry.Labels()))
	}
	// add all the range entries
	cacheInstanceCtx.ranges.List(func(key store.Key, t table.Table) {
		for _, entry := range t.GetAll() {
			entries = append(entries, r.entryFromCacheFn(k, key.Name, entry.ID().String(), entry.Labels()))
		}
	})

	return entries, nil
}

func (r *be) deleteEntries(ctx context.Context, k store.Key) error {
	log := log.FromContext(ctx).With("key", k.String())

	entries, err := r.listEntries(ctx, k)
	if err != nil {
		log.Error("cannot list entries", "error", err)
		return err
	}

	var errm error
	for _, curEntry := range entries {
		if err := r.bestorage.DeleteEntry(ctx, curEntry); err != nil {
			log.Error("saveAll delete failed", "name", curEntry.GetName(), "error", err.Error())
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

func (r *be) listEntries(ctx context.Context, k store.Key) ([]backend.EntryObject, error) {
	return r.bestorage.ListEntries(ctx, k)
}

func (r *be) listClaims(ctx context.Context, k store.Key) (map[string]backend.ClaimObject, error) {
	return r.bestorage.ListClaims(ctx, k)
}

func (r *be) restoreClaims(ctx context.Context, cacheInstanceCtx *CacheInstanceContext, entries []backend.EntryObject, kind string, claimType backend.ClaimType, claimmap map[string]backend.ClaimObject) error {
	log := log.FromContext(ctx)
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		if (kind == r.indexKind && entry.IsIndexEntry() && claimType == entry.GetClaimType()) ||
			(kind != r.indexKind && !entry.IsIndexEntry() && claimType == entry.GetClaimType()) {
			claimName := ""
			if len(entry.GetOwnerReferences()) > 0 {
				claimName = entry.GetOwnerReferences()[0].Name
			}
			nsn := types.NamespacedName{Namespace: entry.GetNamespace(), Name: claimName}
			claim, ok := claimmap[nsn.String()]
			if ok {
				log.Debug("restore claim", "kind", kind, "claimType", claimType, "claim", claim)
				if err := r.restoreClaim(ctx, cacheInstanceCtx, claim); err != nil {
					return err
				}
				// remove the entry since it is processed
				entries = append(entries[:i], entries[i+1:]...)
				delete(claimmap, nsn.String()) // delete the entry to optimize
			}
		}

		for _, ownerref := range entry.GetOwnerReferences() {
			if ownerref.Kind == r.claimKind {
				if claimType == entry.GetClaimType() {
					nsn := types.NamespacedName{Namespace: entry.GetNamespace(), Name: ownerref.Name}
					claim, ok := claimmap[nsn.String()]
					if ok {
						if err := r.restoreClaim(ctx, cacheInstanceCtx, claim); err != nil {
							return err
						}
						// remove the entry since it is processed
						entries = append(entries[:i], entries[i+1:]...)
						delete(claimmap, nsn.String()) // delete the entry to optimize
					}
				}
			}
		}
	}
	return nil
}

func (r *be) restoreClaim(ctx context.Context, cacheInstanceCtx *CacheInstanceContext, claim backend.ClaimObject) error {
	ctx = bebackend.InitClaimContext(ctx, "restore", claim)
	a, err := getApplicator(ctx, cacheInstanceCtx, claim)
	if err != nil {
		return err
	}
	
	// validate is needed, mainly for addresses since the parent route determines
	// e.g. the fact the address belongs to a range or not
	errList := claim.ValidateSyntax(cacheInstanceCtx.Type())
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

func (r *be) updateIndexClaims(ctx context.Context, index backend.IndexObject) error {
	log := log.FromContext(ctx)
	log.Debug("updateIPIndexClaims", "key", index.GetKey().String())
	key := index.GetKey()

	newClaims := index.GetClaims()

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
				log.Error("updateIndexClaims create failed", "name", newClaim.GetName(), "error", err.Error())
				errm = errors.Join(errm, err)
				continue
			}
			continue
		}
		if err := r.bestorage.UpdateClaim(ctx, newClaim, oldClaim); err != nil {
			log.Error("updateIndexClaims create failed", "name", newClaim.GetName(), "error", err.Error())
			errm = errors.Join(errm, err)
			continue
		}
		delete(existingClaims, newClaim.GetNamespacedName().String())
	}

	for _, claim := range existingClaims {
		log.Debug("updateIndexClaims: delete existing claims", "claim", claim.GetName())
		if err := r.bestorage.DeleteClaim(ctx, claim); err != nil {
			log.Error("updateIndexClaims delete failed", "name", claim.GetName(), "error", err.Error())
			errm = errors.Join(errm, err)
			continue
		}
	}
	if errm != nil {
		return errm
	}
	return r.saveAll(ctx, key)
}

func EntryTransformer(_ context.Context, newObj runtime.Object, oldObj runtime.Object) (runtime.Object, error) {
	// Type assertion to specific object types, assuming we are working with a type that has Spec and Status fields
	new, ok := newObj.(backend.EntryObject)
	if !ok {
		return nil, fmt.Errorf("newObj is not of type EntryObject")
	}
	old, ok := oldObj.(backend.EntryObject)
	if !ok {
		return nil, fmt.Errorf("oldObj is not of type EntryObject")
	}

	new.SetResourceVersion(old.GetResourceVersion())
	new.SetUID(old.GetUID())

	return new, nil
}

func ClaimTransformer(_ context.Context, newObj runtime.Object, oldObj runtime.Object) (runtime.Object, error) {
	// Type assertion to specific object types, assuming we are working with a type that has Spec and Status fields
	new, ok := newObj.(backend.ClaimObject)
	if !ok {
		return nil, fmt.Errorf("newObj is not of type ClaimObject, got: %v", reflect.TypeOf(newObj).Name())
	}
	old, ok := oldObj.(backend.ClaimObject)
	if !ok {
		return nil, fmt.Errorf("oldObj is not of type ClaimObject, got: %v", reflect.TypeOf(newObj).Name())
	}

	new.SetResourceVersion(old.GetResourceVersion())
	new.SetUID(old.GetUID())

	return new, nil
}

func (r *be) listIndexClaims(ctx context.Context, k store.Key) (map[string]backend.ClaimObject, error) {
	return r.bestorage.ListClaims(ctx, k, &ListOptions{
		OwnerKind: r.indexKind,
	})
}
