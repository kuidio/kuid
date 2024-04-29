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

package vlan

import (
	"context"
	"fmt"

	"github.com/henderiw/idxtable/pkg/tree/id32"
	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	vlanbev1alpha1 "github.com/kuidio/kuid/apis/backend/vlan/v1alpha1"
	"k8s.io/utils/ptr"
)

type dynamicVLANApplicator struct {
	name string
	applicator
	claimID        *uint32
	parentTreeName string
}

func (r *dynamicVLANApplicator) Validate(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error {
	if err := r.reclaimID(ctx, claim); err != nil {
		return err
	}
	// if the id is unknown we need to get the parent context
	// to determine if the claim is within the main tree or
	// within a range
	if r.claimID == nil {
		if err := r.getParentContext(ctx, claim); err != nil {
			return err
		}
	}
	return nil
}

func (r *dynamicVLANApplicator) reclaimID(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error {
	existingEntries, err := r.getEntriesByOwner(ctx, claim)
	if err != nil {
		return err
	}

	var claimID *uint32
	var claimTreeName string
	if claim.Status.ID != nil {
		claimID, claimTreeName = reclaimIDFromExisitingEntries(existingEntries, *claim.Status.ID)
	} else {
		claimID, claimTreeName = claimIDFromExisitingEntries(existingEntries)
	}
	// remove the existing entries that done match the claimed ID
	// should be none, but just in case
	if err := r.deleteNonClaimedEntries(ctx, existingEntries, claimID, claimTreeName); err != nil {
		return err
	}

	r.claimID = claimID
	if claimID != nil {
		r.parentTreeName = claimTreeName
	}
	return nil
}

// There are 2 scenario's.
// without a label selector: this claims from the main tree
// with a label selector: we expect a parent as this claims from a range
func (r *dynamicVLANApplicator) getParentContext(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error {
	//log := log.FromContext(ctx).With("name", claim.GetName())
	// get entries by labelSelector if the label selector is defined
	if claim.Spec.Selector == nil {
		// this is a dyanmic claim for the main tree
		return nil
	}
	// this is a allocation for a range
	parentEntries := r.getEntriesByLabelSelector(ctx, claim)

	if len(parentEntries) == 0 {
		return fmt.Errorf("no parent found")
	}

	// validate if all parents are from the same range
	for _, parentEntry := range parentEntries {
		labels := parentEntry.Labels()
		parentClaimType := vlanbev1alpha1.GetIPClaimTypeFromString(labels[backend.KuidVLANClaimTypeKey])
		if parentClaimType == vlanbev1alpha1.VLANClaimType_Range {
			if r.parentTreeName != "" && r.parentTreeName != labels[backend.KuidClaimNameKey] {
				return fmt.Errorf("a dynamic claim can only come from 1 parent range got %s and %s", r.parentTreeName, labels[backend.KuidClaimNameKey])
			}
			r.parentTreeName = labels[backend.KuidClaimNameKey]
		} else {
			return fmt.Errorf("a parent can only be a range, got: %s", string(parentClaimType))
		}
	}

	return nil
}

func (r *dynamicVLANApplicator) Apply(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error {
	log := log.FromContext(ctx).With("name", claim.GetName())
	log.Info("dynamic vlan claim")
	if isReserved(r.parentTreeName, claim.Spec.Index) {
		return fmt.Errorf("cannot claim from a reserved range")
	}
	if r.parentTreeName == "" {
		if r.claimID != nil {
			if err := r.cacheCtx.tree.Update(id32.NewID(*r.claimID, 32), claim.GetClaimLabels()); err != nil {
				return err
			}
		} else {
			e, err := r.cacheCtx.tree.ClaimFree(claim.GetClaimLabels())
			if err != nil {
				return err
			}
			r.claimID = ptr.To[uint32](uint32(e.ID().ID()))
		}
	} else {
		k := store.ToKey(r.parentTreeName)
		table, err := r.cacheCtx.ranges.Get(ctx, k)
		if err != nil {
			return fmt.Errorf("selectAddress range does not have corresponding range table: err: %s", err.Error())
		}
		if r.claimID != nil {
			if err := table.Update(uint16(*r.claimID), claim.GetClaimLabels()); err != nil {
				return err
			}
		} else {
			e, err := table.ClaimFree(claim.GetClaimLabels())
			if err != nil {
				return err
			}
			r.claimID = ptr.To[uint32](uint32(e.ID().ID()))
		}
	}

	if r.claimID == nil {
		return fmt.Errorf("claimed failed, no claim ID found")
	}
	claim.Status.ID = ptr.To[uint32](uint32(*r.claimID))
	return nil
}

func (r *dynamicVLANApplicator) Delete(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error {
	// this is a generic delete by owner
	return r.delete(ctx, claim)
}
