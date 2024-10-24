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

/*

import (
	"context"
	"fmt"

	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	"k8s.io/utils/ptr"
)

type dynamicApplicator struct {
	name string
	applicator
	claimID        *uint64
	parentTreeName string
}

func (r *dynamicApplicator) Validate(ctx context.Context, claim backend.ClaimObject) error {
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

func (r *dynamicApplicator) reclaimID(ctx context.Context, claim backend.ClaimObject) error {
	existingEntries, err := r.getEntriesByOwner(ctx, claim)
	if err != nil {
		return err
	}

	var claimID *uint64
	var claimTreeName string
	if claim.GetStatusID() != nil {
		claimID, claimTreeName = reclaimIDFromExisitingEntries(existingEntries, *claim.GetStatusID())
	} else {
		claimID, claimTreeName = claimIDFromExisitingEntries(existingEntries)
	}
	// remove the existing entries that don't match the claimed ID
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
func (r *dynamicApplicator) getParentContext(ctx context.Context, claim backend.ClaimObject) error {
	//log := log.FromContext(ctx).With("name", claim.GetName())
	// get entries by labelSelector if the label selector is defined
	if claim.GetSelector() == nil {
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
		parentClaimType := backend.GetClaimTypeFromString(labels[backend.KuidClaimTypeKey])
		if parentClaimType == backend.ClaimType_Range {
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

func (r *dynamicApplicator) Apply(ctx context.Context, claim backend.ClaimObject) error {
	log := log.FromContext(ctx).With("name", claim.GetName())
	log.Info("dynamic claim")
	if isReserved(r.parentTreeName, claim.GetIndex()) {
		return fmt.Errorf("cannot claim from a reserved range")
	}

	if r.parentTreeName == "" {
		if r.claimID != nil {
			if err := r.cacheCtx.tree.Update(claim.GetClaimID(r.cacheCtx.Type(), *r.claimID), claim.GetClaimLabels()); err != nil {
				return err
			}
		} else {
			e, err := r.cacheCtx.tree.ClaimFree(claim.GetClaimLabels())
			if err != nil {
				return err
			}
			r.claimID = ptr.To[uint64](e.ID().ID())
		}
	} else {
		k := store.ToKey(r.parentTreeName)
		table, err := r.cacheCtx.ranges.Get(ctx, k)
		if err != nil {
			return fmt.Errorf("selectAddress range does not have corresponding range table: err: %s", err.Error())
		}
		if r.claimID != nil {
			if err := table.Update(*r.claimID, claim.GetClaimLabels()); err != nil {
				return err
			}
		} else {
			e, err := table.ClaimFree(claim.GetClaimLabels())
			if err != nil {
				return err
			}
			r.claimID = ptr.To[uint64](e.ID().ID())
		}
	}
	//fmt.Println("apply", r.claimID)
	if r.claimID == nil {
		return fmt.Errorf("claimed failed, no claim ID found")
	}
	claim.SetStatusID(r.claimID)
	return nil
}

func (r *dynamicApplicator) Delete(ctx context.Context, claim backend.ClaimObject) error {
	// this is a generic delete by owner
	return r.delete(ctx, claim)
}
*/
