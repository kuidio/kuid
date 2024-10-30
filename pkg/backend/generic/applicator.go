package generic

import (
	"context"
	"errors"
	"fmt"

	"github.com/henderiw/idxtable/pkg/table"
	"github.com/henderiw/idxtable/pkg/tree"
	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/ptr"
)

type Applicator interface {
	Validate(ctx context.Context, claim backend.ClaimObject) error
	Apply(ctx context.Context, claim backend.ClaimObject) error
	Delete(ctx context.Context, claim backend.ClaimObject) error
}

type applicator struct {
	cacheInstanceCtx *CacheInstanceContext
}

func (r *applicator) getEntriesByOwner(_ context.Context, claim backend.ClaimObject) (map[string]tree.Entries, error) {
	treeEntries := map[string]tree.Entries{}
	ownerSelector, err := claim.GetOwnerSelector()
	if err != nil {
		return nil, err
	}
	claimType := claim.GetClaimType()
	// treeEntries with empty string is the root tree
	treeEntries[""] = r.cacheInstanceCtx.tree.GetByLabel(ownerSelector)
	if claimType == backend.ClaimType_Range {
		// for range claims we return
		return treeEntries, nil
	}
	// this is a NOT a range claim -> static or dynamic claim
	if len(treeEntries) != 0 && len(treeEntries[""]) > 1 {
		return treeEntries, fmt.Errorf("multiple entries match the owner, %v", treeEntries[""])
	}
	var errs error
	r.cacheInstanceCtx.ranges.List(func(k store.Key, t table.Table) {
		treeEntries[k.Name] = t.GetByLabel(ownerSelector)
		if len(treeEntries[k.Name]) > 1 {
			errs = errors.Join(errs, fmt.Errorf("multiple entries match the owner, %v", treeEntries[k.Name]))
			return
		}
	})
	if errs != nil {
		return nil, errs
	}
	return treeEntries, nil
}

func (r *applicator) delete(ctx context.Context, claim backend.ClaimObject) error {
	existingEntries, err := r.getEntriesByOwner(ctx, claim)
	if err != nil {
		return err
	}

	for treeName, existingEntries := range existingEntries {
		for _, existingEntry := range existingEntries {
			if treeName == "" {
				r.cacheInstanceCtx.tree.ReleaseID(existingEntry.ID())
			} else {
				k := store.ToKey(treeName)
				if table, err := r.cacheInstanceCtx.ranges.Get(k); err == nil {
					if err := table.Release(existingEntry.ID().ID()); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (r *applicator) getEntriesByLabelSelector(ctx context.Context, claim backend.ClaimObject) tree.Entries {
	log := log.FromContext(ctx)
	labelSelector, err := claim.GetLabelSelector()
	if err != nil {
		log.Error("cannot get label selector", "error", err.Error())
		return nil
	}
	return r.cacheInstanceCtx.tree.GetByLabel(labelSelector)
}

func reclaimIDFromExisitingEntries(existingEntries map[string]tree.Entries, id uint64) (*uint64, string) {
	for treeName, existingEntries := range existingEntries {
		for _, existingEntry := range existingEntries {
			if id == existingEntry.ID().ID() {
				return &id, treeName
			}
		}
	}
	return nil, ""
}

func claimIDFromExisitingEntries(existingEntries map[string]tree.Entries) (*uint64, string) {
	for treeName, existingEntries := range existingEntries {
		for _, existingEntry := range existingEntries {
			return ptr.To[uint64](existingEntry.ID().ID()), treeName
		}
	}
	return nil, ""
}

func (r *applicator) deleteNonClaimedEntries(_ context.Context, existingEntries map[string]tree.Entries, id *uint64, reclaimTreeName string) error {
	for treeName, existingEntries := range existingEntries {
		for _, existingEntry := range existingEntries {
			if id != nil && *id == existingEntry.ID().ID() && reclaimTreeName == treeName {
				continue
			}
			if treeName == "" {
				r.cacheInstanceCtx.tree.ReleaseID(existingEntry.ID())
			} else {
				k := store.ToKey(treeName)
				if table, err := r.cacheInstanceCtx.ranges.Get(k); err == nil {
					if err := table.Release(existingEntry.ID().ID()); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func isReserved(parentName, index string) bool {
	// given we use the ownerreference with index in the kind the parentName and index
	// match when we check for reserved fields
	// This means a range cannot be defined using the name of the index
	return parentName == index
	/*
		return parentName == fmt.Sprintf("%s.%s", index, backend.IndexReservedMaxName) ||
			parentName == fmt.Sprintf("%s.%s", index, backend.IndexReservedMinName)
	*/
}

func (r *applicator) getExistingCLaimSet(ctx context.Context, claim backend.ClaimObject) (sets.Set[tree.ID], error) {
	oldClaimIDSet := sets.New[tree.ID]()

	existingEntries, err := r.getEntriesByOwner(ctx, claim)
	if err != nil {
		return nil, err
	}
	// delete the entries from the claimSet that overlap
	for treeName, existingEntries := range existingEntries {
		if treeName != "" {
			return nil, fmt.Errorf("cannot have a range in non root tree: %s", treeName)
		}
		for _, existingEntry := range existingEntries {
			oldClaimIDSet.Insert(existingEntry.ID())
		}
	}
	return oldClaimIDSet, nil
}
