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
	"strings"

	"github.com/hansthienpondt/nipam/pkg/table"
	"github.com/henderiw/iputil"
	"github.com/henderiw/logger/log"
	"github.com/kuidio/kuid/apis/backend"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/utils/ptr"
)

type Applicator interface {
	Apply(ctx context.Context, claim *ipambev1alpha1.IPClaim) error
	Delete(ctx context.Context, claim *ipambev1alpha1.IPClaim) error
}

type applicator struct {
	rib *table.RIB
	//pi  *iputil.Prefix
}

func (r *applicator) apply(ctx context.Context, claim *ipambev1alpha1.IPClaim, pi *iputil.Prefix) error {
	log := log.FromContext(ctx)
	// check if the prefix/claim already exists in the routing table
	// based on the name of the claim
	existingRoutes, err := r.getRoutesByOwner(ctx, claim)
	if err != nil {
		return err
	}
	// get the new routes from claim and claimed prefix
	// for network prefixes the routes can get expanded
	newRoutes := r.getRoutesFromClaim(ctx, claim, pi)
	for _, newRoute := range newRoutes {
		exists := false
		var curRoute table.Route
		for i, existingRoute := range existingRoutes {
			if existingRoute.Prefix().String() == newRoute.Prefix().String() {
				// remove the route from the existing route list as we will delete the remaining
				// existing routes later on
				existingRoutes = append(existingRoutes[:i], existingRoutes[i+1:]...)
				exists = true
				curRoute = existingRoute
			}
		}
		log.Info("apply route", "newRoute", newRoute.Prefix().String(), "exists", exists, "existsingRoutes", getExistingRoutes(existingRoutes))
		if exists {
			// update
			if err := r.updateRib(ctx, newRoute, curRoute); err != nil {
				return err
			}
		} else {
			// add
			if err := r.addRib(ctx, newRoute); err != nil {
				return err
			}
		}
	}
	for _, existingRoute := range existingRoutes {
		log.Info("delete existsingRoute", "route", existingRoute.Prefix().String())
		if err := r.rib.Delete(existingRoute); err != nil {
			log.Error("cannot delete route from rib", "route", existingRoute, "error", err.Error())
		}
	}
	r.updateClaimStatus(ctx, claim, pi)
	return nil
}

func (r *applicator) addRib(ctx context.Context, route table.Route) error {
	log := log.FromContext(ctx)
	if err := r.rib.Add(route); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			log.Error("cannot add prefix")
			return errors.Wrap(err, "cannot add prefix")
		}
	}
	return nil
}

func (r *applicator) updateRib(ctx context.Context, newRoute, existingRoute table.Route) error {
	log := log.FromContext(ctx)
	// check if the labels changed
	// if changed inform the owner GVKs through the watch
	if !labels.Equals(newRoute.Labels(), existingRoute.Labels()) {
		// workaround -> should become an atomic update
		//route = route.DeleteLabels()
		//route = route.UpdateLabel(lbls)
		log.Info("update rib with new label info", "route prefix", newRoute.Prefix().String(), "newRoute labels", newRoute.Labels(), "existsingRoute labels", existingRoute.Labels())
		if err := r.rib.Set(newRoute); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				log.Error("cannot update prefix", "error", err.Error())
				return errors.Wrap(err, "cannot update prefix")
			}
		}
		// this is an update where the labels changed
		// only update when not initializing
		// only update when the prefix is a non /32 or /128
		pi := iputil.NewPrefixInfo(newRoute.Prefix())
		if pi != nil && !pi.IsAddressPrefix() {
			log.Info("inform children of the change/update", "existingRoute", existingRoute.Prefix().String(), "labels", newRoute.Labels())
			// delete the children from the rib
			// update the once that have a nsn different from the origin
			childRoutesToBeUpdated := []table.Route{}
			for _, childRoute := range existingRoute.Children(r.rib) {
				log.Info("inform children of the change/update", "existingRoute", existingRoute.Prefix().String(), "child route", childRoute)
				if childRoute.Labels()[backend.KuidClaimNameKey] != newRoute.Labels()[backend.KuidClaimNameKey] {
					childRoutesToBeUpdated = append(childRoutesToBeUpdated, childRoute)
					if err := r.rib.Delete(childRoute); err != nil {
						log.Error("cannot delete route from rib", "route", childRoute, "error", err.Error())
						continue
					}
				}
			}
			// handler watch update to the source owner controller
			log.Info("inform children of the change/update", "existingRoute", existingRoute.Prefix().String(), "child routes", childRoutesToBeUpdated)
		}
	}
	return nil
}

func (r *applicator) updateClaimStatus(ctx context.Context, claim *ipambev1alpha1.IPClaim, pi *iputil.Prefix) {
	// update the status
	claim.Status.Prefix = ptr.To[string](pi.String())
	if claim.Spec.Kind == ipambev1alpha1.PrefixKindNetwork {
		gateway := r.getGateway(ctx, claim, *claim.Status.Prefix)
		if gateway != "" {
			claim.Status.Gateway = ptr.To[string](gateway)
		}
	}
	claim.SetConditions(conditionv1alpha1.Ready())
}

// getRoutesFromClaim return the reoutes with the assocated labels from the claim
// for network based prefixes multiple routes can be returned as they might get expanded
func (r *applicator) getRoutesFromClaim(_ context.Context, claim *ipambev1alpha1.IPClaim, pi *iputil.Prefix) []table.Route {
	routes := []table.Route{}

	labels := claim.Spec.GetUserDefinedLabels()
	labels[backend.KuidIPAMKindKey] = string(claim.Spec.Kind)
	labels[backend.KuidIPAMddressFamilyKey] = string(pi.GetAddressFamily())
	labels[backend.KuidIPAMSubnetKey] = pi.GetSubnetName()
	labels[backend.KuidClaimNameKey] = claim.Name
	labels[backend.KuidOwnerGroupKey] = claim.Spec.Owner.Group
	labels[backend.KuidOwnerVersionKey] = claim.Spec.Owner.Version
	labels[backend.KuidOwnerKindKey] = claim.Spec.Owner.Kind
	labels[backend.KuidOwnerNamespaceKey] = claim.Spec.Owner.Namespace
	labels[backend.KuidOwnerNameKey] = claim.Spec.Owner.Name
	if claim.Spec.Gateway != nil && *claim.Spec.Gateway {
		labels[backend.KuidIPAMGatewayKey] = "true"
	}

	prefix := pi.GetIPPrefix()
	if claim.Spec.Kind == ipambev1alpha1.PrefixKindNetwork {
		if claim.Spec.CreatePrefix != nil {
			switch {
			case pi.GetAddressFamily() == iputil.AddressFamilyIpv4 && pi.GetPrefixLength().Int() == 31,
				pi.GetAddressFamily() == iputil.AddressFamilyIpv6 && pi.GetPrefixLength().Int() == 127:
				routes = append(routes, r.getNetworkNetRoute(labels, pi))
			case pi.IsNorLastNorFirst():
				routes = append(routes, r.getNetworkNetRoute(labels, pi))
				routes = append(routes, r.getNetworIPAddressRoute(labels, pi))
				routes = append(routes, r.getNetworFirstAddressRoute(labels, pi))
				routes = append(routes, r.getNetworLastAddressRoute(labels, pi))
			case pi.IsFirst():
				routes = append(routes, r.getNetworkNetRoute(labels, pi))
				routes = append(routes, r.getNetworIPAddressRoute(labels, pi))
				routes = append(routes, r.getNetworLastAddressRoute(labels, pi))
			case pi.IsLast():
				routes = append(routes, r.getNetworkNetRoute(labels, pi))
				routes = append(routes, r.getNetworIPAddressRoute(labels, pi))
				routes = append(routes, r.getNetworFirstAddressRoute(labels, pi))
			}
			return routes
		} else {
			// return address
			//labels[ipamv1alpha1.NephioParentPrefixLengthKey] = r.pi.GetPrefixLength().String()
			prefix = pi.GetIPAddressPrefix()
		}
	}
	routes = append(routes, table.NewRoute(prefix, labels, map[string]any{}))
	return routes
}

func (r *applicator) getNetworkNetRoute(l map[string]string, pi *iputil.Prefix) table.Route {
	labels := map[string]string{}
	for k, v := range l {
		labels[k] = v
	}
	delete(labels, backend.KuidIPAMGatewayKey)
	return table.NewRoute(pi.GetIPSubnet(), labels, map[string]any{})
}

func (r *applicator) getNetworIPAddressRoute(l map[string]string, pi *iputil.Prefix) table.Route {
	labels := map[string]string{}
	for k, v := range l {
		labels[k] = v
	}
	return table.NewRoute(pi.GetIPAddressPrefix(), labels, map[string]any{})
}

func (r *applicator) getNetworFirstAddressRoute(l map[string]string, pi *iputil.Prefix) table.Route {
	labels := map[string]string{}
	for k, v := range l {
		labels[k] = v
	}
	delete(labels, backend.KuidIPAMGatewayKey)
	return table.NewRoute(pi.GetFirstIPPrefix(), labels, map[string]any{})
}

func (r *applicator) getNetworLastAddressRoute(l map[string]string, pi *iputil.Prefix) table.Route {
	labels := map[string]string{}
	for k, v := range l {
		labels[k] = v
	}
	delete(labels, backend.KuidIPAMGatewayKey)
	return table.NewRoute(pi.GetLastIPPrefix(), labels, map[string]any{})
}

func (r *applicator) getGateway(ctx context.Context, claim *ipambev1alpha1.IPClaim, prefix string) string {
	log := log.FromContext(ctx)
	pi, err := iputil.New(prefix)
	if err != nil {
		log.Error("cannot get gateway parent rpefix", "error", err.Error())
		return ""
	}

	gatewaySelector, err := claim.GetGatewayLabelSelector(string(pi.GetSubnetName()))
	if err != nil {
		log.Error("cannot get gateway label selector", "error", err.Error())
		return ""
	}
	log.Debug("gateway", "gatewaySelector", gatewaySelector)
	routes := r.rib.GetByLabel(gatewaySelector)
	if len(routes) > 0 {
		log.Debug("gateway", "routes", routes)
		return routes[0].Prefix().Addr().String()
	}
	return ""
}

func (r *applicator) getRoutesByOwner(_ context.Context, claim *ipambev1alpha1.IPClaim) (table.Routes, error) {
	// check if the prefix/claim already exists in the routing table
	// based on the owner and the name of the claim
	ownerSelector, err := claim.GetOwnerSelector()
	if err != nil {
		return []table.Route{}, err
	}

	routes := r.rib.GetByLabel(ownerSelector)
	if len(routes) != 0 {
		// for a prefixkind network with create prefix flag set it is possible that multiple
		// routes are returned since they were expanded
		// otherwise we expect a single route
		if len(routes) > 1 && !(claim.Spec.CreatePrefix != nil && claim.Spec.Kind == ipambev1alpha1.PrefixKindNetwork) {
			return []table.Route{}, fmt.Errorf("multiple prefixes match the nsn labelselector, %v", routes)
		}
		// route found
		return routes, nil
	}
	// no route found
	return []table.Route{}, nil
}

func (r *applicator) getRoutesByLabel(ctx context.Context, claim *ipambev1alpha1.IPClaim) table.Routes {
	log := log.FromContext(ctx)
	labelSelector, err := claim.GetLabelSelector()
	if err != nil {
		log.Error("cannot get label selector", "error", err.Error())
		return []table.Route{}
	}
	return r.rib.GetByLabel(labelSelector)
}

// Delete deletes the claimation based on the ownerslector and deletes all prefixes associated with the ownerseelctor
// if no prefixes are found, no error is returned
func (r *applicator) Delete(ctx context.Context, claim *ipambev1alpha1.IPClaim) error {
	log := log.FromContext(ctx)
	log.Info("delete")

	existingRoutes, err := r.getRoutesByOwner(ctx, claim)
	if err != nil {
		return err
	}

	for _, route := range existingRoutes {
		log = log.With("route prefix", route.Prefix())

		// this is a delete
		// only update when not initializing
		// only update when the prefix is a non /32 or /128
		// only update when the parent is a create prefix type
		if !iputil.NewPrefixInfo(route.Prefix()).IsAddressPrefix() && (claim.Spec.CreatePrefix != nil) {
			log.Info("route exists", "handle update for route", route)
			// delete the children from the rib
			// update the once that have a nsn different from the origin
			childRoutesToBeUpdated := []table.Route{}
			for _, childRoute := range route.Children(r.rib) {
				log.Info("route exists", "handle update for route", route, "child route", childRoute)
				if childRoute.Labels()[backend.KuidClaimNameKey] != claim.Name {
					childRoutesToBeUpdated = append(childRoutesToBeUpdated, childRoute)
					if err := r.rib.Delete(childRoute); err != nil {
						log.Error("cannot delete route from rib", "route", childRoute, "error", err.Error())
					}
				}
			}
			// handler watch update to the source owner controller
			log.Info("route exists", "handle update for route", route, "child routes", childRoutesToBeUpdated)
		}

		if err := r.rib.Delete(route); err != nil {
			return err
		}
	}
	return nil
}

func (r *applicator) getSelectedRouteWithPrefixLength(_ context.Context, claim *ipambev1alpha1.IPClaim, routes table.Routes, prefixLength uint8) *table.Route {
	//log := log.FromContext(ctx)
	//log.Info("claim w/o prefix", "routes", routes)

	if prefixLength == 32 || prefixLength == 128 {
		ownKindRoutes := make([]table.Route, 0)
		otherKindRoutes := make([]table.Route, 0)
		for _, route := range routes {
			switch claim.Spec.Kind {
			case ipambev1alpha1.PrefixKindLoopback:
				if route.Labels()[backend.KuidIPAMKindKey] == string(claim.Spec.Kind) {
					ownKindRoutes = append(ownKindRoutes, route)
				}
				if route.Labels()[backend.KuidIPAMKindKey] == string(ipambev1alpha1.PrefixKindAggregate) {
					otherKindRoutes = append(otherKindRoutes, route)
				}
			case ipambev1alpha1.PrefixKindNetwork:
				if route.Labels()[backend.KuidIPAMKindKey] == string(claim.Spec.Kind) {
					ownKindRoutes = append(ownKindRoutes, route)
				}
				if route.Labels()[backend.KuidIPAMKindKey] == string(ipambev1alpha1.PrefixKindAggregate) {
					otherKindRoutes = append(otherKindRoutes, route)
				}
			case ipambev1alpha1.PrefixKindPool:
				if route.Labels()[backend.KuidIPAMKindKey] == string(claim.Spec.Kind) {
					ownKindRoutes = append(ownKindRoutes, route)
				}
			default:
				// aggregates with dynamic claim always have a prefix
			}
		}
		if len(ownKindRoutes) > 0 {
			return &ownKindRoutes[0]
		}
		if len(otherKindRoutes) > 0 {
			return &otherKindRoutes[0]
		}
		return nil
	} else {
		for _, route := range routes {
			if route.Prefix().Bits() < int(prefixLength) {
				return &route
			}
		}
	}
	return nil
}

func (r *applicator) getPrefixLengthFromRoute(_ context.Context, claim *ipambev1alpha1.IPClaim, route table.Route) iputil.PrefixLength {
	if claim.Spec.PrefixLength != nil {
		return iputil.PrefixLength(*claim.Spec.PrefixLength)
	}
	// return either 32 for ipv4 and 128 for ipv6
	return iputil.PrefixLength(route.Prefix().Addr().BitLen())
}

func (r *applicator) updatePrefixInfo(_ context.Context, claim *ipambev1alpha1.IPClaim, pi *iputil.Prefix, p netip.Prefix, prefixLength iputil.PrefixLength) *iputil.Prefix {
	if claim.Spec.Kind == ipambev1alpha1.PrefixKindNetwork {
		if claim.Spec.CreatePrefix != nil {
			return iputil.NewPrefixInfo(p)
		}
		return iputil.NewPrefixInfo(netip.PrefixFrom(p.Addr(), int(pi.GetPrefixLength())))
	}
	return iputil.NewPrefixInfo(netip.PrefixFrom(p.Addr(), prefixLength.Int()))
}

// claimPrefix claims a prefix from the rib based on the claim (dynamic)
func (r *applicator) claimPrefix(ctx context.Context, claim *ipambev1alpha1.IPClaim) (*iputil.Prefix, error) {
	log := log.FromContext(ctx)

	// first check if the resource is already claimed
	existingRoutes, err := r.getRoutesByOwner(ctx, claim)
	if err != nil {
		return nil, err
	}
	found := false
	var spi *iputil.Prefix
	for _, existingRoute := range existingRoutes {
		if claim.Status.Prefix != nil {
			spi, err = iputil.New(*claim.Status.Prefix)
			if err != nil {
				return nil, err
			}
			epi := iputil.NewPrefixInfo(existingRoute.Prefix())
			if spi.GetIPAddress() == epi.GetIPAddress() {
				found = true
				break
			}
		}
	}
	if found {
		return spi, nil
	}

	// if not claimed, try to claim the ip 
	routes := r.getRoutesByLabel(ctx, claim)
	if len(routes) == 0 {
		return nil, fmt.Errorf("dynamic claim: no available routes based on the selector labels %v", claim.Spec.GetSelectorLabels())
	}

	// try to reclaim the prefix if the prefix was already claimed
	if claim.Status.Prefix != nil {
		pi, err := iputil.New(*claim.Status.Prefix)
		if err != nil {
			return nil, err
		}
		log.Info("refresh claimed prefix",
			"claimedPrefix", claim.Status.Prefix,
			"prefixlength", pi.GetPrefixLength())

		// check if the prefix is available
		p := r.rib.GetAvailablePrefixByBitLen(pi.GetIPPrefix(), uint8(pi.GetPrefixLength()))
		if p.IsValid() {
			log.Info("refresh claimed prefix finished",
				"claimedPrefix", claim.Status.Prefix)
			// previously claimed prefix is available and reassigned
			return iputil.NewPrefixInfo(p), nil
		}
		log.Info("refresh claim prefix not available",
			"claimedPrefix", claim.Status.Prefix,
			"prefixlength", pi.GetPrefixLength())
	}

	// If there was no previously claimed prefix or the reclaim of the prefix failed
	// we try to claim a new prefix
	// prefixlength is either set by the claim request, if not it is derived from the
	// returned prefix and address family (32 for ipv4 and 128 for ipv6)
	prefixLength := r.getPrefixLengthFromRoute(ctx, claim, routes[0])
	selectedRoute := r.getSelectedRouteWithPrefixLength(ctx, claim, routes, uint8(prefixLength.Int()))
	if selectedRoute == nil {
		return nil, fmt.Errorf("no route found with requested prefixLength: %d", prefixLength)
	}
	pi := iputil.NewPrefixInfo(selectedRoute.Prefix())
	log.Info("new claim", "selectedRoute", selectedRoute)
	p := r.rib.GetAvailablePrefixByBitLen(pi.GetIPPrefix(), uint8(prefixLength.Int()))
	if !p.IsValid() {
		return nil, errors.New("no free prefix found")
	}
	log.Info("new claim",
		"pi prefix", pi,
		"p prefix", p,
		"prefixLength", pi.GetPrefixLength(),
	)
	pi = r.updatePrefixInfo(ctx, claim, pi, p, prefixLength)
	log.Info("new claim",
		"claimedPrefix", pi.Prefix.String())
	return pi, nil
}

func getExistingRoutes(existingRoutes table.Routes) []string {
	routes := []string{}
	for _, existingRoute := range existingRoutes {
		routes = append(routes, existingRoute.Prefix().String())
	}
	return routes
}
