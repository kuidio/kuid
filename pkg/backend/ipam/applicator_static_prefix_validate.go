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
	"github.com/kuidio/kuid/apis/backend"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
)

type staticPrefixApplicator struct {
	name string
	applicator
}

func (r *staticPrefixApplicator) Validate(ctx context.Context, claim *ipambev1alpha1.IPClaim) error {
	log := log.FromContext(ctx)
	pi, err := iputil.New(*claim.Spec.Prefix)
	if err != nil {
		return err
	}
	// get dryrun rib
	dryrunRib := r.cacheCtx.rib.Clone()

	// There is 2 scenario's:
	// a regular prefix w/o address: 10.0.0.0/24
	// an address based prefix: 10.0.0.1/24 -> only allowed for claimType network
	// check if the net prefix/subnet exists
	route, ok := dryrunRib.Get(pi.GetIPSubnet())
	if ok {
		// route exists -> the owner need to match
		routeLabels := route.Labels()
		if err := claim.ValidateOwner(routeLabels); err != nil {
			return err
		}
		// for an address based prefix (which is only allowed) for networks
		// we need to also validate if the address part does not exist
		if pi.GetIPSubnet().String() != pi.GetIPPrefix().String() {
			// check if the existing net prefix/subnet is of type network
			if claim.GetIPPrefixType() != ipambev1alpha1.IPPrefixType_Network {
				return fmt.Errorf("a static address based prefix (net <> address) is only allowed for claimType network")
			}
			route, ok = dryrunRib.Get(pi.GetIPAddressPrefix())
			if ok {
				// if the route exists validate the owner
				routeLabels := route.Labels()
				if err := claim.ValidateOwner(routeLabels); err != nil {
					return err
				}
			}
		}
		return nil
	}
	// Route does not exist -> dry run
	route = table.NewRoute(
		pi.GetIPSubnet(),
		claim.GetDummyLabelsFromPrefix(*pi),
		map[string]any{},
	)
	if err := dryrunRib.Add(route); err != nil {
		log.Error("cannot add route", "route", route, "error", err.Error())
		return err
	}
	// get the route again and check for children
	route, ok = dryrunRib.Get(pi.GetIPSubnet())
	if !ok {
		err := fmt.Errorf("cannot get route %s which just got addded", pi.GetIPSubnet())
		log.Error(err.Error())
		return err
	}
	// check for children
	routes := route.Children(dryrunRib)
	if len(routes) > 0 {
		if err := r.validateExistingChildren(claim, routes); err != nil {
			return err
		}
	}
	// get parents
	routes = route.Parents(dryrunRib)
	if len(routes) == 0 {
		// no parents exist
		if err := validateNoParent(claim); err != nil {
			return err
		}
		return nil
	}
	// parents exist
	parentRoute := findParent(routes)
	if err := r.validateExistingParent(claim, pi, parentRoute); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (r *staticPrefixApplicator) validateExistingChildren(claim *ipambev1alpha1.IPClaim, routes table.Routes) error {
	prefixType := claim.GetIPPrefixType()

	for _, route := range routes {
		routeLabels := route.Labels()
		childClaimSummaryType := routeLabels[backend.KuidIPAMClaimSummaryTypeKey]
		childPrefixType := routeLabels[backend.KuidIPAMIPPrefixTypeKey]
		switch prefixType {
		case ipambev1alpha1.IPPrefixType_Aggregate, ipambev1alpha1.IPPrefixType_Other: // the claim is of type aggregate
			// we only allow prefixes -> validate aggregate type
			if childClaimSummaryType == string(ipambev1alpha1.IPClaimSummaryType_Address) ||
				childClaimSummaryType == string(ipambev1alpha1.IPClaimSummaryType_Range) {
				return fmt.Errorf("child with addressing %s not allowed in claim of type %s", childClaimSummaryType, prefixType)
			}
			if childPrefixType == string(ipambev1alpha1.IPPrefixType_Aggregate) {
				return fmt.Errorf("nesting %s is not possible", childPrefixType)
			}
		case ipambev1alpha1.IPPrefixType_Network, ipambev1alpha1.IPPrefixType_Pool:
			// we only allow range and addresses -> these dont have a claimType
			if childClaimSummaryType == string(ipambev1alpha1.IPClaimSummaryType_Prefix) {
				return fmt.Errorf("child with addressing %s not allowed in claim of type %s", childClaimSummaryType, prefixType)
			}
		default:
			return fmt.Errorf("invalid claimType: %s", prefixType)
		}
	}
	return nil
}

func (r *staticPrefixApplicator) validateExistingParent(claim *ipambev1alpha1.IPClaim, _ *iputil.Prefix, route table.Route) error {
	prefixType := claim.GetIPPrefixType()

	routeLabels := route.Labels()
	parentClaimSummaryType := routeLabels[backend.KuidIPAMClaimSummaryTypeKey]
	parentClaimType := routeLabels[backend.KuidIPAMIPPrefixTypeKey]
	switch prefixType {
	case ipambev1alpha1.IPPrefixType_Aggregate:
		return fmt.Errorf("parent %s/%s nesting %s/%s is not possible", route.Prefix().String(), *claim.Spec.Prefix, parentClaimType, prefixType)
	case ipambev1alpha1.IPPrefixType_Other: // the claim is of type aggregate
		// we only allow prefixes
		if parentClaimSummaryType == string(ipambev1alpha1.IPClaimSummaryType_Address) ||
			parentClaimSummaryType == string(ipambev1alpha1.IPClaimSummaryType_Range) {
			return fmt.Errorf("parent %s not allowed in claim of type %s", parentClaimSummaryType, prefixType)
		}
		if parentClaimType == string(ipambev1alpha1.IPPrefixType_Network) ||
			parentClaimType == string(ipambev1alpha1.IPPrefixType_Pool) {
			return fmt.Errorf("parent %s/%s nesting %s/%s is not possible", route.Prefix().String(), *claim.Spec.Prefix, parentClaimType, prefixType)
		}
	case ipambev1alpha1.IPPrefixType_Network, ipambev1alpha1.IPPrefixType_Pool:
		// we only allow range and addresses -> these dont have a claimType
		if parentClaimSummaryType == string(ipambev1alpha1.IPClaimSummaryType_Address) ||
			parentClaimSummaryType == string(ipambev1alpha1.IPClaimSummaryType_Range) {
			return fmt.Errorf("parent %s not allowed in claim of type %s", parentClaimSummaryType, prefixType)
		}
		if parentClaimType == string(ipambev1alpha1.IPPrefixType_Network) ||
			parentClaimType == string(ipambev1alpha1.IPPrefixType_Pool) {
			return fmt.Errorf("parent %s/%s nesting %s/%s is not possible", route.Prefix().String(), *claim.Spec.Prefix, parentClaimType, prefixType)
		}
	default:
		return fmt.Errorf("invalid claimType: %s", prefixType)
	}
	return nil
}
