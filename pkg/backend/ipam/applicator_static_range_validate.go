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
	"go4.org/netipx"
)

type staticRangeApplicator struct {
	name string
	applicator
}

func (r *staticRangeApplicator) Validate(ctx context.Context, claim *ipambev1alpha1.IPClaim) error {
	log := log.FromContext(ctx)
	ipRange, err := netipx.ParseIPRange(*claim.Spec.Range)
	if err != nil {
		return err
	}
	// get dryrun rib
	dryrunRib := r.cacheCtx.rib.Clone()

	// expand the prefixes
	exists := false
	for _, prefix := range ipRange.Prefixes() {
		route, ok := dryrunRib.Get(prefix)
		if ok {
			exists = true
			// if the route exists validate the owner
			routeLabels := route.Labels()
			if err := claim.ValidateOwner(routeLabels); err != nil {
				return err
			}
		} else {
			// check if some exists and others not
			if exists {
				return fmt.Errorf("some routes in the range exists and others not, it should be all or nothing")
			}
		}
	}
	if exists {
		return nil // all prefixes of the range exists so we are good
	}
	// insert in the routing table
	var singlParentPrefix string
	for _, prefix := range ipRange.Prefixes() {
		//dryrunRib := r.cacheCtx.rib.Clone()
		route := table.NewRoute(
			prefix,
			map[string]string{},
			map[string]any{},
		)
		if err := dryrunRib.Add(route); err != nil {
			log.Error("cannot add route", "route", route, "error", err.Error())
			return err
		}
		// get the route again and check for children
		route, ok := dryrunRib.Get(prefix)
		if !ok {
			err := fmt.Errorf("cannot get route %s which just got addded", prefix)
			log.Error(err.Error())
			return err
		}
		// check for children
		routes := route.Children(dryrunRib)
		if len(routes) > 0 {
			err := fmt.Errorf("cannot have children for a range %s", prefix)
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
		if err := r.validateExistingParent(claim, iputil.NewPrefixInfo(prefix), parentRoute); err != nil {
			log.Error(err.Error())
			return err
		}
		if singlParentPrefix != "" {
			if parentRoute.Prefix().String() != singlParentPrefix {
				err := fmt.Errorf("a range has to fit into a single parent prefix got %s and %s", singlParentPrefix, parentRoute.Prefix().String())
				log.Error(err.Error())
				return err
			}
		} else {
			singlParentPrefix = parentRoute.Prefix().String()
		}
	}

	return nil
}

func (r *staticRangeApplicator) validateExistingParent(_ *ipambev1alpha1.IPClaim, _ *iputil.Prefix, route table.Route) error {
	routeLabels := route.Labels()
	parentClaimType := routeLabels[backend.KuidIPAMTypeKey]

	if parentClaimType == string(ipambev1alpha1.IPClaimType_Aggregate) {
		return fmt.Errorf("a range is not possible with a parent of type %s", parentClaimType)
	}

	return nil
}
