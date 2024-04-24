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

	vlanbev1alpha1 "github.com/kuidio/kuid/apis/backend/vlan/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/utils/ptr"
)

type Applicator interface {
	Validate(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error
	Apply(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error
	Delete(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error
}

type applicator struct {
	cacheCtx *CacheContext
}

type dynamicVLANApplicator struct {
	name string
	applicator
}

func (r *dynamicVLANApplicator) Validate(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error {
	return nil
}

func (r *dynamicVLANApplicator) Apply(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error {
	entries, err := r.getEntriesByOwner(ctx, claim)
	if err != nil {
		return err
	}

	claimedID := getEntry(entries, claim)
	for id := range entries { // release the unallocated items
		if id != claimedID {
			if err := r.cacheCtx.table.Release(id); err != nil {
				return err
			}
		}
	}
	if claimedID != -1 {
		// claim id can be reused -> reapply claim to update labels
		if err := r.cacheCtx.table.Update(claimedID, claim.GetClaimLabels()); err != nil {
			return err
		}
		claim.Status.ID = ptr.To[uint32](uint32(claimedID))
		return nil
	}
	// dynamic claim
	id, err := r.cacheCtx.table.ClaimDynamic(claim.GetClaimLabels())
	if err != nil {
		return err
	}
	claim.Status.ID = ptr.To[uint32](uint32(id))
	return nil
}

func (r *dynamicVLANApplicator) Delete(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error {
	return r.delete(ctx, claim)
}

type staticVLANApplicator struct {
	name string
	applicator
}

func (r *staticVLANApplicator) Validate(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error {
	return nil
}

func (r *staticVLANApplicator) Apply(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error {
	entries, err := r.getEntriesByOwner(ctx, claim)
	if err != nil {
		return err
	}

	claimedID := getStaticEntry(entries, claim)
	for id := range entries { // release the unallocated items
		if id != claimedID {
			if err := r.cacheCtx.table.Release(id); err != nil {
				return err
			}
		}
	}
	if claimedID != -1 {
		// claim id can be reused -> reapply claim to update labels
		if err := r.cacheCtx.table.Update(claimedID, claim.GetClaimLabels()); err != nil {
			return err
		}
		claim.Status.ID = ptr.To[uint32](uint32(claimedID))
		return nil
	}
	// static claim
	if err := r.cacheCtx.table.Claim(int64(*claim.Spec.ID), claim.GetClaimLabels()); err != nil {
		return err
	}
	claim.Status.ID = claim.Spec.ID
	return nil
}

func (r *staticVLANApplicator) Delete(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error {
	return r.delete(ctx, claim)
}

type rangeVLANApplicator struct {
	name string
	applicator
}

func (r *rangeVLANApplicator) Validate(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error {
	return nil
}

func (r *rangeVLANApplicator) Apply(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error {
	entries, err := r.getEntriesByOwner(ctx, claim)
	if err != nil {
		return err
	}
	start, end := claim.GetVLANRange()
	if len(entries) == (end - start + 1) {
		_, startok := entries[int64(start)]
		_, endok := entries[int64(start)]
		if startok && endok {
			for id := range entries {
				if err := r.cacheCtx.table.Update(id, claim.GetClaimLabels()); err != nil {
					return err
				}
			}
			claim.Status.Range = claim.Spec.Range
			return nil
		}
	}
	// there was a change, delete the entries and reclaim if possible
	for id := range entries {
		if err := r.cacheCtx.table.Release(id); err != nil {
			return err
		}
	}
	size := end - start + 1
	if err := r.cacheCtx.table.ClaimRange(int64(start), int64(size), claim.GetClaimLabels()); err != nil {
		return err
	}
	claim.Status.Range = claim.Spec.Range
	return nil
}

func (r *rangeVLANApplicator) Delete(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error {
	return r.delete(ctx, claim)
}

type sizeVLANApplicator struct {
	name string
	applicator
}

func (r *sizeVLANApplicator) Validate(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error {
	return nil
}

func (r *sizeVLANApplicator) Apply(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error {
	entries, err := r.getEntriesByOwner(ctx, claim)
	if err != nil {
		return err
	}
	size := *claim.Spec.VLANSize

	if len(entries) == int(size) {
		for id := range entries {
			if err := r.cacheCtx.table.Update(id, claim.GetClaimLabels()); err != nil {
				return err
			}
		}
		// TODO how to reflect status
		return nil
	}
	// there was a change, delete the entries and reclaim if possible
	for id := range entries {
		if err := r.cacheCtx.table.Release(id); err != nil {
			return err
		}
	}
	if err := r.cacheCtx.table.ClaimSize(int64(size), claim.GetClaimLabels()); err != nil {
		return err
	}
	// TODO how to reflect status
	return nil

}

func (r *sizeVLANApplicator) Delete(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error {
	return r.delete(ctx, claim)
}

func (r *applicator) getEntriesByOwner(_ context.Context, claim *vlanbev1alpha1.VLANClaim) (map[int64]labels.Set, error) {
	ownerSelector, err := claim.GetOwnerSelector()
	if err != nil {
		return nil, err
	}
	entries := r.cacheCtx.table.GetByLabel(ownerSelector)
	if len(entries) != 0 {
		return entries, nil
	}
	return map[int64]labels.Set{}, nil
}

func getEntry(entries map[int64]labels.Set, claim *vlanbev1alpha1.VLANClaim) int64 {
	if claim.Status.ID != nil {
		if _, ok := entries[int64(*claim.Status.ID)]; ok {
			return int64(*claim.Status.ID)
		}
	}
	for id := range entries {
		return id
	}
	return -1
}

func getStaticEntry(entries map[int64]labels.Set, claim *vlanbev1alpha1.VLANClaim) int64 {
	if claim.Spec.ID != nil {
		if _, ok := entries[int64(*claim.Spec.ID)]; ok {
			return int64(*claim.Spec.ID)
		}
	}
	return -1
}

func (r *applicator) delete(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error {
	entries, err := r.getEntriesByOwner(ctx, claim)
	if err != nil {
		return err
	}
	for id := range entries {
		if err := r.cacheCtx.table.Release(id); err != nil {
			return err
		}
	}
	return nil
}
