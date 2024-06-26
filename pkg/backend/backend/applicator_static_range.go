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

package backend

import (
	"context"
	"fmt"

	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
)

type rangeApplicator struct {
	name string
	applicator
	rangeExists bool
}

// when a range changes the start and stop we delete the range
// and recreate it. All the children will be deleted as well
func (r *rangeApplicator) Validate(ctx context.Context, claim backend.ClaimObject) error {
	//log := log.FromContext(ctx)

	// reclaimRange gets the existing entries based on owner
	// -> 3 scenarios: none exist, they all exist, some exist
	// -> none exist -> claim them
	// -> they all exist -> reclaim them
	// -> same exist -> delete them including the children and reclaim the new range aftre the entries
	// have been deleted
	if err := r.reclaimRange(ctx, claim); err != nil {
		return err
	}
	if !r.rangeExists {
		// we need to validate if there are no children
		if err := r.validateRangeOverlap(ctx, claim); err != nil {
			return err
		}
	}

	return nil
}

// reclaimRange gets the existing entries based on owner
// -> 3 scenarios: none exist, they all exist, some exist
// -> none exist -> claim them
// -> they all exist -> reclaim them
// -> same exist -> delete them including the children and reclaim the new range aftre the entries
// have been deleted
func (r *rangeApplicator) reclaimRange(ctx context.Context, claim backend.ClaimObject) error {
	//if claim.Spec.Range == nil {
	//	return fmt.Errorf("cannot claim a range w/o a range")
	//}
	arange, err := claim.GetRangeID(r.cacheCtx.Type())
	if err != nil {
		return err
	}
	claimSet := map[string]struct{}{}
	for _, rangeID := range arange.IDs() {
		r.rangeExists = true // we are optimistic and set the claimset to true since we have entries
		claimSet[rangeID.String()] = struct{}{}
	}
	//fmt.Println("claimSet before", claim.GetName(), claimSet)

	existingEntries, err := r.getEntriesByOwner(ctx, claim)
	if err != nil {
		return err
	}
	//fmt.Println("existingEntries", existingEntries)
	// delete the
	for treeName, existingEntries := range existingEntries {
		if treeName != "" {
			return fmt.Errorf("cannot have a range in non root tree: %s", treeName)
		}
		for _, existingEntry := range existingEntries {
			//fmt.Println("existingEntry", existingEntry.ID().String())
			delete(claimSet, existingEntry.ID().String())
		}
	}

	//fmt.Println("claimSet after", claim.GetName(), claimSet)
	if len(claimSet) != 0 {
		// cleanup
		r.rangeExists = false
		// remove all entries as the range change
		if err := r.deleteNonClaimedEntries(ctx, existingEntries, nil, ""); err != nil {
			return err
		}
		k := store.ToKey(claim.GetName())
		if _, err := r.cacheCtx.ranges.Get(ctx, k); err == nil {
			// exists
			if err := r.cacheCtx.ranges.Delete(ctx, k); err != nil {
				return err
			}
		}
	}
	// all good -> they either all exist or none exists or we cleaned up
	return nil
}

func (r *rangeApplicator) validateRangeOverlap(_ context.Context, claim backend.ClaimObject) error {
	arange, err := claim.GetRangeID(r.cacheCtx.Type())
	if err != nil {
		return err
	}
	for _, id := range arange.IDs() {
		entry, err := r.cacheCtx.tree.Get(id)
		if err == nil {
			// this shouls always fail since the range existance was already validated
			labels := entry.Labels()
			if err := claim.ValidateOwner(labels); err != nil {
				return err
			}
		}
		childEntries := r.cacheCtx.tree.Children(id)
		if len(childEntries) != 0 {
			return fmt.Errorf("range overlaps with children: %v", childEntries)
		}
		parentEntries := r.cacheCtx.tree.Parents(id)
		if len(parentEntries) > 0 {
			return fmt.Errorf("range overlaps with parent: %v", parentEntries)
		}
	}
	return nil
}

func (r *rangeApplicator) Apply(ctx context.Context, claim backend.ClaimObject) error {
	arange, err := claim.GetRangeID(r.cacheCtx.Type())
	if err != nil {
		return err
	}
	for _, id := range arange.IDs() {
		if r.rangeExists {
			if err := r.cacheCtx.tree.Update(id, claim.GetClaimLabels()); err != nil {
				return err
			}
		} else {
			if err := r.cacheCtx.tree.ClaimID(id, claim.GetClaimLabels()); err != nil {
				return err
			}
		}

	}
	k := store.ToKey(claim.GetName())
	if _, err := r.cacheCtx.ranges.Get(ctx, k); err != nil {
		//table := table.New(uint32(arange.From().ID()), uint32(arange.To().ID()))
		table := claim.GetTable(r.cacheCtx.Type(), arange.From().ID(), arange.To().ID())
		if err := r.cacheCtx.ranges.Create(ctx, k, table); err != nil {
			return err
		}
	}
	claim.SetStatusRange(claim.GetRange())
	return nil
}

func (r *rangeApplicator) Delete(ctx context.Context, claim backend.ClaimObject) error {
	return r.delete(ctx, claim)
}
