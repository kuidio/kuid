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
	"errors"
	"fmt"
	"strings"

	"github.com/hansthienpondt/nipam/pkg/table"
	"github.com/henderiw/idxtable/pkg/iptable"
	"github.com/henderiw/iputil"
	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	ipamresv1alpha1 "github.com/kuidio/kuid/apis/resource/ipam/v1alpha1"
	"go4.org/netipx"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/utils/ptr"
)

type Validator interface {
	Validate(ctx context.Context, claim *ipambev1alpha1.IPClaim) error
}

type Applicator interface {
	Validate(ctx context.Context, claim *ipambev1alpha1.IPClaim) error
	Apply(ctx context.Context, claim *ipambev1alpha1.IPClaim) error
	Delete(ctx context.Context, claim *ipambev1alpha1.IPClaim) error
}

type applicator struct {
	cacheCtx *CacheContext
}

// apply only works on the main rib
func (r *applicator) apply(ctx context.Context, claim *ipambev1alpha1.IPClaim, pis []*iputil.Prefix, networkParent bool, parentLabels map[string]string) error {
	log := log.FromContext(ctx)
	// check if the prefix/claim already exists in the routing table
	// based on the name of the claim
	existingRoutes, err := r.getRoutesByOwner(ctx, claim)
	if err != nil {
		return err
	}
	// get the new routes from claim and claimed prefix
	// for network prefixes the routes can get expanded
	newRoutes := table.Routes{}
	for _, pi := range pis {
		pi := pi
		newRoutes = append(newRoutes, getRoutesFromClaim(ctx, claim, pi, networkParent, parentLabels)...)

	}
	for _, newRoute := range newRoutes {
		fmt.Println("newRoute", newRoute.Prefix().String())
		newRoute := newRoute
		exists := false
		var curRoute table.Route
		for i, existingRoute := range existingRoutes[""] {
			if existingRoute.Prefix().String() == newRoute.Prefix().String() {
				// remove the route from the existing route list as we will delete the remaining
				// existing routes later on
				existingRoutes[""] = append(existingRoutes[""][:i], existingRoutes[""][i+1:]...)
				exists = true
				curRoute = existingRoute
			}
		}
		log.Info("apply route", "newRoute", newRoute.Prefix().String(), "exists", exists, "existsingRoutes", getExistingRoutes(existingRoutes[""]))
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
	for _, existingRoute := range existingRoutes[""] {
		log.Info("delete existsingRoute", "route", existingRoute.Prefix().String())
		if err := r.cacheCtx.rib.Delete(existingRoute); err != nil {
			log.Error("cannot delete route from rib", "route", existingRoute, "error", err.Error())
		}
	}
	return nil
}

func (r *applicator) applyRange(ctx context.Context, claim *ipambev1alpha1.IPClaim, ipRange netipx.IPRange) error {
	k := store.ToKey(claim.Name)
	if _, err := r.cacheCtx.ranges.Get(ctx, k); err != nil {
		ipTable, err := iptable.New(ipRange.From(), ipRange.To())
		if err != nil {
			return err
		}
		if err := r.cacheCtx.ranges.Create(ctx, k, ipTable); err != nil {
			return err
		}
	}
	return nil
}

func (r *applicator) applyAddressInRange(ctx context.Context, claim *ipambev1alpha1.IPClaim, pi *iputil.Prefix, rangeName string, labels map[string]string) error {
	k := store.ToKey(rangeName)
	ipTable, err := r.cacheCtx.ranges.Get(ctx, k)
	if err != nil {
		return err
	}
	routes := getRoutesFromClaim(ctx, claim, pi, false, labels)
	addr := pi.Addr().String()
	route, err := ipTable.Get(addr)
	if err != nil {
		ipTable.Claim(pi.Addr().String(), routes[0])
		return nil
	}
	if err := claim.ValidateOwner(route.Labels()); err != nil {
		return err
	}
	return ipTable.Update(addr, routes[0])
}

func (r *applicator) addRib(ctx context.Context, route table.Route) error {
	log := log.FromContext(ctx)
	if err := r.cacheCtx.rib.Add(route); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			log.Error("cannot add prefix")
			return fmt.Errorf("cannot add prefix, err: %s", err.Error())
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
		if err := r.cacheCtx.rib.Set(newRoute); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				log.Error("cannot update prefix", "error", err.Error())
				return fmt.Errorf("cannot update prefix, err: %s", err.Error())
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
			for _, childRoute := range existingRoute.Children(r.cacheCtx.rib) {
				log.Info("inform children of the change/update", "existingRoute", existingRoute.Prefix().String(), "child route", childRoute)
				if childRoute.Labels()[backend.KuidClaimNameKey] != newRoute.Labels()[backend.KuidClaimNameKey] {
					childRoutesToBeUpdated = append(childRoutesToBeUpdated, childRoute)
					if err := r.cacheCtx.rib.Delete(childRoute); err != nil {
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

func (r *applicator) updateClaimAddressStatus(ctx context.Context, claim *ipambev1alpha1.IPClaim, pi *iputil.Prefix, networkParent bool) {
	// update the status
	claim.Status.Address = ptr.To[string](pi.String())
	if claim.GetIPPrefixType() == ipambev1alpha1.IPPrefixType_Network || networkParent {
		defaultGateway := r.getDefaultGateway(ctx, claim, pi)
		if defaultGateway != "" {
			claim.Status.DefaultGateway = ptr.To[string](defaultGateway)
		}
	}
	claim.SetConditions(conditionv1alpha1.Ready())
}

func (r *applicator) updateClaimPrefixStatus(ctx context.Context, claim *ipambev1alpha1.IPClaim, pi *iputil.Prefix) {
	// update the status
	claim.Status.Prefix = ptr.To[string](pi.String())
	if claim.GetIPPrefixType() == ipambev1alpha1.IPPrefixType_Network {
		defaultGateway := r.getDefaultGateway(ctx, claim, pi)
		if defaultGateway != "" {
			claim.Status.DefaultGateway = ptr.To[string](defaultGateway)
		}
	}
	claim.SetConditions(conditionv1alpha1.Ready())
}

func (r *applicator) updateClaimRangeStatus(_ context.Context, claim *ipambev1alpha1.IPClaim) {
	// update the status
	claim.Status.Range = claim.Spec.Range
	claim.SetConditions(conditionv1alpha1.Ready())
}

// getRoutesFromClaim return the reoutes with the assocated labels from the claim
// for network based prefixes multiple routes can be returned as they might get expanded
func getRoutesFromClaim(_ context.Context, claim *ipambev1alpha1.IPClaim, pi *iputil.Prefix, networkParent bool, parentLabels map[string]string) []table.Route {
	routes := []table.Route{}

	ipClaimType, _ := claim.GetIPClaimType()
	labels := claim.Spec.GetUserDefinedLabels()
	for k, v := range parentLabels {
		labels[k] = v
	}
	// system defined labels
	labels[backend.KuidIPAMIPPrefixTypeKey] = string(claim.GetIPPrefixType())
	labels[backend.KuidIPAMClaimSummaryTypeKey] = string(claim.GetIPClaimSummaryType())
	labels[backend.KuidIPAMClaimTypeKey] = string(ipClaimType)
	labels[backend.KuidIPAMddressFamilyKey] = string(pi.GetAddressFamily())
	labels[backend.KuidIPAMSubnetKey] = pi.GetSubnetName()
	labels[backend.KuidClaimNameKey] = claim.Name
	labels[backend.KuidOwnerGroupKey] = claim.Spec.Owner.Group
	labels[backend.KuidOwnerVersionKey] = claim.Spec.Owner.Version
	labels[backend.KuidOwnerKindKey] = claim.Spec.Owner.Kind
	labels[backend.KuidOwnerNamespaceKey] = claim.Spec.Owner.Namespace
	labels[backend.KuidOwnerNameKey] = claim.Spec.Owner.Name
	if claim.Spec.DefaultGateway != nil && *claim.Spec.DefaultGateway {
		labels[backend.KuidIPAMDefaultGatewayKey] = "true"
	}

	prefix := pi.GetIPPrefix()
	// networkParent is there for dynamic addresses as we dont know ahead of time
	// if the dynamic address matches a network or other parent prefix
	if claim.GetIPPrefixType() == ipambev1alpha1.IPPrefixType_Network || networkParent {
		if claim.Spec.CreatePrefix != nil {
			switch {
			case pi.GetAddressFamily() == iputil.AddressFamilyIpv4 && pi.GetPrefixLength().Int() == 31,
				pi.GetAddressFamily() == iputil.AddressFamilyIpv6 && pi.GetPrefixLength().Int() == 127:
				routes = append(routes, getNetworkNetRoute(labels, pi))
			case pi.IsNorLastNorFirst():
				routes = append(routes, getNetworkNetRoute(labels, pi))
				routes = append(routes, getNetworIPAddressRoute(labels, pi))
				routes = append(routes, getNetworFirstAddressRoute(labels, pi))
				routes = append(routes, getNetworLastAddressRoute(labels, pi))
			case pi.IsFirst():
				routes = append(routes, getNetworkNetRoute(labels, pi))
				routes = append(routes, getNetworIPAddressRoute(labels, pi))
				routes = append(routes, getNetworLastAddressRoute(labels, pi))
			case pi.IsLast():
				routes = append(routes, getNetworkNetRoute(labels, pi))
				routes = append(routes, getNetworIPAddressRoute(labels, pi))
				routes = append(routes, getNetworFirstAddressRoute(labels, pi))
			}
			return routes
		} else {
			// return address
			//labels[ipamv1alpha1.NephioParentPrefixLengthKey] = r.pi.GetPrefixLength().String()
			//fmt.Println("getRoutesFromClaim addressPrefix")
			prefix = pi.GetIPAddressPrefix()
		}
	}
	//fmt.Println("getRoutesFromClaim", claim.GetInfo(), pi.Prefix.String())
	routes = append(routes, table.NewRoute(prefix, labels, map[string]any{}))
	return routes
}

func getNetworkNetRoute(l map[string]string, pi *iputil.Prefix) table.Route {
	labels := map[string]string{}
	for k, v := range l {
		labels[k] = v
	}
	delete(labels, backend.KuidIPAMDefaultGatewayKey)
	delete(labels, backend.KuidINVEndpointKey)
	return table.NewRoute(pi.GetIPSubnet(), labels, map[string]any{})
}

func getNetworIPAddressRoute(l map[string]string, pi *iputil.Prefix) table.Route {
	labels := map[string]string{}
	for k, v := range l {
		labels[k] = v
	}
	if pi.IsFirst() || pi.IsLast() {
		delete(labels, backend.KuidIPAMDefaultGatewayKey)
	}
	return table.NewRoute(pi.GetIPAddressPrefix(), labels, map[string]any{})
}

func getNetworFirstAddressRoute(l map[string]string, pi *iputil.Prefix) table.Route {
	labels := map[string]string{}
	for k, v := range l {
		labels[k] = v
	}
	delete(labels, backend.KuidIPAMDefaultGatewayKey)
	delete(labels, backend.KuidINVEndpointKey)
	return table.NewRoute(pi.GetFirstIPPrefix(), labels, map[string]any{})
}

func getNetworLastAddressRoute(l map[string]string, pi *iputil.Prefix) table.Route {
	labels := map[string]string{}
	for k, v := range l {
		labels[k] = v
	}
	delete(labels, backend.KuidIPAMDefaultGatewayKey)
	delete(labels, backend.KuidINVEndpointKey)
	return table.NewRoute(pi.GetLastIPPrefix(), labels, map[string]any{})
}

func (r *applicator) getDefaultGateway(ctx context.Context, claim *ipambev1alpha1.IPClaim, pi *iputil.Prefix) string {
	log := log.FromContext(ctx)
	/*
		pi, err := iputil.New(prefix)
		if err != nil {
			log.Error("cannot get gateway parent rpefix", "error", err.Error())
			return ""
		}
	*/

	defaultGatewaySelector, err := claim.GetDefaultGatewayLabelSelector(string(pi.GetSubnetName()))
	if err != nil {
		log.Error("cannot get gateway label selector", "error", err.Error())
		return ""
	}
	log.Debug("defaultGateway", "defaultGatewaySelector", defaultGatewaySelector)
	routes := r.cacheCtx.rib.GetByLabel(defaultGatewaySelector)
	if len(routes) > 0 {
		log.Debug("defaultGateway", "routes", routes)
		return routes[0].Prefix().Addr().String()
	}
	return ""
}

func (r *applicator) getRoutesByOwner(ctx context.Context, claim *ipambev1alpha1.IPClaim) (map[string]table.Routes, error) {
	ribRoutes := map[string]table.Routes{}
	// check if the prefix/claim already exists in the routing table
	// based on the owner and the name of the claim
	ownerSelector, err := claim.GetOwnerSelector()
	if err != nil {
		return ribRoutes, err
	}

	claimSummaryType := claim.GetIPClaimSummaryType()
	claimPrefixType := claim.GetIPPrefixType()

	ribRoutes[""] = r.cacheCtx.rib.GetByLabel(ownerSelector)
	if len(ribRoutes[""]) != 0 {
		// ranges and prefixes using network type can have multiple plrefixes
		if len(ribRoutes[""]) > 1 && (claimSummaryType == ipambev1alpha1.IPClaimSummaryType_Address ||
			claimSummaryType == ipambev1alpha1.IPClaimSummaryType_Prefix && claimPrefixType != ipambev1alpha1.IPPrefixType_Network) {
			return ribRoutes, fmt.Errorf("multiple prefixes match the owner, %v", ribRoutes[""])
		}
		// route found
		//return "", routes, nil
	}
	// add the search in the iptable
	if claimSummaryType == ipambev1alpha1.IPClaimSummaryType_Address {
		var errm error
		r.cacheCtx.ranges.List(ctx, func(ctx context.Context, k store.Key, ipTable iptable.IPTable) {
			ribRoutes[k.Name] = ipTable.GetByLabel(ownerSelector)
			if len(ribRoutes[k.Name]) > 1 {
				errm = errors.Join(errm, fmt.Errorf("multiple address match the owner, %v", ribRoutes[k.Name]))
				return
			}
		})
		if errm != nil {
			return ribRoutes, errm
		}
	}
	return ribRoutes, nil
}

func (r *applicator) getRoutesByLabel(ctx context.Context, claim *ipambev1alpha1.IPClaim) table.Routes {
	log := log.FromContext(ctx)
	labelSelector, err := claim.GetLabelSelector()
	if err != nil {
		log.Error("cannot get label selector", "error", err.Error())
		return []table.Route{}
	}
	return r.cacheCtx.rib.GetByLabel(labelSelector)
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

	for ribName, existingRoutes := range existingRoutes {
		if ribName == "" {
			for _, existingRoute := range existingRoutes {
				log = log.With("route prefix", existingRoute.Prefix())
				// this is a delete
				// only update when not initializing
				// only update when the prefix is a non /32 or /128
				// only update when the parent is a create prefix type
				pi := iputil.NewPrefixInfo(existingRoute.Prefix())
				if pi != nil && !pi.IsAddressPrefix() {
					log.Info("inform children of the delete", "existingRoute", existingRoute.Prefix().String(), "labels", existingRoute.Labels())
					// delete the children from the rib
					// update the once that have a nsn different from the origin
					childRoutesToBeUpdated := []table.Route{}
					for _, childRoute := range existingRoute.Children(r.cacheCtx.rib) {
						log.Info("route exists", "handle delete for route", existingRoute, "child route", childRoute)
						if childRoute.Labels()[backend.KuidClaimNameKey] != claim.Name {
							childRoutesToBeUpdated = append(childRoutesToBeUpdated, childRoute)
							if err := r.cacheCtx.rib.Delete(childRoute); err != nil {
								log.Error("cannot delete route from rib", "route", childRoute, "error", err.Error())
							}
						}
					}
					// handler watch update to the source owner controller
					log.Info("route exists", "handle update for route", existingRoute, "child routes", childRoutesToBeUpdated)
				}

				if err := r.cacheCtx.rib.Delete(existingRoute); err != nil {
					return err
				}

				// check if the route was a range -> if so delete the range table
				routeLabels := existingRoute.Labels()
				//parentClaimType := ipambev1alpha1.GetIPClaimTypeFromString(routeLabels[backend.KuidIPAMTypeKey])
				parentSummaryType := ipambev1alpha1.GetIPClaimSummaryTypeFromString(routeLabels[backend.KuidIPAMClaimSummaryTypeKey])
				parentClaimName := routeLabels[backend.KuidClaimNameKey]

				if parentSummaryType == ipambev1alpha1.IPClaimSummaryType_Range {
					k := store.ToKey(parentClaimName) // this is the name of the range
					if _, err := r.cacheCtx.ranges.Get(ctx, k); err == nil {
						// the table exists -> delete it
						if err := r.cacheCtx.ranges.Delete(ctx, k); err != nil {
							return err
						}
					}
				}
			}
		} else {
			k := store.ToKey(ribName)
			if len(existingRoutes) > 0 {
				if ipTable, err := r.cacheCtx.ranges.Get(ctx, k); err == nil {
					// the table exists
					for _, existingRoute := range existingRoutes {
						if _, err := ipTable.Get(existingRoute.Prefix().Addr().String()); err == nil {
							if err := ipTable.Release(existingRoute.Prefix().Addr().String()); err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func isParentRouteSelectable(route table.Route, prefixLength uint8) bool {
	// return the first route that has a routes with the prefixlength available
	return route.Prefix().Bits() < int(prefixLength)
}

func getExistingRoutes(existingRoutes table.Routes) []string {
	routes := []string{}
	for _, existingRoute := range existingRoutes {
		routes = append(routes, existingRoute.Prefix().String())
	}
	return routes
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

func getUserDefinedLabels(labels map[string]string) map[string]string {
	udmLabels := map[string]string{}
	for k, v := range labels {
		fmt.Println("getUserDefinedLabels", k, v)
		if backend.BackendIPAMSystemKeys.Has(k) {
			continue
		}
		if backend.BackendSystemKeys.Has(k) {
			continue
		}
		udmLabels[k] = v
	}
	return udmLabels
}
