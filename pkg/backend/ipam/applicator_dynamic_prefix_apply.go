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

	"github.com/henderiw/iputil"
	"github.com/henderiw/logger/log"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	"github.com/pkg/errors"
)

func (r *dynamicPrefixApplicator) Apply(ctx context.Context, claim *ipambev1alpha1.IPClaim) error {
	log := log.FromContext(ctx).With("name", claim.GetName())
	log.Info("dynamic prefix claim")

	// claim a prefix
	pi, err := r.claimPrefix(ctx, claim)
	if err != nil {
		return err
	}

	if err := r.apply(ctx, claim, []*iputil.Prefix{pi}, false); err != nil {
		return err
	}
	r.updateClaimPrefixStatus(ctx, claim, pi)
	return nil
}

// claimPrefix claims a prefix from the rib based on the claim (dynamic)
func (r *dynamicPrefixApplicator) claimPrefix(ctx context.Context, claim *ipambev1alpha1.IPClaim) (*iputil.Prefix, error) {
	log := log.FromContext(ctx)

	// first check if the resource is already claimed
	existingRoutes, err := r.getRoutesByOwner(ctx, claim)
	if err != nil {
		return nil, err
	}
	found := false
	var spi *iputil.Prefix
	for _, existingRoute := range existingRoutes {
		// validate if the existing prefix/address is in the routing
		// table -> if so we return -> apply takes care of the cleanup
		if claim.Status.Prefix != nil {
			spi, err = iputil.New(*claim.Status.Prefix)
			if err != nil {
				return nil, err
			}
			epi := iputil.NewPrefixInfo(existingRoute.Prefix())
			if spi.GetIPAddress() == epi.GetIPAddress() {
				found = true
				break
			}
		}
	}
	if found {
		return spi, nil
	}

	// if not claimed, try to claim the ip
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
		log.Info("refresh claimed prefix",
			"claimedPrefix", claim.Status.Prefix,
			"prefixlength", pi.GetPrefixLength())

		// check if the prefix is available
		p := r.cacheCtx.rib.GetAvailablePrefixByBitLen(pi.GetIPPrefix(), uint8(pi.GetPrefixLength()))
		if p.IsValid() {
			log.Info("refresh claimed prefix finished",
				"claimedPrefix", claim.Status.Prefix)
			// previously claimed prefix is available and reassigned
			return iputil.NewPrefixInfo(p), nil
		}
		log.Info("refresh claim prefix not available",
			"claimedPrefix", claim.Status.Prefix,
			"prefixlength", pi.GetPrefixLength())
	}

	// A prefix claim always need a prefix length
	prefixLength := iputil.PrefixLength(*claim.Spec.PrefixLength)
	for _, parentRoute := range parentRoutes {
		if isParentRouteSelectable(parentRoute, uint8(prefixLength)) {
			pi := iputil.NewPrefixInfo(parentRoute.Prefix())
			p := r.cacheCtx.rib.GetAvailablePrefixByBitLen(pi.GetIPPrefix(), uint8(prefixLength.Int()))
			if p.IsValid() {
				// success
				return iputil.NewPrefixInfo(p), nil
			}
		}
	}
	return nil, errors.New("no free prefix found")

}
