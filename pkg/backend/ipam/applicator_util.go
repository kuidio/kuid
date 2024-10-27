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
	"github.com/hansthienpondt/nipam/pkg/table"
	"github.com/kform-dev/choreo/apis/kuid/backend"
)

func isParentRouteSelectable(route table.Route, prefixLength uint8) bool {
	// return the first route that has a routes with the prefixlength available
	return route.Prefix().Bits() < int(prefixLength)
}

func findMostSpecificParent(routes table.Routes) table.Route {
	parentRoute := routes[0]
	for _, route := range routes {
		if route.Prefix().Bits() > parentRoute.Prefix().Bits() {
			parentRoute = route
		}
	}
	return parentRoute
}

func getUserDefinedLabels(labels map[string]string) map[string]string {
	udmLabels := map[string]string{}
	for k, v := range labels {
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

func getExistingRoutes(existingRoutes table.Routes) []string {
	routes := []string{}
	for _, existingRoute := range existingRoutes {
		routes = append(routes, existingRoute.Prefix().String())
	}
	return routes
}
