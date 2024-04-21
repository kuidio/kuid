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
	fmt.Println("claimAddress", pi.Prefix.String(), r.parentClaimSummaryType, r.parentRangeName, r.parentNetwork)
	if r.parentClaimSummaryType == ipambev1alpha1.IPClaimSummaryType_Range {
		if err := r.applyAddressInRange(ctx, claim, pi, r.parentRangeName); err != nil {
			return err
		}
	} else {
		// the claimType is coming from the parent for addresses
		if err := r.apply(ctx, claim, []*iputil.Prefix{pi}, r.parentNetwork); err != nil {
			return err
		}
	}
	fmt.Println("claimAddress after apply", pi.Prefix.String())
	r.updateClaimAddressStatus(ctx, claim, pi, r.parentNetwork)
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

	for ribName, existingRoutes := range existingRoutes {
		if len(existingRoutes) > 1 {
			return nil, fmt.Errorf("cannot have multiple routes for an address entry")
		}
		if ribName == "" {
			for _, existingRoute := range existingRoutes {
				// Now we have only 1 route
				if claim.Status.Address != nil {
					spi, err := iputil.New(*claim.Status.Address)
					if err != nil {
						return nil, err
					}
					// since network based addresses return the parent prefixlength
					// we need to return the same address from the status iso the one
					// in the rib
					// e.g. 10.0.0.1/24 is in the rib 10.0.0.1/32 but the status reflects
					// 10.0.0.1/24
					fmt.Println("claim Address", spi.Addr().String(), existingRoute.Prefix().Addr().String())
					if spi.Addr().String() == existingRoute.Prefix().Addr().String() {
						if spi.GetPrefixLength() != spi.GetAddressPrefixLength() {
							r.parentNetwork = true
						}
						return spi, nil
					}
				}
				// we delete the route if the claim status is empty or does not match
				// and reallocate
				if err := r.cacheCtx.rib.Delete(existingRoutes[0]); err != nil {
					return nil, err
				}
			}
		} else {
			k := store.ToKey(ribName)
			if len(existingRoutes) > 0 {
				if ipTable, err := r.cacheCtx.ranges.Get(ctx, k); err == nil {
					// the table exists
					for _, existingRoute := range existingRoutes {
						r.parentRangeName = ribName
						r.parentClaimSummaryType = ipambev1alpha1.IPClaimSummaryType_Range
						if claim.Status.Address != nil {
							spi, err := iputil.New(*claim.Status.Address)
							if err != nil {
								return nil, err
							}
							if spi.Addr().String() == existingRoute.Prefix().Addr().String() {
								return spi, nil
							}
						}
						// we delete the route if the claim status is empty or does not match
						// and reallocate
						if err := ipTable.Release(existingRoute.Prefix().Addr().String()); err != nil {
							return nil, err
						}
					}
				}
			}
		}
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
		parentIPPrefixType := ipambev1alpha1.GetIPPrefixTypeFromString(routeLabels[backend.KuidIPAMIPPrefixTypeKey])
		parentClaimSummaryType := ipambev1alpha1.GetIPClaimSummaryTypeFromString(routeLabels[backend.KuidIPAMClaimSummaryTypeKey])
		parentClaimName := routeLabels[backend.KuidClaimNameKey]
		// update the context such that the applicator can use this information to apply the IP
		r.parentClaimSummaryType = parentClaimSummaryType
		r.parentRangeName = parentClaimName
		if parentIPPrefixType != nil && *parentIPPrefixType == ipambev1alpha1.IPPrefixType_Network {
			r.parentNetwork = true
		}

		pi := iputil.NewPrefixInfo(parentRoute.Prefix())

		switch parentClaimSummaryType {
		case ipambev1alpha1.IPClaimSummaryType_Range:
			// lookup range -> try to claim an ip from the range
			k := store.ToKey(parentClaimName)
			ipTable, err := r.cacheCtx.ranges.Get(ctx, k)
			if err != nil {
				return nil, fmt.Errorf("selectAddress range does not have corresponding range table: err: %s", err.Error())
			}
			if claim.Status.Address != nil {
				statuspi, err := iputil.New(*claim.Status.Address)
				if err != nil {
					return nil, err
				}

				route, err := ipTable.Get(statuspi.GetIPAddress().String())
				if err == nil { // error means found
					if err := claim.ValidateOwner(route.Labels()); err == nil {
						return statuspi, nil // route already exists
					}
				}
			}
			addr, err := ipTable.FindFree()
			if err != nil {
				return nil, err
			}
			return iputil.NewPrefixInfo(netip.PrefixFrom(addr, int(pi.GetAddressPrefixLength()))), nil

		case ipambev1alpha1.IPClaimSummaryType_Prefix:
			if parentIPPrefixType != nil && (*parentIPPrefixType == ipambev1alpha1.IPPrefixType_Network || *parentIPPrefixType == ipambev1alpha1.IPPrefixType_Pool) {
				parentpi := iputil.NewPrefixInfo(parentRoute.Prefix())
				if claim.Status.Address != nil {
					statuspi, err := iputil.New(*claim.Status.Address)
					if err != nil {
						return nil, err
					}
					fmt.Println("address status not empty", statuspi.Prefix.String())
					// check if the route is free in the rib
					prefixLength := pi.GetAddressPrefixLength()
					if _, ok := r.cacheCtx.rib.Get(netip.PrefixFrom(statuspi.Addr(), prefixLength.Int())); !ok {
						return statuspi, nil
					}
				}

				// gather the prefixLength - use address based prefixLength /32 or /128 to validate the rib
				// for netowork allocations use the parent prefixLength
				prefixLength := pi.GetAddressPrefixLength()
				fmt.Println("prefixLength", prefixLength)
				if isParentRouteSelectable(parentRoute, uint8(prefixLength)) {
					p := r.cacheCtx.rib.GetAvailablePrefixByBitLen(pi.GetIPPrefix(), uint8(prefixLength.Int()))
					fmt.Println("addr", p.Addr().String())
					if p.IsValid() {
						// success, parentClaimType was already checked for non nil
						if *parentIPPrefixType == ipambev1alpha1.IPPrefixType_Network {
							fmt.Println("parent prefixLength", parentpi.GetPrefixLength())
							return iputil.NewPrefixInfo(netip.PrefixFrom(p.Addr(), int(parentpi.GetPrefixLength()))), nil
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
