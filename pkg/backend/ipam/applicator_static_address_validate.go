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
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
)

type staticAddressApplicator struct {
	name string
	applicator
	parentClaimSummaryType ipambev1alpha1.IPClaimSummaryType
	parentRangeName        string
	parentNetwork          bool
	parentLabels           map[string]string
}

func (r *staticAddressApplicator) Validate(ctx context.Context, claim *ipambev1alpha1.IPClaim) error {
	log := log.FromContext(ctx)
	pi, err := iputil.New(*claim.Spec.Address)
	if err != nil {
		return err
	}
	// get dryrun rib
	dryrunRib := r.cacheCtx.rib.Clone()

	// There is 2 scenario's:
	// an address with /32 or /128 prefixLength: 10.0.0.1/32 -> address prefix
	// an address with a dedicated prefixLength: 10.0.0.1/24 (onlky allowed for network)
	// check the /32 or /128 equivalent in the rib
	route, ok := dryrunRib.Get(pi.GetIPAddressPrefix())
	if ok {
		fmt.Println("static address route exists", *claim.Spec.Address, route.Prefix().String())
		// if the route exists validate the owner
		routeLabels := route.Labels()
		// a range is an exception as it can overlap with an address
		if routeLabels[backend.KuidIPAMClaimSummaryTypeKey] == string(ipambev1alpha1.IPClaimSummaryType_Range) {
			r.parentClaimSummaryType = ipambev1alpha1.IPClaimSummaryType_Range
			r.parentRangeName = routeLabels[backend.KuidClaimNameKey]
			return nil
		} else {
			if err := claim.ValidateOwner(routeLabels); err != nil {
				return err
			}
			return nil
		}
	}
	route = table.NewRoute(
		pi.GetIPAddressPrefix(),
		map[string]string{},
		map[string]any{},
	)
	if err := dryrunRib.Add(route); err != nil {
		log.Error("cannot add route", "route", route, "error", err.Error())
		return err
	}
	// get the route again and check for children
	route, ok = dryrunRib.Get(pi.GetIPAddressPrefix())
	if !ok {
		err := fmt.Errorf("cannot get route %s which just got addded", pi.GetIPSubnet())
		log.Error(err.Error())
		return err
	}
	// check for children
	routes := route.Children(dryrunRib)
	if len(routes) > 0 {
		err := fmt.Errorf("cannot have children for an address %s", pi.GetIPPrefix())
		log.Error(err.Error())
		return err
	}
	// get parents
	routes = route.Parents(dryrunRib)
	if len(routes) == 0 {
		// no parents exist
		if err := validateNoParent(claim); err != nil {
			return err
		}
	}
	// parents exist
	parentRoute := findParent(routes)
	fmt.Println("static address parent route", parentRoute.Prefix(), parentRoute.Labels())
	if err := r.validateExistingParent(ctx, claim, pi, parentRoute); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (r *staticAddressApplicator) validateExistingParent(ctx context.Context, claim *ipambev1alpha1.IPClaim, pi *iputil.Prefix, route table.Route) error {
	routeLabels := route.Labels()
	parentIPPrefixType := ipambev1alpha1.GetIPPrefixTypeFromString(routeLabels[backend.KuidIPAMIPPrefixTypeKey])
	parentClaimSummaryType := ipambev1alpha1.GetIPClaimSummaryTypeFromString(routeLabels[backend.KuidIPAMClaimSummaryTypeKey])
	parentClaimName := routeLabels[backend.KuidClaimNameKey]
	// update the context such that the applicator can use this information to apply the IP
	r.parentClaimSummaryType = parentClaimSummaryType
	r.parentRangeName = parentClaimName
	r.parentLabels = getUserDefinedLabels(routeLabels)

	if pi.IsAddressPrefix() {
		// 32 or /128 -> cannot be claimed in a network or aggregate
		if parentIPPrefixType != nil &&
			(*parentIPPrefixType == ipambev1alpha1.IPPrefixType_Network || *parentIPPrefixType == ipambev1alpha1.IPPrefixType_Aggregate) {
			return fmt.Errorf("a /32 or /128 address is not possible with a parent of type %s", *parentIPPrefixType)
		}
		if parentClaimSummaryType == ipambev1alpha1.IPClaimSummaryType_Range {
			k := store.ToKey(parentClaimName)
			ipTable, err := r.applicator.cacheCtx.ranges.Get(ctx, k)
			if err != nil {
				return err
			}
			route, err := ipTable.Get(pi.GetIPAddress().String())
			if err == nil { // error means not found
				fmt.Println("range address labels", route.Labels())
				if err := claim.ValidateOwner(route.Labels()); err != nil {
					fmt.Println("owner error", err.Error())
					return fmt.Errorf("address is already allocated in range %s, err: %s", parentClaimName, err.Error())
				}
			}
			return nil
		}
	} else {
		// an address with a dedicated prefixLength is only possible for network prefix parents
		if parentIPPrefixType != nil && *parentIPPrefixType != ipambev1alpha1.IPPrefixType_Network {
			return fmt.Errorf("a prefix based address is not possible with a parent of type %s", *parentIPPrefixType)
		}
		if parentClaimSummaryType == ipambev1alpha1.IPClaimSummaryType_Range {
			return fmt.Errorf("a prefix based address is not possible for a %v", parentIPPrefixType)
		}
		r.parentNetwork = true
	}
	return nil
}
