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
	"net/netip"

	"github.com/hansthienpondt/nipam/pkg/table"
	"github.com/henderiw/iputil"
	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
)

func (r *dynamicAddressApplicator) Apply(ctx context.Context, claim *ipambev1alpha1.IPClaim) error {
	log := log.FromContext(ctx).With("name", claim.GetName())
	log.Info("dynamic address claim")

	// claim a address
	pi, err := r.claimAddress(ctx, claim)
	if err != nil {
		return err
	}
	if r.parentClaimInfo == ipambev1alpha1.IPClaimInfo_Range {
		if err := r.applyAddressInRange(ctx, claim, pi, r.parentRangeName); err != nil {
			return err
		}
	} else {
		// the claimType is coming from the parent for addresses
		if err := r.apply(ctx, claim, []*iputil.Prefix{pi}, r.parentNetwork); err != nil {
			return err
		}
	}

	r.updateClaimAddressStatus(ctx, claim, pi)
	return nil
}

// claimPrefix claims a prefix from the rib based on the claim (dynamic)
func (r *dynamicAddressApplicator) claimAddress(ctx context.Context, claim *ipambev1alpha1.IPClaim) (*iputil.Prefix, error) {
	//log := log.FromContext(ctx)

	// first check if the resource is already claimed
	existingRoutes, err := r.getRoutesByOwner(ctx, claim)
	if err != nil {
		return nil, err
	}

	if len(existingRoutes) > 1 {
		return nil, fmt.Errorf("cannot have multiple routes for an address entry")
	}
	if len(existingRoutes) == 1 {
		// This should be equal to the status
		return iputil.NewPrefixInfo(existingRoutes[0].Prefix()), nil
	}

	// if not claimed, try to claim an address
	parentRoutes := r.getRoutesByLabel(ctx, claim)
	if len(parentRoutes) == 0 {
		return nil, fmt.Errorf("dynamic claim: no available routes based on the selector labels %v", claim.Spec.GetSelectorLabels())
	}

	return r.selectAddress(ctx, claim, parentRoutes)
}

/*
walk over the routes
*/
func (r *dynamicAddressApplicator) selectAddress(ctx context.Context, claim *ipambev1alpha1.IPClaim, parentRoutes table.Routes) (*iputil.Prefix, error) {
	routes := make([]string, 0, len(parentRoutes))
	for _, parentRoute := range parentRoutes {
		routes = append(routes, parentRoute.Prefix().String())
		routeLabels := parentRoute.Labels()
		parentClaimType := ipambev1alpha1.GetIPClaimTypeFromString(routeLabels[backend.KuidIPAMTypeKey])
		parentClaimInfo := ipambev1alpha1.GetIPClaimInfoFromString(routeLabels[backend.KuidIPAMInfoKey])
		parentClaimName := routeLabels[backend.KuidClaimNameKey]
		// update the context such that the applicator can use this information to apply the IP
		r.parentClaimInfo = parentClaimInfo
		r.parentRangeName = parentClaimName

		pi := iputil.NewPrefixInfo(parentRoute.Prefix())

		switch parentClaimInfo {
		case ipambev1alpha1.IPClaimInfo_Range:
			// lookup range -> try to claim an ip from the range
			k := store.ToKey(parentClaimName)
			ipTable, err := r.cacheCtx.ranges.Get(ctx, k)
			if err != nil {
				return nil, fmt.Errorf("selectAddress range does not have corresponding range table: err: %s", err.Error())
			}
			if claim.Status.Address != nil {
				pi, err := iputil.New(*claim.Status.Address)
				if err != nil {
					return nil, err
				}

				if ipTable.IsFree(pi.GetIPAddress().String()) {
					return pi, nil
				}
			}
			addr, err := ipTable.FindFree()
			if err != nil {
				return nil, err
			}
			return iputil.NewPrefixInfo(netip.PrefixFrom(addr, int(pi.GetAddressPrefixLength()))), nil

		case ipambev1alpha1.IPClaimInfo_Prefix:
			if parentClaimType != nil && (*parentClaimType == ipambev1alpha1.IPClaimType_Network || *parentClaimType == ipambev1alpha1.IPClaimType_Pool) {
				if claim.Status.Address != nil {
					pi, err := iputil.New(*claim.Status.Address)
					if err != nil {
						return nil, err
					}
					// check if the prefix is available
					p := r.cacheCtx.rib.GetAvailablePrefixByBitLen(pi.GetIPPrefix(), uint8(pi.GetPrefixLength()))
					if p.IsValid() {
						// previously claimed prefix is available and reassigned
						return iputil.NewPrefixInfo(p), nil
					}
				}

				// gather the prefixLength - use address based prefixLength /32 or /128
				// for netowork allocations use the parent prefixLength
				prefixLength := pi.GetAddressPrefixLength()
				fmt.Println("prefixLength", prefixLength)
				if isParentRouteSelectable(parentRoute, uint8(prefixLength)) {
					pi := iputil.NewPrefixInfo(parentRoute.Prefix())
					p := r.cacheCtx.rib.GetAvailablePrefixByBitLen(pi.GetIPPrefix(), uint8(prefixLength.Int()))
					fmt.Println("addr", p.Addr().String())
					if p.IsValid() {
						// success, parentClaimType was already checked for non nil
						if *parentClaimType == ipambev1alpha1.IPClaimType_Network {
							r.parentNetwork = true
							return iputil.NewPrefixInfo(netip.PrefixFrom(p.Addr(), int(pi.GetPrefixLength()))), nil
						} else {
							return iputil.NewPrefixInfo(p), nil
						}
					}
				}
			}
		default:
		}
	}
	return nil, fmt.Errorf("no free addresses found in routes: %v", routes)
}
