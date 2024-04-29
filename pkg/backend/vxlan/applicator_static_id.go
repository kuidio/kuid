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
	"fmt"

	"github.com/henderiw/idxtable/pkg/tree/id32"
	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	vxlanbev1alpha1 "github.com/kuidio/kuid/apis/backend/vxlan/v1alpha1"
)

type staticVXLANApplicator struct {
	name string
	applicator
	claimID        *uint32
	parentTreeName string
}

func (r *staticVXLANApplicator) Validate(ctx context.Context, claim *vxlanbev1alpha1.VXLANClaim) error {
	//log := log.FromContext(ctx)
	if claim.Spec.ID == nil {
		return fmt.Errorf("cannot claim a static id w/o an id")
	}

	if err := r.reclaimID(ctx, claim); err != nil {
		return err
	}

	if r.claimID == nil {
		// when the claim did not exist we need to check the
		// parent contect to know from which tree/table to
		// claim the ID
		if err := r.getParentContext(ctx, claim); err != nil {
			return err
		}
	}
	return nil
}

// reclaimID will validate if the id specified in the claim exists already.
// If so it will reclaim it and update the parentTreeName (this is to ensure the claim reclaims it from the proper tree/table)
// if no entry exist it will return an empty
// On top the entries that were not claimed are cleaned up, such that we delete entries that are
// void. E.g. this takes care of the fact that a use changed the static ID as the reclaim failed
// so the remaining entry is cleaned up
func (r *staticVXLANApplicator) reclaimID(ctx context.Context, claim *vxlanbev1alpha1.VXLANClaim) error {
	// reclaim ID
	existingEntries, err := r.getEntriesByOwner(ctx, claim)
	if err != nil {
		return err
	}

	claimID, claimTreeName := reclaimIDFromExisitingEntries(existingEntries, *claim.Spec.ID)
	fmt.Println("static id", claim.Name, claimID, claimTreeName)
	// remove the existing entries that done match the claimed ID
	// should be none, but just in case
	if err := r.deleteNonClaimedEntries(ctx, existingEntries, claimID, claimTreeName); err != nil {
		return err
	}

	r.claimID = claimID
	if claimID != nil {
		// here we are sure we got the same ID as the static ID
		// otherwise it would have been nil
		r.parentTreeName = claimTreeName
	}

	return nil
}

func (r *staticVXLANApplicator) getParentContext(ctx context.Context, claim *vxlanbev1alpha1.VXLANClaim) error {
	log := log.FromContext(ctx).With("name", claim.GetName())

	entry, err := r.cacheCtx.tree.Get(id32.NewID(*claim.Spec.ID, 32))
	if err == nil {
		// entry exists
		labels := entry.Labels()
		// a range can overlap so we return the entry as a parent if the entry match and it is a range
		claimType := vxlanbev1alpha1.GetClaimTypeFromString(labels[backend.KuidVXLANClaimTypeKey])
		if claimType == vxlanbev1alpha1.VXLANClaimType_Range {
			r.parentTreeName = labels[backend.KuidClaimNameKey]
			return nil
		} else {
			// This should always result in a different owner
			// since we checked the claimed entries before
			if err := claim.ValidateOwner(labels); err != nil {
				return err
			}
			return nil
		}
	}
	parentEntries := r.cacheCtx.tree.Parents(id32.NewID(*claim.Spec.ID, 32))
	if len(parentEntries) > 1 {
		log.Error("got multiple parent entries", "entries", parentEntries)
		return fmt.Errorf("multiple parent entries %v", parentEntries)
	}
	for _, parentEntry := range parentEntries {
		labels := parentEntry.Labels()
		parentClaimType := vxlanbev1alpha1.GetClaimTypeFromString(labels[backend.KuidVXLANClaimTypeKey])
		if parentClaimType == vxlanbev1alpha1.VXLANClaimType_Range {
			r.parentTreeName = labels[backend.KuidClaimNameKey]
			break
		} else {
			log.Error("got parent which is not a range", "entry", parentEntry)
			return fmt.Errorf("got parent which is not a range %s", labels[backend.KuidVXLANClaimTypeKey])
		}
	}
	return nil
}

func (r *staticVXLANApplicator) Apply(ctx context.Context, claim *vxlanbev1alpha1.VXLANClaim) error {
	if r.parentTreeName == "" {
		// a claim in the main tree
		if r.claimID != nil {
			if err := r.cacheCtx.tree.Update(id32.NewID(*r.claimID, 32), claim.GetClaimLabels()); err != nil {
				return err
			}
		} else {
			if err := r.cacheCtx.tree.Claim(id32.NewID(*claim.Spec.ID, 32), claim.GetClaimLabels()); err != nil {
				return err
			}
		}
	} else {
		fmt.Println("static id in range", r.parentTreeName)
		// a claim in a range
		k := store.ToKey(r.parentTreeName)
		table, err := r.cacheCtx.ranges.Get(ctx, k)
		if err != nil {
			return fmt.Errorf("selectAddress range does not have corresponding range table: err: %s", err.Error())
		}
		if r.claimID != nil {
			if err := table.Update(*claim.Spec.ID, claim.GetClaimLabels()); err != nil {
				return err
			}
		} else {
			if err := table.Claim(*claim.Spec.ID, claim.GetClaimLabels()); err != nil {
				return err
			}
		}

	}

	claim.Status.ID = claim.Spec.ID
	return nil
}

func (r *staticVXLANApplicator) Delete(ctx context.Context, claim *vxlanbev1alpha1.VXLANClaim) error {
	return r.delete(ctx, claim)
}
