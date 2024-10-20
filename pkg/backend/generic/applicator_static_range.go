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

package generic

import (
	"context"
	"fmt"

	"github.com/henderiw/store"
	"github.com/kform-dev/choreo/apis/condition"
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
	// -> same exist -> delete them including the children and reclaim the new range after the entries
	// have been deleted
	exists, err := r.validateExists(ctx, claim)
	if err != nil {
		return err
	}
	if !exists {
		// we need to validate if there are no children
		if err := r.validateRangeOverlap(ctx, claim); err != nil {
			return err
		}
	}
	r.rangeExists = exists
	return nil
}

// reclaimRange gets the existing entries based on owner
// -> 3 scenarios: none exist, they all exist, some exist
// -> none exist -> claim them
// -> they all exist -> reclaim them
// -> same exist -> delete them including the children and reclaim the new range aftre the entries
// have been deleted
func (r *rangeApplicator) validateExists(ctx context.Context, claim backend.ClaimObject) (bool, error) {
	arange, err := claim.GetRangeID(r.cacheInstanceCtx.Type())
	if err != nil {
		return false, err
	}
	rangeExists := true // we are optimistic and set the claimset to true since we have entries
	claimSet := map[string]struct{}{}
	for _, rangeID := range arange.IDs() {
		claimSet[rangeID.String()] = struct{}{}
	}

	existingEntries, err := r.getEntriesByOwner(ctx, claim)
	if err != nil {
		return false, err
	}
	// delete the
	for treeName, existingEntries := range existingEntries {
		if treeName != "" {
			return false, fmt.Errorf("cannot have a range in non root tree: %s", treeName)
		}
		for _, existingEntry := range existingEntries {
			delete(claimSet, existingEntry.ID().String())
		}
	}

	if len(claimSet) != 0 {
		// cleanup
		rangeExists = false
		// remove all entries as the range change
		if err := r.deleteNonClaimedEntries(ctx, existingEntries, nil, ""); err != nil {
			return false, err
		}
		k := store.ToKey(claim.GetName())
		if _, err := r.cacheInstanceCtx.ranges.Get(k); err == nil {
			// exists
			if err := r.cacheInstanceCtx.ranges.Delete(k); err != nil {
				return false, err
			}
		}
	}
	// all good -> they either all exist or none exists or we cleaned up
	return rangeExists, nil
}

func (r *rangeApplicator) validateRangeOverlap(_ context.Context, claim backend.ClaimObject) error {
	arange, err := claim.GetRangeID(r.cacheInstanceCtx.Type())
	if err != nil {
		return err
	}
	for _, id := range arange.IDs() {
		entry, err := r.cacheInstanceCtx.tree.Get(id)
		if err == nil {
			// this shouls always fail since the range existance was already validated
			labels := entry.Labels()
			if err := claim.ValidateOwner(labels); err != nil {
				return err
			}
		}
		childEntries := r.cacheInstanceCtx.tree.Children(id)
		if len(childEntries) != 0 {
			return fmt.Errorf("range overlaps with children: %v", childEntries)
		}
		parentEntries := r.cacheInstanceCtx.tree.Parents(id)
		if len(parentEntries) > 0 {
			return fmt.Errorf("range overlaps with parent: %v", parentEntries)
		}
	}
	return nil
}

func (r *rangeApplicator) Apply(ctx context.Context, claim backend.ClaimObject) error {
	arange, err := claim.GetRangeID(r.cacheInstanceCtx.Type())
	if err != nil {
		return err
	}
	for _, id := range arange.IDs() {
		if r.rangeExists {
			if err := r.cacheInstanceCtx.tree.Update(id, claim.GetClaimLabels()); err != nil {
				return err
			}
		} else {
			if err := r.cacheInstanceCtx.tree.ClaimID(id, claim.GetClaimLabels()); err != nil {
				return err
			}
		}

	}
	k := store.ToKey(claim.GetName())
	if _, err := r.cacheInstanceCtx.ranges.Get(k); err != nil {
		//table := table.New(uint32(arange.From().ID()), uint32(arange.To().ID()))
		table := claim.GetTable(r.cacheInstanceCtx.Type(), arange.From().ID(), arange.To().ID())
		if err := r.cacheInstanceCtx.ranges.Create(k, table); err != nil {
			return err
		}
	}
	claim.SetStatusRange(claim.GetRange())
	claim.SetConditions(condition.Ready())
	return nil
}

func (r *rangeApplicator) Delete(ctx context.Context, claim backend.ClaimObject) error {
	return r.delete(ctx, claim)
}
