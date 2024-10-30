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
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	"github.com/kuidio/kuid/apis/backend/ipam"
)

type staticAddressApplicator struct {
	name string
	applicator
	exists                 bool
	parentClaimSummaryType ipam.IPClaimSummaryType
	parentRangeName        string
	parentNetwork          bool
	parentLabels           map[string]string
}

func (r *staticAddressApplicator) Validate(ctx context.Context, claim *ipam.IPClaim) error {
	log := log.FromContext(ctx)
	log.Debug("validate")
	if claim.Spec.Address == nil {
		return fmt.Errorf("cannot claim a static ip address without an address")
	}
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

func (r *staticAddressApplicator) validateExists(_ context.Context, claim *ipam.IPClaim) (bool, error) {
	pi, err := iputil.New(*claim.Spec.Address)
	if err != nil {
		return false, err
	}

	// There is 2 scenario's:
	// an address with /32 or /128 prefixLength: 10.0.0.1/32 -> address prefix
	// an address with a dedicated prefixLength: 10.0.0.1/24 (only allowed for network)
	// check the /32 or /128 equivalent in the rib
	route, ok := r.cacheInstanceCtx.rib.Get(pi.GetIPAddressPrefix())
	if ok {
		// if the route exists validate the owner
		routeLabels := route.Labels()
		// a range is an exception as it can overlap with an address in the main rib
		if routeLabels[backend.KuidIPAMClaimSummaryTypeKey] == string(ipam.IPClaimSummaryType_Range) {
			r.parentClaimSummaryType = ipam.IPClaimSummaryType_Range
			r.parentRangeName = routeLabels[backend.KuidClaimNameKey]
			return true, nil
		} else {
			if err := claim.ValidateOwner(routeLabels); err != nil {
				return false, err
			}
			return true, nil
		}
	}
	// address does not exist
	return false, nil
}

func (r *staticAddressApplicator) validateParents(ctx context.Context, claim *ipam.IPClaim) error {
	pi, err := iputil.New(*claim.Spec.Address)
	if err != nil {
		return err
	}

	parentRoutes := r.cacheInstanceCtx.rib.Parents(pi.GetIPAddressPrefix())
	if len(parentRoutes) == 0 {
		return fmt.Errorf("a prefix range always needs a parent")
	}

	// the library returns all parent routes, but we need to most specific
	parentRoute := findMostSpecificParent(parentRoutes)
	if err := r.validateParent(ctx, parentRoute, claim); err != nil {
		return err
	}
	return nil
}

func (r *staticAddressApplicator) validateParent(_ context.Context, route table.Route, claim *ipam.IPClaim) error {
	pi, err := iputil.New(*claim.Spec.Address)
	if err != nil {
		return err
	}

	routeLabels := route.Labels()
	parentIPPrefixType := ipam.GetIPPrefixTypeFromString(routeLabels[backend.KuidIPAMIPPrefixTypeKey])
	parentClaimSummaryType := ipam.GetIPClaimSummaryTypeFromString(routeLabels[backend.KuidIPAMClaimSummaryTypeKey])
	parentClaimName := routeLabels[backend.KuidClaimNameKey]
	// update the context such that the applicator can use this information to apply the IP
	r.parentClaimSummaryType = parentClaimSummaryType
	r.parentRangeName = parentClaimName
	r.parentLabels = getUserDefinedLabels(routeLabels)

	if pi.IsAddressPrefix() {
		// 32 or /128 -> cannot be claimed in a network or aggregate
		if parentIPPrefixType != nil &&
			(*parentIPPrefixType == ipam.IPPrefixType_Network) {
			return fmt.Errorf("a /32 or /128 address is not possible with a parent of type %s", *parentIPPrefixType)
		}
		if parentClaimSummaryType == ipam.IPClaimSummaryType_Range {
			k := store.ToKey(parentClaimName)
			ipTable, err := r.applicator.cacheInstanceCtx.ranges.Get(k)
			if err != nil {
				return err
			}
			route, err := ipTable.Get(pi.GetIPAddress().String())
			if err == nil { // error means not found
				if err := claim.ValidateOwner(route.Labels()); err != nil {
					return fmt.Errorf("address is already allocated in range %s, err: %s", parentClaimName, err.Error())
				}
			}
			return nil
		}
	} else {
		// an address with a dedicated prefixLength is only possible for network prefix parents
		if parentIPPrefixType != nil && *parentIPPrefixType != ipam.IPPrefixType_Network {
			return fmt.Errorf("a prefix based address is not possible with a parent of type %s", *parentIPPrefixType)
		}
		if parentClaimSummaryType == ipam.IPClaimSummaryType_Range {
			return fmt.Errorf("a prefix based address is not possible for a %v", parentIPPrefixType)
		}
		r.parentNetwork = true
	}
	return nil
}

func (r *staticAddressApplicator) validateChildren(_ context.Context, claim *ipam.IPClaim) error {
	pi, err := iputil.New(*claim.Spec.Address)
	if err != nil {
		return err
	}
	childRoutes := r.cacheInstanceCtx.rib.Children(pi.GetIPAddressPrefix())
	if len(childRoutes) > 0 {
		return fmt.Errorf("an address based prefix %s cannot have children", *claim.Spec.Address)
	}

	return nil
}

func (r *staticAddressApplicator) Apply(ctx context.Context, claim *ipam.IPClaim) error {
	log := log.FromContext(ctx)
	log.Debug("apply")

	pi, err := iputil.New(*claim.Spec.Address)
	if err != nil {
		return err
	}
	if r.parentClaimSummaryType == ipam.IPClaimSummaryType_Range {
		if err := r.applyAddressInRange(ctx, claim, pi, r.parentRangeName, r.parentLabels); err != nil {
			return err
		}
	} else {
		if err := r.apply(ctx, claim, []*iputil.Prefix{pi}, false, r.parentLabels); err != nil {
			return err
		}
	}

	r.updateClaimAddressStatus(ctx, claim, pi, r.parentNetwork)
	return nil
}

func (r *staticAddressApplicator) Delete(ctx context.Context, claim *ipam.IPClaim) error {
	log := log.FromContext(ctx)
	log.Debug("delete")
	return r.delete(ctx, claim)
}
