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
	parentClaimInfo ipambev1alpha1.IPClaimInfo
	parentRangeName string
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
		// if the route exists validate the owner
		routeLabels := route.Labels()
		// a range is an exception as it can overlap with an address
		if routeLabels[backend.KuidIPAMInfoKey] == string(ipambev1alpha1.IPClaimInfo_Range) {
			r.parentClaimInfo = ipambev1alpha1.IPClaimInfo_Range
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
	if err := r.validateExistingParent(ctx, claim, pi, parentRoute); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (r *staticAddressApplicator) validateExistingParent(ctx context.Context, _ *ipambev1alpha1.IPClaim, pi *iputil.Prefix, route table.Route) error {
	routeLabels := route.Labels()
	parentClaimInfo := routeLabels[backend.KuidIPAMInfoKey]
	parentClaimType := routeLabels[backend.KuidIPAMTypeKey]
	parentClaimName := routeLabels[backend.KuidClaimNameKey]

	if pi.IsAddressPrefix() {
		// 32 or /128 -> cannot be claimed in a network or aggregate
		if parentClaimType == string(ipambev1alpha1.IPClaimType_Network) ||
			parentClaimType == string(ipambev1alpha1.IPClaimType_Aggregate) {
			return fmt.Errorf("a /32 or /128 address is not possible with a parent of type %s", parentClaimType)
		}
		if parentClaimInfo == string(ipambev1alpha1.IPClaimInfo_Range) {
			k := store.ToKey(parentClaimName)
			ipTable, err := r.applicator.cacheCtx.ranges.Get(ctx, k)
			if err != nil {
				return err
			}
			if !ipTable.IsFree(pi.GetIPAddress().String()) {
				return fmt.Errorf("address is already allocated in range %s", parentClaimName)
			}
			r.parentClaimInfo = ipambev1alpha1.IPClaimInfo_Range
			r.parentRangeName = parentClaimName
		}
	} else {
		// an address with a dedicated prefixLength is only possible for network prefix parents
		if parentClaimType != string(ipambev1alpha1.IPClaimType_Network) {
			return fmt.Errorf("a prefix based address is not possible with a parent of type %s", parentClaimType)
		}
		if parentClaimInfo == string(ipambev1alpha1.IPClaimInfo_Range) {
			return fmt.Errorf("a prefix based address is not possible for a %s", parentClaimInfo)
		}
	}
	return nil
}
