/*
Copyright 2023 The Nephio Authors.

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

package backend

import "k8s.io/apimachinery/pkg/util/sets"

const (
	// system defined common
	KuidOwnerGroupKey     = "be.kuid.dev/owner-group"
	KuidOwnerVersionKey   = "be.kuid.dev/owner-version"
	KuidOwnerKindKey      = "be.kuid.dev/owner-kind"
	KuidOwnerNameKey      = "be.kuid.dev/owner-name"
	KuidOwnerNamespaceKey = "be.kuid.dev/owner-namespace"
	KuidClaimNameKey      = "be.kuid.dev/claim-name"
	// system defined ipam
	KuidIPAMKindKey         = "ipam.be.kuid.dev/kind"
	KuidIPAMddressFamilyKey = "ipam.be.kuid.dev/address-family"
	KuidIPAMSubnetKey       = "ipam.be.kuid.dev/subnet" // this is the subnet in prefix annotation used for GW selection
	//KuidIPAMPoolKey         = "ipam.be.kuid.dev/pool"
	KuidIPAMGatewayKey = "ipam.be.kuid.dev/gateway"
	KuidIPAMIndexKey   = "ipam.be.kuid.dev/index"
	// user defined common
	//NephioClusterNameKey       = "nephio.org/cluster-name"
	//NephioSiteNameKey          = "nephio.org/site-name"
	//NephioRegionKey            = "nephio.org/region"
	//NephioAvailabilityZoneKey  = "nephio.org/availability-zone"
	//NephioInterfaceKey         = "nephio.org/interface"
	//NephioNetworkNameKey       = "nephio.org/network-name"
	//NephioPurposeKey           = "nephio.org/purpose"
	//NephioApplicationPartOfKey = "app.kubernetes.io/part-of"
	//NephioIndexKey = "nephio.org/index"
	// status ipam
	//NephioClaimedPrefix  = "nephio.org/claimed-prefix"
	//NephioClaimedGateway = "nephio.org/claimed-gateway"
)

var BackendSystemKeys = sets.New[string](KuidOwnerGroupKey,
	KuidOwnerVersionKey,
	KuidOwnerKindKey,
	KuidOwnerNameKey,
	KuidOwnerNamespaceKey,
	KuidClaimNameKey,
)

var BackendIPAMSystemKeys = sets.New[string](KuidOwnerGroupKey,
	KuidIPAMKindKey,
	KuidIPAMddressFamilyKey,
	KuidIPAMSubnetKey,
	KuidIPAMGatewayKey,
	KuidIPAMIndexKey,
)
