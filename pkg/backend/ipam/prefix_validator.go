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
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	ipamresv1alpha1 "github.com/kuidio/kuid/apis/resource/ipam/v1alpha1"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
)

type prefixClaimValidator struct {
	rib *table.RIB
}

func (r *prefixClaimValidator) Validate(ctx context.Context, claim *ipambev1alpha1.IPClaim) error {
	log := log.FromContext(ctx)
	pi, err := iputil.New(*claim.Spec.Prefix)
	if err != nil {
		return err
	}
	// get dryrun rib
	dryrunRib := r.rib.Clone()

	// check if the prefix/subnet exists
	route, ok := dryrunRib.Get(pi.GetIPSubnet())
	if ok {
		// route exists
		routeLabels := route.Labels()
		if claim.Spec.CreatePrefix != nil {
			// this is a create prefix, not an address prefix
			if err := claim.ValidateOwner(routeLabels); err != nil {
				return err
			}
			if claim.Spec.Kind != ipambev1alpha1.PrefixKindNetwork {
				// for non network prefixes we can stop,
				// for network prefixes we still need to validate the address portion
				return nil
			}
		}
		// in case of a network prefix we need to validate the address
		if claim.Spec.Kind == ipambev1alpha1.PrefixKindNetwork {
			route, ok = dryrunRib.Get(pi.GetIPAddressPrefix())
			if ok {
				// if the route exists validate the owner
				if err := claim.ValidateOwner(routeLabels); err != nil {
					return err
				}
			}
		}
		// get labels again
		routeLabels = route.Labels()
		// this is an address prefix
		if err := claim.ValidateOwner(routeLabels); err != nil {
			return err
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
		if err := claim.ValidateExistingChildren(routes[0]); err != nil {
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
	if err := claim.ValidateExistingParent(pi, parentRoute); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func findParent(routes table.Routes) table.Route {
	parentRoute := routes[0]
	for _, route := range routes {
		if route.Prefix().Bits() > parentRoute.Prefix().Bits() {
			parentRoute = route
		}
	}
	return parentRoute
}

func validateNoParent(ipClaim *ipambev1alpha1.IPClaim) error {
	fmt.Println("validateNoParent claim", ipClaim)
	if ipClaim.Spec.Owner.Group != ipamresv1alpha1.SchemeGroupVersion.Group ||
		ipClaim.Spec.Owner.Version != ipamresv1alpha1.SchemeGroupVersion.Version ||
		ipClaim.Spec.Owner.Kind != ipamresv1alpha1.NetworkInstanceKind {
		ownerRef := commonv1alpha1.OwnerReference{
			Group:   ipamresv1alpha1.SchemeGroupVersion.Group,
			Version: ipamresv1alpha1.SchemeGroupVersion.Version,
			Kind:    ipamresv1alpha1.NetworkInstanceKind,
		}
		return fmt.Errorf("an agregate route is required %s/%s", ipClaim.Spec.Owner.String(), ownerRef)
	}
	return nil // an aggregate coming from a network Instance can be created
}
