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
	"errors"
	"fmt"

	"github.com/henderiw/iputil"
	"github.com/henderiw/logger/log"
	"github.com/kuidio/kuid/apis/backend/ipam"
)

type dynamicPrefixApplicator struct {
	name string
	applicator
}

func (r *dynamicPrefixApplicator) Validate(ctx context.Context, claim *ipam.IPClaim) error {
	return nil
}

func (r *dynamicPrefixApplicator) Apply(ctx context.Context, claim *ipam.IPClaim) error {
	log := log.FromContext(ctx)
	log.Debug("apply")

	pi, err := r.validateExists(ctx, claim)
	if err != nil {
		return err
	}
	if pi == nil {
		// we need to claim an ip address
		pi, err = r.claimPrefix(ctx, claim)
		if err != nil {
			return err
		}
	}

	// the claimType is coming from the parent for addresses
	if err := r.apply(ctx, claim, []*iputil.Prefix{pi}, false, map[string]string{}); err != nil {
		return err
	}

	r.updateClaimPrefixStatus(ctx, claim, pi)
	return nil
}

func (r *dynamicPrefixApplicator) validateExists(ctx context.Context, claim *ipam.IPClaim) (*iputil.Prefix, error) {
	existingRoutes, err := r.getRoutesByOwner(ctx, claim)
	if err != nil {
		return nil, err
	}
	for _, existingRoute := range existingRoutes[""] {
		// validate if the existing prefix/address is in the routing
		// table -> if so we return -> apply takes care of the cleanup
		if claim.Status.Prefix != nil {
			spi, err := iputil.New(*claim.Status.Prefix)
			if err != nil {
				return nil, err
			}
			epi := iputil.NewPrefixInfo(existingRoute.Prefix())
			if spi.GetIPAddress() == epi.GetIPAddress() {
				return spi, nil
			}
		}
	}
	return nil, nil
}

func (r *dynamicPrefixApplicator) claimPrefix(ctx context.Context, claim *ipam.IPClaim) (*iputil.Prefix, error) {
	log := log.FromContext(ctx)
	// if not claimed, try to claim an address
	parentRoutes := r.getRoutesByLabel(ctx, claim)
	if len(parentRoutes) == 0 {
		return nil, fmt.Errorf("dynamic claim: no available routes based on the selector labels %v", claim.Spec.GetSelectorLabels())
	}
	// try to reclaim the prefix if the prefix was already claimed
	if claim.Status.Prefix != nil {
		pi, err := iputil.New(*claim.Status.Prefix)
		if err != nil {
			return nil, err
		}
		log.Debug("refresh claimed prefix",
			"claimedPrefix", claim.Status.Prefix,
			"prefixlength", pi.GetPrefixLength())

		// check if the prefix is available
		p := r.cacheInstanceCtx.rib.GetAvailablePrefixByBitLen(pi.GetIPPrefix(), uint8(pi.GetPrefixLength()))
		if p.IsValid() {
			log.Debug("refresh claimed prefix finished",
				"claimedPrefix", claim.Status.Prefix)
			// previously claimed prefix is available and reassigned
			return iputil.NewPrefixInfo(p), nil
		}
		log.Debug("refresh claim prefix not available",
			"claimedPrefix", claim.Status.Prefix,
			"prefixlength", pi.GetPrefixLength())
	}

	// A prefix claim always need a prefix length
	prefixLength := iputil.PrefixLength(*claim.Spec.PrefixLength)
	for _, parentRoute := range parentRoutes {
		if isParentRouteSelectable(parentRoute, uint8(prefixLength)) {
			pi := iputil.NewPrefixInfo(parentRoute.Prefix())
			p := r.cacheInstanceCtx.rib.GetAvailablePrefixByBitLen(pi.GetIPPrefix(), uint8(prefixLength.Int()))
			if p.IsValid() {
				// success
				return iputil.NewPrefixInfo(p), nil
			}
		}
	}
	return nil, errors.New("no free prefix found")
}

func (r *dynamicPrefixApplicator) Delete(ctx context.Context, claim *ipam.IPClaim) error {
	log := log.FromContext(ctx)
	log.Debug("delete")
	return r.delete(ctx, claim)
}
