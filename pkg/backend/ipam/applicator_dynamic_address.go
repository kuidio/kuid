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
	"github.com/kuidio/kuid/apis/backend/ipam"
)

type dynamicAddressApplicator struct {
	name string
	applicator
	parentClaimSummaryType ipam.IPClaimSummaryType
	parentRangeName        string
	parentNetwork          bool
	parentLabels           map[string]string
}

func (r *dynamicAddressApplicator) Validate(ctx context.Context, claim *ipam.IPClaim) error {
	return nil
}

func (r *dynamicAddressApplicator) Apply(ctx context.Context, claim *ipam.IPClaim) error {
	log := log.FromContext(ctx)
	log.Debug("apply")

	pi, err := r.validateExists(ctx, claim)
	if err != nil {
		return err
	}
	if pi == nil {
		// we need to claim an ip address
		pi, err = r.claimIP(ctx, claim)
		if err != nil {
			return err
		}
	}

	if r.parentClaimSummaryType == ipam.IPClaimSummaryType_Range {
		if err := r.applyAddressInRange(ctx, claim, pi, r.parentRangeName, r.parentLabels); err != nil {
			return err
		}
	} else {
		// the claimType is coming from the parent for addresses
		if err := r.apply(ctx, claim, []*iputil.Prefix{pi}, r.parentNetwork, r.parentLabels); err != nil {
			return err
		}
	}
	r.updateClaimAddressStatus(ctx, claim, pi, r.parentNetwork)
	return nil
}

func (r *dynamicAddressApplicator) validateExists(ctx context.Context, claim *ipam.IPClaim) (*iputil.Prefix, error) {
	existingRoutes, err := r.getRoutesByOwner(ctx, claim)
	if err != nil {
		return nil, err
	}

	for ribName, existingRoutes := range existingRoutes {
		if len(existingRoutes) > 1 {
			return nil, fmt.Errorf("cannot have multiple routes for an address entry")
		}
		// Now we have only 1 route
		if ribName == "" {
			for _, existingRoute := range existingRoutes {
				//
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
					if spi.Addr().String() == existingRoute.Prefix().Addr().String() {
						if spi.GetPrefixLength() != spi.GetAddressPrefixLength() {
							r.parentLabels = getUserDefinedLabels(findMostSpecificParent(existingRoute.Parents(r.cacheInstanceCtx.rib)).Labels())
							r.parentNetwork = true
						}
						return spi, nil
					}
				}
				// we delete the route if the claim status is empty or does not match
				// and reallocate
				if err := r.cacheInstanceCtx.rib.Delete(existingRoutes[0]); err != nil {
					return nil, err
				}
			}
		} else {
			k := store.ToKey(ribName)
			if len(existingRoutes) > 0 {
				if ipTable, err := r.cacheInstanceCtx.ranges.Get(k); err == nil {
					// the table exists
					for _, existingRoute := range existingRoutes {
						r.parentRangeName = ribName
						r.parentClaimSummaryType = ipam.IPClaimSummaryType_Range
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
	return nil, nil
}

func (r *dynamicAddressApplicator) claimIP(ctx context.Context, claim *ipam.IPClaim) (*iputil.Prefix, error) {
	// if not claimed, try to claim an address
	parentRoutes := r.getRoutesByLabel(ctx, claim)
	if len(parentRoutes) == 0 {
		return nil, fmt.Errorf("dynamic claim: no available routes based on the selector labels %v", claim.Spec.GetSelectorLabels())
	}
	return r.selectAddress(ctx, claim, parentRoutes)
}

func (r *dynamicAddressApplicator) selectAddress(_ context.Context, claim *ipam.IPClaim, parentRoutes table.Routes) (*iputil.Prefix, error) {
	routes := make([]string, 0, len(parentRoutes))
	for _, parentRoute := range parentRoutes {
		routes = append(routes, parentRoute.Prefix().String())
		routeLabels := parentRoute.Labels()
		parentIPPrefixType := ipam.GetIPPrefixTypeFromString(routeLabels[backend.KuidIPAMIPPrefixTypeKey])
		parentClaimSummaryType := ipam.GetIPClaimSummaryTypeFromString(routeLabels[backend.KuidIPAMClaimSummaryTypeKey])
		parentClaimName := routeLabels[backend.KuidClaimNameKey]
		// update the context such that the applicator can use this information to apply the IP
		r.parentClaimSummaryType = parentClaimSummaryType
		r.parentRangeName = parentClaimName
		r.parentLabels = getUserDefinedLabels(routeLabels)
		if parentIPPrefixType != nil && *parentIPPrefixType == ipam.IPPrefixType_Network {
			r.parentNetwork = true
		}

		pi := iputil.NewPrefixInfo(parentRoute.Prefix())

		switch parentClaimSummaryType {
		case ipam.IPClaimSummaryType_Range:
			// lookup range -> try to claim an ip from the range
			k := store.ToKey(parentClaimName)
			ipTable, err := r.cacheInstanceCtx.ranges.Get(k)
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

		case ipam.IPClaimSummaryType_Prefix:
			if parentIPPrefixType != nil && (*parentIPPrefixType == ipam.IPPrefixType_Network || *parentIPPrefixType == ipam.IPPrefixType_Pool) {
				parentpi := iputil.NewPrefixInfo(parentRoute.Prefix())
				if claim.Status.Address != nil {
					statuspi, err := iputil.New(*claim.Status.Address)
					if err != nil {
						return nil, err
					}
					// check if the route is free in the rib
					prefixLength := pi.GetAddressPrefixLength()
					if _, ok := r.cacheInstanceCtx.rib.Get(netip.PrefixFrom(statuspi.Addr(), prefixLength.Int())); !ok {
						return statuspi, nil
					}
				}

				// gather the prefixLength - use address based prefixLength /32 or /128 to validate the rib
				// for netowork allocations use the parent prefixLength
				prefixLength := pi.GetAddressPrefixLength()
				if isParentRouteSelectable(parentRoute, uint8(prefixLength)) {
					p := r.cacheInstanceCtx.rib.GetAvailablePrefixByBitLen(pi.GetIPPrefix(), uint8(prefixLength.Int()))
					if p.IsValid() {
						// success, parentClaimType was already checked for non nil
						if *parentIPPrefixType == ipam.IPPrefixType_Network {
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

func (r *dynamicAddressApplicator) Delete(ctx context.Context, claim *ipam.IPClaim) error {
	log := log.FromContext(ctx)
	log.Debug("delete")
	return r.delete(ctx, claim)
}
