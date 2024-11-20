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
	"errors"
	"fmt"

	"github.com/henderiw/store"
	"github.com/kform-dev/choreo/apis/condition"
	"github.com/kuidio/kuid/apis/backend"
)

type rangeApplicator struct {
	name string
	applicator
}

// when a range changes the start and stop we delete the range
// and recreate it. All the children will be deleted as well
func (r *rangeApplicator) Validate(ctx context.Context, claim backend.ClaimObject) error {
	//log := log.FromContext(ctx)

	// reclaimRange gets the existing entries based on owner
	// -> 3 scenarios: none exist, they all exist, some exist
	// -> none exist -> claim them
	// -> exists, no change -> don't do much other than updating the labels
	// -> exists, change -> check childs and parents; when a parent exists and is from a different claim we stop, if children exist we block
	// have been deleted
	changed, err := r.validateChange(ctx, claim)
	if err != nil {
		return err
	}
	if changed {
		// we need to validate if there are no children
		if err := r.validateRangeOverlap(ctx, claim); err != nil {
			return err
		}
	}
	return nil
}

// validateChange checks if the range changed; change is only reported when the range existed
func (r *rangeApplicator) validateChange(ctx context.Context, claim backend.ClaimObject) (bool, error) {
	_, newClaimSet, err := claim.GetClaimSet(r.cacheInstanceCtx.Type())
	if err != nil {
		return false, err
	}

	_, oldClaimSet, err := r.getExistingCLaimSet(ctx, claim)
	if err != nil {
		return false, err
	}

	newEntries := newClaimSet.Difference(oldClaimSet)
	deletedEntries := oldClaimSet.Difference(newClaimSet)

	if len(newEntries.UnsortedList()) != 0 || len(deletedEntries.UnsortedList()) != 0 {
		// changed
		return true, nil
	}
	// no change
	return false, nil
}

func (r *rangeApplicator) validateRangeOverlap(_ context.Context, claim backend.ClaimObject) error {
	arange, err := claim.GetRangeID(r.cacheInstanceCtx.Type())
	if err != nil {
		return err
	}
	var errm error
	for _, id := range arange.IDs() {
		childEntries := r.cacheInstanceCtx.tree.Children(id)
		if len(childEntries) != 0 {
			errm = errors.Join(errm, fmt.Errorf("range id %s overlaps with children: %v", id.String(), childEntries))
		}
		parentEntries := r.cacheInstanceCtx.tree.Parents(id)
		if len(parentEntries) > 0 {
			for _, entry := range parentEntries {
				if !claim.IsOwner(entry.Labels()) {
					errm = errors.Join(errm, fmt.Errorf("range id %s overlaps with parent: %v", id.String(), parentEntries))
				}
			}
		}
	}
	return errm
}

func (r *rangeApplicator) Apply(ctx context.Context, claim backend.ClaimObject) error {
	arange, err := claim.GetRangeID(r.cacheInstanceCtx.Type())
	if err != nil {
		return err
	}
	newClaimMap, newClaimSet, err := claim.GetClaimSet(r.cacheInstanceCtx.Type())
	if err != nil {
		return err
	}

	oldClaimMap, oldClaimSet, err := r.getExistingCLaimSet(ctx, claim)
	if err != nil {
		return err
	}

	newEntries := newClaimSet.Difference(oldClaimSet)
	existingEntries := newClaimSet.Intersection(oldClaimSet)
	deletedEntries := oldClaimSet.Difference(newClaimSet)

	for idstr := range deletedEntries {
		if err := r.cacheInstanceCtx.tree.ReleaseID(oldClaimMap[idstr]); err != nil {
			return err
		}
	}
	for idstr := range newEntries {
		if err := r.cacheInstanceCtx.tree.ClaimID(newClaimMap[idstr], claim.GetClaimLabels()); err != nil {
			return err
		}
	}
	for idstr := range existingEntries {
		if err := r.cacheInstanceCtx.tree.Update(newClaimMap[idstr], claim.GetClaimLabels()); err != nil {
			return err
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
