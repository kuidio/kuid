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
)

type staticApplicator struct {
	name string
	applicator
}

func (r *staticApplicator) Validate(ctx context.Context, claim backend.ClaimObject) error {
	//log := log.FromContext(ctx)
	if claim.GetStaticID() == nil {
		return fmt.Errorf("cannot claim a static id without an id")
	}

	return nil
}

func (r *staticApplicator) Apply(ctx context.Context, claim backend.ClaimObject) error {
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
		return fmt.Errorf("cannot claim an id from a reserved range")
	}
	if parentTreeName == "" {
		// a claim in the main tree
		if claimID != nil {
			/*
				//r.cacheInstanceCtx.tree.PrintNodes()
				//r.cacheInstanceCtx.tree.PrintValues()

				//fmt.Println("staticApplicator update", "type", r.cacheInstanceCtx.Type(), "id", claim.GetClaimID(r.cacheInstanceCtx.Type(), *claimID))
			*/
			if err := r.cacheInstanceCtx.tree.Update(claim.GetClaimID(r.cacheInstanceCtx.Type(), *claimID), claim.GetClaimLabels()); err != nil {
				return err
			}
		} else {
			/*
				r.cacheInstanceCtx.tree.PrintNodes()
				r.cacheInstanceCtx.tree.PrintValues()

				fmt.Println("staticApplicator claim", "type", r.cacheInstanceCtx.Type(), claim.GetStaticTreeID(r.cacheInstanceCtx.Type()), claim.GetClaimLabels())

				statusID := "nil"
				if claim.GetStatusID() != nil {
					statusID = fmt.Sprintf("%d", *claim.GetStatusID())
				}
				staticID := "nil"
				if claim.GetStaticID() != nil {
					staticID = fmt.Sprintf("%d", *claim.GetStaticID())
				}

				fmt.Println("staticApplicator claim", "status", statusID, "id", staticID)
			*/
			if err := r.cacheInstanceCtx.tree.ClaimID(claim.GetStaticTreeID(r.cacheInstanceCtx.Type()), claim.GetClaimLabels()); err != nil {
				return err
			}
		}
	} else {
		// a claim in a range
		k := store.ToKey(parentTreeName)
		table, err := r.cacheInstanceCtx.ranges.Get(k)
		if err != nil {
			return fmt.Errorf("selectAddress range does not have corresponding range table: err: %s", err.Error())
		}
		if claimID != nil {
			if err := table.Update(*claim.GetStaticID(), claim.GetClaimLabels()); err != nil {
				return err
			}
		} else {
			if err := table.Claim(*claim.GetStaticID(), claim.GetClaimLabels()); err != nil {
				return err
			}
		}
	}
	claim.SetStatusID(claim.GetStaticID())
	claim.SetConditions(condition.Ready())
	return nil
}

// validateExists will validate if the id specified in the claim exists already.
// If so it will reclaim it and update the parentTreeName (this is to ensure the claim reclaims it from the proper tree/table)
// if no entry exist it will keep r.claimID to nil
// On top the entries that were not claimed are cleaned up, such that we delete entries that are
// void. E.g. this takes care of the fact that a use changed the static ID as the reclaim failed
// so the remaining entry is cleaned up
func (r *staticApplicator) validateExists(ctx context.Context, claim backend.ClaimObject) (*uint64, string, error) {
	existingEntries, err := r.getEntriesByOwner(ctx, claim)
	if err != nil {
		return nil, "", err
	}

	claimID, claimTreeName := reclaimIDFromExisitingEntries(existingEntries, *claim.GetStaticID())
	// remove the existing entries that don't match the claimed ID
	// should be none, but just in case
	if err := r.deleteNonClaimedEntries(ctx, existingEntries, claimID, claimTreeName); err != nil {
		return nil, "", err
	}
	return claimID, claimTreeName, nil
}

func (r *staticApplicator) getParentContext(ctx context.Context, claim backend.ClaimObject) (string, error) {
	log := log.FromContext(ctx).With("name", claim.GetName())

	entry, err := r.cacheInstanceCtx.tree.Get(claim.GetStaticTreeID(r.cacheInstanceCtx.Type()))
	if err == nil {
		// entry exists
		labels := entry.Labels()
		// a range can overlap so we return the entry as a parent if the entry match and it is a range
		claimType := backend.GetClaimTypeFromString(labels[backend.KuidClaimTypeKey])
		if claimType == backend.ClaimType_Range {
			return labels[backend.KuidClaimNameKey], nil
		} else {
			// This should always result in a different owner
			// since we checked the claimed entries before
			if err := claim.ValidateOwner(labels); err != nil {
				return "", err
			}
			return "", nil
		}
	}
	parentEntries := r.cacheInstanceCtx.tree.Parents(claim.GetStaticTreeID(r.cacheInstanceCtx.Type()))
	if len(parentEntries) > 1 {
		log.Error("got multiple parent entries", "entries", parentEntries)
		return "", fmt.Errorf("multiple parent entries %v", parentEntries)
	}
	for _, parentEntry := range parentEntries {
		labels := parentEntry.Labels()
		parentClaimType := backend.GetClaimTypeFromString(labels[backend.KuidClaimTypeKey])
		if parentClaimType == backend.ClaimType_Range {
			return labels[backend.KuidClaimNameKey], nil
		} else {
			log.Error("got parent which is not a range", "entry", parentEntry)
			return "", fmt.Errorf("got parent which is not a range %s", labels[backend.KuidClaimTypeKey])
		}
	}
	return "", nil
}

func (r *staticApplicator) Delete(ctx context.Context, claim backend.ClaimObject) error {
	return r.delete(ctx, claim)
}
