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

package ipam

import (
	"context"
	"fmt"

	"github.com/hansthienpondt/nipam/pkg/table"
	"github.com/henderiw/iputil"
	"github.com/henderiw/logger/log"
	"github.com/kuidio/kuid/apis/backend/ipam"
	"go4.org/netipx"
)

type staticRangeApplicator struct {
	name string
	applicator
	exists bool
}

func (r *staticRangeApplicator) Validate(ctx context.Context, claim *ipam.IPClaim) error {
	log := log.FromContext(ctx)
	log.Debug("validate")

	exists, err := r.validateExists(ctx, claim)
	if err != nil {
		return err
	}
	if exists {
		// we can return since we trust the previous insertion in the tree
		r.exists = exists
		return nil
	}

	if err := r.validateParents(ctx, claim); err != nil {
		return err
	}

	if err := r.validateChildren(ctx, claim); err != nil {
		return err
	}

	return nil
}

func (r *staticRangeApplicator) validateExists(_ context.Context, claim *ipam.IPClaim) (bool, error) {
	ipRange, err := netipx.ParseIPRange(*claim.Spec.Range)
	if err != nil {
		return false, err
	}
	exists := true
	for _, prefix := range ipRange.Prefixes() {
		route, ok := r.cacheInstanceCtx.rib.Get(prefix)
		if ok {
			// entry exists; validate the owner to see if someone else owns this prefix
			labels := route.Labels()
			if err := claim.ValidateOwner(labels); err != nil {
				return exists, err
			}
		} else {
			exists = false
		}
	}
	return exists, nil
}

func (r *staticRangeApplicator) validateParents(ctx context.Context, claim *ipam.IPClaim) error {
	ipRange, err := netipx.ParseIPRange(*claim.Spec.Range)
	if err != nil {
		return err
	}

	var singlParentPrefix string
	for _, prefix := range ipRange.Prefixes() {
		parentRoutes := r.cacheInstanceCtx.rib.Parents(prefix)
		if len(parentRoutes) == 0 {
			return fmt.Errorf("a prefix range always needs a parent")
		}

		// the library returns all parent routes, but we need to most specific
		parentRoute := findMostSpecificParent(parentRoutes)
		if err := r.validateParent(ctx, parentRoute, claim); err != nil {
			return err
		}

		// check if the routes all belong to the same parent
		if singlParentPrefix != "" {
			if parentRoute.Prefix().String() != singlParentPrefix {
				return fmt.Errorf("a range has to fit into a single parent prefix got %s and %s", singlParentPrefix, parentRoute.Prefix().String())
			}
		} else {
			singlParentPrefix = parentRoute.Prefix().String()
		}
	}
	return nil
}

func (r *staticRangeApplicator) validateParent(_ context.Context, _ table.Route, _ *ipam.IPClaim) error {
	// a range can be allocated on any parent
	return nil
}

func (r *staticRangeApplicator) validateChildren(_ context.Context, claim *ipam.IPClaim) error {
	ipRange, err := netipx.ParseIPRange(*claim.Spec.Range)
	if err != nil {
		return err
	}
	for _, prefix := range ipRange.Prefixes() {
		childRoutes := r.cacheInstanceCtx.rib.Children(prefix)
		if len(childRoutes) > 0 {
			return fmt.Errorf("cannot have children for a range %s", prefix)
		}
	}
	return nil
}

func (r *staticRangeApplicator) Apply(ctx context.Context, claim *ipam.IPClaim) error {
	ipRange, err := netipx.ParseIPRange(*claim.Spec.Range)
	if err != nil {
		return err
	}

	// add each prefix in the routing table -> we convey them all together
	pis := make([]*iputil.Prefix, 0, len(ipRange.Prefixes()))
	for _, prefix := range ipRange.Prefixes() {
		pis = append(pis, iputil.NewPrefixInfo(prefix))
	}
	if err := r.apply(ctx, claim, pis, false, map[string]string{}); err != nil {
		return err
	}
	if err := r.applyRange(ctx, claim, ipRange); err != nil {
		return err
	}

	r.updateClaimRangeStatus(ctx, claim)
	return nil
}

func (r *staticRangeApplicator) Delete(ctx context.Context, claim *ipam.IPClaim) error {
	return r.delete(ctx, claim)
}
