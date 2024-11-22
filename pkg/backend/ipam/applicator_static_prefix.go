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
	"github.com/kuidio/kuid/apis/backend/ipam"
)

type staticPrefixApplicator struct {
	name string
	applicator
	exists bool
}

func (r *staticPrefixApplicator) Validate(ctx context.Context, claim *ipam.IPClaim) error {
	log := log.FromContext(ctx)
	log.Debug("validate")
	if claim.Spec.Prefix == nil {
		return fmt.Errorf("cannot claim a static ip prefix without a prefix")
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

func (r *staticPrefixApplicator) validateExists(_ context.Context, claim *ipam.IPClaim) (bool, error) {
	pi, err := iputil.New(*claim.Spec.Prefix)
	if err != nil {
		return false, err
	}

	// There is 2 scenario's:
	// a regular prefix w/o address: 10.0.0.0/24
	// an address based prefix (only allowed for claimType network): 10.0.0.1/24
	// to accomodate for both scenario;s we check the subnet prefix 10.0.0.0/24 for either scenario
	route, ok := r.cacheInstanceCtx.rib.Get(pi.GetIPSubnet())
	if ok {
		// entry exists; validate the owner to see if someone else owns this prefix
		labels := route.Labels()
		if err := claim.ValidateOwner(labels); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (r *staticPrefixApplicator) validateParents(ctx context.Context, claim *ipam.IPClaim) error {
	pi, err := iputil.New(*claim.Spec.Prefix)
	if err != nil {
		return err
	}

	parentRoutes := r.cacheInstanceCtx.rib.Parents(pi.GetIPSubnet())
	if len(parentRoutes) == 0 {
		// a parent route is always required unless you are an aggregate route owned
		// by an IP Index
		if !claim.IsOwnedByIPIndex() {
			return fmt.Errorf("no parent found, only possible for routes owned by IPIndex")
		}
		return nil
	}
	parentRoute := findMostSpecificParent(parentRoutes)
	if err := r.validateParent(ctx, parentRoute, claim); err != nil {
		return err
	}
	return nil

}

func (r *staticPrefixApplicator) validateParent(_ context.Context, route table.Route, claim *ipam.IPClaim) error {
	prefixType := claim.GetIPPrefixType()

	routeLabels := route.Labels()
	parentClaimSummaryType := routeLabels[backend.KuidIPAMClaimSummaryTypeKey]
	parentClaimPrefixType := routeLabels[backend.KuidIPAMIPPrefixTypeKey]
	switch prefixType {
	case ipam.IPPrefixType_Other:
		if parentClaimPrefixType == string(ipam.IPPrefixType_Network) {
			return fmt.Errorf("parent %s/%s nesting %s/%s is not possible", route.Prefix().String(), *claim.Spec.Prefix, parentClaimPrefixType, prefixType)
		}
		return nil
	case ipam.IPPrefixType_Network:
		// we only allow range and addresses -> these dont have a claimType
		if parentClaimSummaryType == string(ipam.IPClaimSummaryType_Address) ||
			parentClaimSummaryType == string(ipam.IPClaimSummaryType_Range) {
			return fmt.Errorf("parent %s not allowed in claim of type %s", parentClaimSummaryType, prefixType)
		}
		if parentClaimPrefixType == string(ipam.IPPrefixType_Network) {
			return fmt.Errorf("parent %s/%s nesting %s/%s is not possible", route.Prefix().String(), *claim.Spec.Prefix, parentClaimPrefixType, prefixType)
		}
	default:
		return fmt.Errorf("invalid prefixType: %s", prefixType)
	}
	return nil
}

func (r *staticPrefixApplicator) validateChildren(_ context.Context, claim *ipam.IPClaim) error {
	// network, other
	prefixType := claim.GetIPPrefixType()
	pi, err := iputil.New(*claim.Spec.Prefix)
	if err != nil {
		return err
	}

	// There is 2 scenario's:
	// a regular prefix w/o address: 10.0.0.0/24
	// an address based prefix (only allowed for claimType network): 10.0.0.1/24
	// to accomodate for both scenario;s we check the subnet prefix 10.0.0.0/24 for either scenario
	childRoutes := r.cacheInstanceCtx.rib.Children(pi.GetIPSubnet())
	for _, childRoute := range childRoutes {
		routeLabels := childRoute.Labels()
		childClaimSummaryType := routeLabels[backend.KuidIPAMClaimSummaryTypeKey]
		//childPrefixType := routeLabels[backend.KuidIPAMIPPrefixTypeKey]
		switch prefixType {
		case ipam.IPPrefixType_Other: // the claim is of type aggregate
			// TODO insertion of prefixes
			/*
			if childClaimSummaryType == string(ipam.IPClaimSummaryType_Address) ||
				childClaimSummaryType == string(ipam.IPClaimSummaryType_Range) {
				return fmt.Errorf("child with addressing %s not allowed in claim of type %s", childClaimSummaryType, prefixType)
			}
			if childPrefixType == string(ipam.IPPrefixType_Aggregate) {
				return fmt.Errorf("nesting %s is not possible", childPrefixType)
			}
				*/
		case ipam.IPPrefixType_Network:
			// we only allow range and addresses -> these dont have a claimType
			if childClaimSummaryType == string(ipam.IPClaimSummaryType_Prefix) {
				return fmt.Errorf("child with addressing %s not allowed in claim of type %s", childClaimSummaryType, prefixType)
			}
		default:
			return fmt.Errorf("invalid claimType: %s", prefixType)
		}
	}
	return nil
}

func (r *staticPrefixApplicator) Apply(ctx context.Context, claim *ipam.IPClaim) error {
	log := log.FromContext(ctx)
	log.Debug("apply")
	pi, err := iputil.New(*claim.Spec.Prefix)
	if err != nil {
		return err
	}

	if err := r.apply(ctx, claim, []*iputil.Prefix{pi}, false, map[string]string{}); err != nil {
		return err
	}
	r.updateClaimPrefixStatus(ctx, claim, pi)
	return nil
}

func (r *staticPrefixApplicator) Delete(ctx context.Context, claim *ipam.IPClaim) error {
	log := log.FromContext(ctx)
	log.Debug("delete")
	return r.delete(ctx, claim)
}
