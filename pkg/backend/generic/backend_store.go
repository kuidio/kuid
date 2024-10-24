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
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
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

	if err := r.restoreMinMaxRanges(ctx, cacheInstanceCtx, curEntries, index); err != nil {
		return err
	}

	if err := r.restoreClaims(ctx, cacheInstanceCtx, curEntries, backend.ClaimType_Range, claimmap); err != nil {
		return err
	}
	if err := r.restoreClaims(ctx, cacheInstanceCtx, curEntries, backend.ClaimType_StaticID, claimmap); err != nil {
		return err
	}
	if err := r.restoreClaims(ctx, cacheInstanceCtx, curEntries, backend.ClaimType_DynamicID, claimmap); err != nil {
		return err
	}

	log.Debug("restore entries left", "items", len(curEntries))

	return nil
}

func (r *be) saveAll(ctx context.Context, k store.Key) error {
	log := log.FromContext(ctx)
	log.Debug("SaveAll")

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
		var oldEntry backend.EntryObject
		for idx, curEntry := range curEntries {
			idx := idx
			curEntry := curEntry
			if curEntry.GetNamespacedName() == newEntry.GetNamespacedName() {
				// delete the current entry
				curEntries = append(curEntries[:idx], curEntries[idx+1:]...)
				found = true
				oldEntry = curEntry
				break
			}
		}

		ctx = genericapirequest.WithNamespace(ctx, newEntry.GetNamespace())
		if !found {
			if _, err := r.entryStorage.Create(ctx, newEntry, nil, &metav1.CreateOptions{
				FieldManager: "backend",
			}); err != nil {
				log.Error("saveAll create failed", "name", newEntry.GetName(), "error", err.Error())
				return err
			}
			continue
		}

		defaultObjInfo := rest.DefaultUpdatedObjectInfo(oldEntry, entryTransformer)
		if _, _, err := r.entryStorage.Update(ctx, oldEntry.GetName(), defaultObjInfo, nil, nil, false, &metav1.UpdateOptions{
			FieldManager: "backend",
		}); err != nil {
			fmt.Println("update err", err)
			return err
		}
	}

	for _, curEntry := range curEntries {
		ctx = genericapirequest.WithNamespace(ctx, curEntry.GetNamespace())
		if _, _, err := r.entryStorage.Delete(ctx, curEntry.GetName(), nil, &metav1.DeleteOptions{}); err != nil {
			return err
		}
	}
	return nil
}

// Destroy removes the store db
func (r *be) destroy(ctx context.Context, k store.Key) error {
	// no need to delete the index as this is what this fn is supposed to do
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
		entry := entry
		entries = append(entries, r.entryFromCacheFn(k, "", entry.ID().String(), entry.Labels()))
	}
	// add all the range entries
	cacheInstanceCtx.ranges.List(func(key store.Key, t table.Table) {
		for _, entry := range t.GetAll() {
			entry := entry
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
	for _, entry := range entries {
		ctx = genericapirequest.WithNamespace(ctx, entry.GetNamespace())
		if _, _, err := r.entryStorage.Delete(ctx, entry.GetName(), nil, &metav1.DeleteOptions{}); err != nil {
			log.Error("cannot delete entry", "error", err)
			errm = errors.Join(errm, err)
			continue
		}
	}
	return errm
}

func (r *be) listEntries(ctx context.Context, k store.Key) ([]backend.EntryObject, error) {
	log := log.FromContext(ctx).With("key", k.String())
	/*
			selector, err := selector.ExprSelectorAsSelector(
				&selectorv1alpha1.ExpressionSelector{
					Match: map[string]string{
						"spec.index": k.Name,
					},
				},
			)
		if err != nil {
			return nil, err
		}
	*/
	list, err := r.entryStorage.List(ctx, &internalversion.ListOptions{})
	if err != nil {
		return nil, err
	}

	items, err := meta.ExtractList(list)
	if err != nil {
		return nil, err
	}

	entryList := make([]backend.EntryObject, 0)
	var errm error
	for _, obj := range items {
		entryObj, ok := obj.(backend.EntryObject)
		if !ok {
			log.Error("obj is not an EntryObject", "obj", reflect.TypeOf(obj).Name())
			errm = errors.Join(errm, err)
			continue
		}
		if entryObj.GetIndex() == k.Name {
			entryList = append(entryList, entryObj)
		}
	}
	return entryList, errm
}

func (r *be) listClaims(ctx context.Context, k store.Key) (map[string]backend.ClaimObject, error) {
	log := log.FromContext(ctx).With("key", k.String())
	/*
		selector, err := selector.ExprSelectorAsSelector(
			&selectorv1alpha1.ExpressionSelector{
				Match: map[string]string{
					"spec.index": k.Name,
				},
			},
		)
		if err != nil {
			return nil, err
		}
	*/
	list, err := r.claimStorage.List(ctx, &internalversion.ListOptions{})
	if err != nil {
		return nil, err
	}

	items, err := meta.ExtractList(list)
	if err != nil {
		return nil, err
	}

	claimMap := make(map[string]backend.ClaimObject)
	var errm error
	for _, obj := range items {
		claimObj, ok := obj.(backend.ClaimObject)
		if !ok {
			log.Error("obj is not an ClaimObject", "obj", reflect.TypeOf(obj).Name())
			errm = errors.Join(errm, err)
			continue
		}
		claimMap[claimObj.GetNamespacedName().String()] = claimObj
	}
	return claimMap, errm
}

func (r *be) restoreMinMaxRanges(ctx context.Context, cacheInstanceCtx *CacheInstanceContext, entries []backend.EntryObject, index backend.IndexObject) error {
	storedEntries := sets.New[string]()
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		for _, ownerref := range entry.GetOwnerReferences() {
			if ownerref.APIVersion == index.GetObjectKind().GroupVersionKind().GroupVersion().Identifier() &&
				ownerref.Kind == index.GetObjectKind().GroupVersionKind().Kind &&
				ownerref.Name == index.GetName() &&
				ownerref.UID == index.GetUID() {
				entries = append(entries[:i], entries[i+1:]...)
				storedEntries.Insert(entry.GetSpecID())
			}
		}
	}

	if index.GetMinID() != nil && *index.GetMinID() != 0 {
		claim := index.GetMinClaim()
		if err := r.restoreClaim(ctx, cacheInstanceCtx, claim); err != nil {
			return err
		}
	}
	if index.GetMaxID() != nil && *index.GetMaxID() != index.GetMax() {
		claim := index.GetMaxClaim()
		if err := r.restoreClaim(ctx, cacheInstanceCtx, claim); err != nil {
			return err
		}
	}
	// At init when there is no entries initialized this allows to store the entries in the database
	if storedEntries.Len() == 0 {
		entries, err := r.getEntriesFromCache(ctx, index.GetKey())
		if err != nil {
			return err
		}
		for _, entry := range entries {
			uobj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(entry)
			if err != nil {
				return err
			}
			u := &unstructured.Unstructured{
				Object: uobj,
			}
			ctx = genericapirequest.WithNamespace(ctx, u.GetNamespace())

			if _, err := r.entryStorage.Create(ctx, u, nil, &metav1.CreateOptions{
				FieldManager: "backend",
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *be) restoreClaims(ctx context.Context, cacheInstanceCtx *CacheInstanceContext, entries []backend.EntryObject, claimType backend.ClaimType, claimmap map[string]backend.ClaimObject) error {
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
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


func entryTransformer(_ context.Context, newObj runtime.Object, oldObj runtime.Object) (runtime.Object, error) {
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