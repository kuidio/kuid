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

	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	"github.com/kform-dev/choreo/apis/condition"
	"github.com/kuidio/kuid/apis/backend"
	"k8s.io/utils/ptr"
)

type dynamicApplicator struct {
	name string
	applicator
}

func (r *dynamicApplicator) Validate(ctx context.Context, claim backend.ClaimObject) error {
	return nil
}

// validateExists will validate if the id specified in the claim exists already.
// If so it will reclaim it and update the parentTreeName (this is to ensure the claim reclaims it from the proper tree/table)
// if no entry exist it will keep r.claimID to nil
// On top the entries that were not claimed are cleaned up, such that we delete entries that are
// void. E.g. this takes care of the fact that a use changed the static ID as the reclaim failed
// so the remaining entry is cleaned up
func (r *dynamicApplicator) validateExists(ctx context.Context, claim backend.ClaimObject) (*uint64, string, error) {
	existingEntries, err := r.getEntriesByOwner(ctx, claim)
	if err != nil {
		return nil, "", err
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
		return nil, "", err
	}

	return claimID, claimTreeName, nil
}

// There are 2 scenario's.
// without a label selector: this claims from the main tree
// with a label selector: we expect a parent as this claims from a range
func (r *dynamicApplicator) getParentContext(ctx context.Context, claim backend.ClaimObject) (string, error) {
	// get entries by labelSelector if the label selector is defined
	if claim.GetSelector() == nil {
		// this is a dynamic claim for the main tree
		return "", nil
	}
	// this is a allocation for a range
	parentEntries := r.getEntriesByLabelSelector(ctx, claim)

	if len(parentEntries) == 0 {
		return "", fmt.Errorf("no parent found")
	}

	// validate if all parents are from the same range
	for _, parentEntry := range parentEntries {
		labels := parentEntry.Labels()
		parentClaimType := backend.GetClaimTypeFromString(labels[backend.KuidClaimTypeKey])
		if parentClaimType == backend.ClaimType_Range {
			return labels[backend.KuidClaimNameKey], nil
		} else {
			return "", fmt.Errorf("a parent can only be a range, got: %s", string(parentClaimType))
		}
	}
	return "", nil
}

func (r *dynamicApplicator) Apply(ctx context.Context, claim backend.ClaimObject) error {
	log := log.FromContext(ctx).With("name", claim.GetName())
	log.Debug("dynamic claim")

	claimID, parentTreeName, err := r.validateExists(ctx, claim)
	if err != nil {
		return err
	}
	if claimID == nil {
		parentTreeName, err = r.getParentContext(ctx, claim)
		if err != nil {
			return err
		}
	}

	if isReserved(parentTreeName, claim.GetIndex()) {
		return fmt.Errorf("cannot claim from a reserved range")
	}

	if parentTreeName == "" {
		// root tree apply
		if claimID != nil {
			if err := r.cacheInstanceCtx.tree.Update(claim.GetClaimID(r.cacheInstanceCtx.Type(), *claimID), claim.GetClaimLabels()); err != nil {
				return err
			}
			claim.SetStatusID(claimID)
			claim.SetConditions(condition.Ready())
			return nil
		}
		if claim.GetStatusID() != nil {
			// TODO check if free ?
			if err := r.cacheInstanceCtx.tree.ClaimID(claim.GetStatusClaimID(r.cacheInstanceCtx.Type()), claim.GetClaimLabels()); err != nil {
				return fmt.Errorf("reclaim status id claim failed, no claim ID found err: %s", err)
			}
			claim.SetStatusID(claim.GetStatusID())
			claim.SetConditions(condition.Ready())
			return nil

		}
		e, err := r.cacheInstanceCtx.tree.ClaimFree(claim.GetClaimLabels())
		if err != nil {
			return fmt.Errorf("claimed failed, no claim ID found err: %s", err)
		}
		claimID = ptr.To[uint64](e.ID().ID())
		claim.SetStatusID(claimID)
		claim.SetConditions(condition.Ready())
		return nil
	}
	// table - range entry apply
	k := store.ToKey(parentTreeName)
	table, err := r.cacheInstanceCtx.ranges.Get(k)
	if err != nil {
		return fmt.Errorf("selectAddress range does not have corresponding range table: err: %s", err.Error())
	}
	if claimID != nil {
		if err := table.Update(*claimID, claim.GetClaimLabels()); err != nil {
			return err
		}
		claim.SetStatusID(claimID)
		claim.SetConditions(condition.Ready())
		return nil
	}
	// try reclaim existing id
	if claim.GetStatusID() != nil {
		if table.IsFree(*claim.GetStatusID()) {
			if err := table.Claim(*claim.GetStatusID(), claim.GetClaimLabels()); err != nil {
				return fmt.Errorf("claimed failed, no claim ID found err: %s", err)
			}
			claim.SetStatusID(claim.GetStatusID())
			claim.SetConditions(condition.Ready())
			return nil
		}
	}
	e, err := table.ClaimFree(claim.GetClaimLabels())
	if err != nil {
		return fmt.Errorf("claimed failed, no claim ID found err: %s", err)
	}
	claimID = ptr.To[uint64](e.ID().ID())
	claim.SetStatusID(claimID)
	claim.SetConditions(condition.Ready())
	return nil

}

func (r *dynamicApplicator) Delete(ctx context.Context, claim backend.ClaimObject) error {
	// this is a generic delete by owner
	return r.delete(ctx, claim)
}
